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
	"golang.org/x/crypto/bcrypt"
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

func (h *AuthHandler) GoogleLogin(w http.ResponseWriter, r *http.Request) {
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

// TokenExchange handles mobile auth: accepts a Google ID token, verifies it,
// and returns a JWT. Used by native mobile apps that get an idToken from Google Sign-In.
func (h *AuthHandler) TokenExchange(w http.ResponseWriter, r *http.Request) {
	var input struct {
		IDToken     string `json:"idToken"`
		AccessToken string `json:"accessToken"`
	}
	if err := decodeJSON(r, &input); err != nil {
		respondError(w, http.StatusBadRequest, "invalid request")
		return
	}

	var userInfoURL string
	if input.AccessToken != "" {
		// Use access token to fetch user info directly
		userInfoURL = "https://www.googleapis.com/oauth2/v2/userinfo?access_token=" + input.AccessToken
	} else if input.IDToken != "" {
		// Verify ID token via tokeninfo endpoint
		userInfoURL = "https://oauth2.googleapis.com/tokeninfo?id_token=" + input.IDToken
	} else {
		respondError(w, http.StatusBadRequest, "idToken or accessToken is required")
		return
	}

	resp, err := http.Get(userInfoURL)
	if err != nil {
		log.Printf("token verify failed: %v", err)
		respondError(w, http.StatusBadRequest, "failed to verify token")
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		respondError(w, http.StatusUnauthorized, "invalid token")
		return
	}

	var tokenInfo struct {
		Sub     string `json:"sub"`
		ID      string `json:"id"`      // userinfo v2 uses "id" instead of "sub"
		Email   string `json:"email"`
		Name    string `json:"name"`
		Picture string `json:"picture"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&tokenInfo); err != nil {
		respondError(w, http.StatusInternalServerError, "failed to decode token info")
		return
	}

	// userinfo v2 uses "id", tokeninfo uses "sub"
	googleID := tokenInfo.Sub
	if googleID == "" {
		googleID = tokenInfo.ID
	}
	if googleID == "" {
		respondError(w, http.StatusBadRequest, "could not determine user identity")
		return
	}

	// Upsert user
	user, err := model.UpsertUser(r.Context(), h.DB, googleID, tokenInfo.Email, tokenInfo.Name, tokenInfo.Picture)
	if err != nil {
		log.Printf("upsert user failed: %v", err)
		respondError(w, http.StatusInternalServerError, "failed to create user")
		return
	}

	// Create default notebook if needed
	var nbCount int
	err = h.DB.QueryRowContext(r.Context(), `SELECT COUNT(*) FROM notebooks WHERE user_id = $1`, user.ID).Scan(&nbCount)
	if err != nil {
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

	respondJSON(w, http.StatusOK, map[string]any{
		"token": signed,
		"user":  user,
	})
}

func (h *AuthHandler) Signup(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Email    string `json:"email"`
		Password string `json:"password"`
		Name     string `json:"name"`
	}
	if err := decodeJSON(r, &input); err != nil {
		respondError(w, http.StatusBadRequest, "invalid request")
		return
	}
	input.Email = strings.TrimSpace(strings.ToLower(input.Email))
	input.Name = strings.TrimSpace(input.Name)
	if input.Email == "" || input.Password == "" || input.Name == "" {
		respondError(w, http.StatusBadRequest, "email, password, and name are required")
		return
	}
	if len(input.Password) < 8 {
		respondError(w, http.StatusBadRequest, "password must be at least 8 characters")
		return
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		respondError(w, http.StatusInternalServerError, "internal error")
		return
	}

	user, err := model.CreateUserWithPassword(r.Context(), h.DB, input.Email, input.Name, string(hash))
	if err != nil {
		if strings.Contains(err.Error(), "duplicate key") || strings.Contains(err.Error(), "unique") {
			respondError(w, http.StatusConflict, "email already registered")
			return
		}
		log.Printf("create user failed: %v", err)
		respondError(w, http.StatusInternalServerError, "failed to create account")
		return
	}

	// Create default notebook
	tx, err := h.DB.BeginTx(r.Context(), nil)
	if err != nil {
		log.Printf("begin tx for default notebook failed: %v", err)
		respondError(w, http.StatusInternalServerError, "failed to create account")
		return
	}
	defer tx.Rollback()
	if _, err := model.CreateDefaultNotebook(r.Context(), tx, user.ID); err != nil {
		log.Printf("create default notebook failed: %v", err)
		respondError(w, http.StatusInternalServerError, "failed to create account")
		return
	}
	if err := tx.Commit(); err != nil {
		log.Printf("commit default notebook failed: %v", err)
		respondError(w, http.StatusInternalServerError, "failed to create account")
		return
	}

	signed := h.signJWT(user)
	respondJSON(w, http.StatusCreated, map[string]any{"token": signed, "user": user})
}

func (h *AuthHandler) EmailLogin(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	if err := decodeJSON(r, &input); err != nil {
		respondError(w, http.StatusBadRequest, "invalid request")
		return
	}
	input.Email = strings.TrimSpace(strings.ToLower(input.Email))
	if input.Email == "" || input.Password == "" {
		respondError(w, http.StatusBadRequest, "email and password are required")
		return
	}

	user, passwordHash, err := model.GetUserByEmail(r.Context(), h.DB, input.Email)
	if err != nil {
		respondError(w, http.StatusUnauthorized, "invalid email or password")
		return
	}
	if passwordHash == "" {
		respondError(w, http.StatusUnauthorized, "this account uses Google sign-in")
		return
	}
	if err := bcrypt.CompareHashAndPassword([]byte(passwordHash), []byte(input.Password)); err != nil {
		respondError(w, http.StatusUnauthorized, "invalid email or password")
		return
	}

	signed := h.signJWT(user)
	respondJSON(w, http.StatusOK, map[string]any{"token": signed, "user": user})
}

func (h *AuthHandler) signJWT(user model.User) string {
	jwtToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub":   user.ID,
		"email": user.Email,
		"iat":   time.Now().Unix(),
		"exp":   time.Now().Add(24 * time.Hour).Unix(),
	})
	signed, _ := jwtToken.SignedString([]byte(h.JWTSecret))
	return signed
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
