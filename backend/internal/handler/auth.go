package handler

import (
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"database/sql"
	"encoding/hex"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"note-thing/backend/internal/model"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

type AuthHandler struct {
	DB        *sql.DB
	JWTSecret string
}

func (h *AuthHandler) oauthConfig() *oauth2.Config {
	return &oauth2.Config{
		ClientID:     os.Getenv("GOOGLE_CLIENT_ID"),
		ClientSecret: os.Getenv("GOOGLE_CLIENT_SECRET"),
		Endpoint:     google.Endpoint,
		RedirectURL:  os.Getenv("GOOGLE_REDIRECT_URL"),
		Scopes:       []string{"openid", "email", "profile"},
	}
}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	state := h.signedState()
	http.Redirect(w, r, h.oauthConfig().AuthCodeURL(state, oauth2.SetAuthURLParam("prompt", "select_account")), http.StatusTemporaryRedirect)
}

func (h *AuthHandler) Callback(w http.ResponseWriter, r *http.Request) {
	state := r.URL.Query().Get("state")
	if !h.verifyState(state) {
		respondError(w, http.StatusBadRequest, "invalid oauth state")
		return
	}

	code := r.URL.Query().Get("code")
	token, err := h.oauthConfig().Exchange(r.Context(), code)
	if err != nil {
		log.Printf("oauth exchange failed: %v", err)
		respondError(w, http.StatusBadRequest, "oauth exchange failed")
		return
	}

	client := h.oauthConfig().Client(r.Context(), token)
	resp, err := client.Get("https://www.googleapis.com/oauth2/v2/userinfo")
	if err != nil {
		log.Printf("fetch userinfo failed: %v", err)
		respondError(w, http.StatusInternalServerError, "failed to fetch user info")
		return
	}
	defer resp.Body.Close()

	var userInfo struct {
		ID      string `json:"id"`
		Email   string `json:"email"`
		Name    string `json:"name"`
		Picture string `json:"picture"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&userInfo); err != nil {
		respondError(w, http.StatusInternalServerError, "failed to decode user info")
		return
	}

	// Upsert user
	user, err := model.UpsertUser(r.Context(), h.DB, userInfo.ID, userInfo.Email, userInfo.Name, userInfo.Picture)
	if err != nil {
		log.Printf("upsert user failed: %v", err)
		respondError(w, http.StatusInternalServerError, "failed to create user")
		return
	}

	// Create default notebook if new user (check if they have any notebooks)
	var nbCount int
	err = h.DB.QueryRowContext(r.Context(), `SELECT COUNT(*) FROM notebooks WHERE user_id = $1`, user.ID).Scan(&nbCount)
	if err != nil {
		log.Printf("count notebooks failed: %v", err)
		respondError(w, http.StatusInternalServerError, "internal error")
		return
	}
	if nbCount == 0 {
		tx, err := h.DB.BeginTx(r.Context(), nil)
		if err != nil {
			respondError(w, http.StatusInternalServerError, "internal error")
			return
		}
		defer tx.Rollback()
		if _, err := model.CreateDefaultNotebook(r.Context(), tx, user.ID); err != nil {
			log.Printf("create default notebook failed: %v", err)
			respondError(w, http.StatusInternalServerError, "internal error")
			return
		}
		tx.Commit()
	}

	// Sign JWT
	jwtToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub":   user.ID,
		"email": user.Email,
		"iat":   time.Now().Unix(),
		"exp":   time.Now().Add(24 * time.Hour).Unix(),
	})
	signed, err := jwtToken.SignedString([]byte(h.JWTSecret))
	if err != nil {
		respondError(w, http.StatusInternalServerError, "failed to sign token")
		return
	}

	// Redirect to frontend with token
	frontendURL := os.Getenv("FRONTEND_URL")
	if frontendURL != "" {
		http.Redirect(w, r, frontendURL+"/auth/callback?token="+signed, http.StatusTemporaryRedirect)
		return
	}

	// Fallback: return JSON
	respondJSON(w, http.StatusOK, map[string]any{
		"token": signed,
		"user":  user,
	})
}

func (h *AuthHandler) Me(w http.ResponseWriter, r *http.Request) {
	userID := userIDFromContext(r)
	user, err := model.GetUserByID(r.Context(), h.DB, userID)
	if err != nil {
		respondError(w, http.StatusNotFound, "user not found")
		return
	}
	respondJSON(w, http.StatusOK, user)
}

func (h *AuthHandler) UpdateMe(w http.ResponseWriter, r *http.Request) {
	userID := userIDFromContext(r)
	var input struct {
		Name string `json:"name"`
	}
	if err := decodeJSON(r, &input); err != nil || strings.TrimSpace(input.Name) == "" {
		respondError(w, http.StatusBadRequest, "name is required")
		return
	}
	user, err := model.UpdateUserName(r.Context(), h.DB, userID, strings.TrimSpace(input.Name))
	if err != nil {
		respondError(w, http.StatusInternalServerError, "failed to update user")
		return
	}
	respondJSON(w, http.StatusOK, user)
}

func (h *AuthHandler) DeleteMe(w http.ResponseWriter, r *http.Request) {
	userID := userIDFromContext(r)
	if err := model.DeleteUser(r.Context(), h.DB, userID); err != nil {
		respondError(w, http.StatusInternalServerError, "failed to delete account")
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

// signedState creates an HMAC-signed random state token.
// Format: <random_hex>.<signature_hex>
func (h *AuthHandler) signedState() string {
	nonce := make([]byte, 16)
	rand.Read(nonce)
	nonceHex := hex.EncodeToString(nonce)

	mac := hmac.New(sha256.New, []byte(h.JWTSecret))
	mac.Write([]byte(nonceHex))
	sig := hex.EncodeToString(mac.Sum(nil))

	return nonceHex + "." + sig
}

// verifyState checks the HMAC signature on the state token.
func (h *AuthHandler) verifyState(state string) bool {
	if len(state) < 34 { // 32 hex chars + "." + at least 1 char
		return false
	}

	dot := 32 // nonce is always 32 hex chars (16 bytes)
	if dot >= len(state) || state[dot] != '.' {
		return false
	}

	nonceHex := state[:dot]
	sig := state[dot+1:]

	mac := hmac.New(sha256.New, []byte(h.JWTSecret))
	mac.Write([]byte(nonceHex))
	expectedSig := hex.EncodeToString(mac.Sum(nil))

	return hmac.Equal([]byte(sig), []byte(expectedSig))
}
