package handler

import (
	"database/sql"
	"encoding/json"
	"errors"
	"io"
	"log"
	"net/http"
	"os"

	"note-thing/backend/internal/billing"
	"note-thing/backend/internal/model"

	"github.com/stripe/stripe-go/v82"
	"github.com/stripe/stripe-go/v82/webhook"
)

type BillingHandler struct {
	DB      *sql.DB
	Billing *billing.Service
}

func (h *BillingHandler) GetPrice(w http.ResponseWriter, r *http.Request) {
	price, err := model.GetActivePrice(r.Context(), h.DB)
	if errors.Is(err, sql.ErrNoRows) {
		respondError(w, http.StatusNotFound, "no active price")
		return
	}
	if err != nil {
		respondError(w, http.StatusInternalServerError, "failed to get price")
		return
	}
	respondJSON(w, http.StatusOK, map[string]any{
		"amountCents": price.AmountCents,
		"currency":    price.Currency,
		"interval":    price.Interval,
		"trialDays":   14,
	})
}

func (h *BillingHandler) Status(w http.ResponseWriter, r *http.Request) {
	userID := userIDFromContext(r)
	sub, err := model.GetSubscriptionByUserID(r.Context(), h.DB, userID)
	if errors.Is(err, sql.ErrNoRows) {
		respondJSON(w, http.StatusOK, map[string]any{
			"subscription":    nil,
			"hasActiveAccess": false,
		})
		return
	}
	if err != nil {
		respondError(w, http.StatusInternalServerError, "failed to get subscription")
		return
	}

	// Get price info
	price, _ := model.GetActivePrice(r.Context(), h.DB)

	hasAccess := sub.Status == "active" || sub.Status == "trialing"
	respondJSON(w, http.StatusOK, map[string]any{
		"subscription": map[string]any{
			"id":                sub.ID,
			"status":            sub.Status,
			"stripePriceId":     sub.StripePriceID,
			"trialEnd":          sub.TrialEnd,
			"currentPeriodEnd":  sub.CurrentPeriodEnd,
			"cancelAtPeriodEnd": sub.CancelAtPeriodEnd,
			"amountCents":       price.AmountCents,
			"currency":          price.Currency,
		},
		"hasActiveAccess": hasAccess,
	})
}

func (h *BillingHandler) CreateCheckout(w http.ResponseWriter, r *http.Request) {
	userID := userIDFromContext(r)
	var input struct {
		SuccessURL string `json:"successUrl"`
		CancelURL  string `json:"cancelUrl"`
	}
	if err := decodeJSON(r, &input); err != nil {
		respondError(w, http.StatusBadRequest, "invalid request")
		return
	}

	user, err := model.GetUserByID(r.Context(), h.DB, userID)
	if err != nil {
		respondError(w, http.StatusInternalServerError, "failed to get user")
		return
	}

	customerID, err := h.Billing.GetOrCreateCustomer(r.Context(), userID, user.Email, user.Name)
	if err != nil {
		respondError(w, http.StatusInternalServerError, "failed to create customer")
		return
	}

	checkoutURL, err := h.Billing.CreateCheckoutSession(r.Context(), userID, customerID, input.SuccessURL, input.CancelURL)
	if err != nil {
		log.Printf("create checkout error: %v", err)
		respondError(w, http.StatusInternalServerError, "failed to create checkout")
		return
	}

	respondJSON(w, http.StatusOK, map[string]string{"checkoutUrl": checkoutURL})
}

func (h *BillingHandler) CreatePortal(w http.ResponseWriter, r *http.Request) {
	userID := userIDFromContext(r)
	var input struct {
		ReturnURL string `json:"returnUrl"`
	}
	if err := decodeJSON(r, &input); err != nil {
		respondError(w, http.StatusBadRequest, "invalid request")
		return
	}

	customerID, err := model.GetUserStripeCustomerID(r.Context(), h.DB, userID)
	if err != nil || customerID == "" {
		respondError(w, http.StatusBadRequest, "no billing account")
		return
	}

	portalURL, err := h.Billing.CreatePortalSession(r.Context(), customerID, input.ReturnURL)
	if err != nil {
		respondError(w, http.StatusInternalServerError, "failed to create portal")
		return
	}

	respondJSON(w, http.StatusOK, map[string]string{"portalUrl": portalURL})
}

func (h *BillingHandler) CancelSubscription(w http.ResponseWriter, r *http.Request) {
	userID := userIDFromContext(r)
	sub, err := model.GetSubscriptionByUserID(r.Context(), h.DB, userID)
	if err != nil {
		respondError(w, http.StatusNotFound, "no subscription")
		return
	}
	if err := h.Billing.CancelSubscription(r.Context(), sub.StripeSubscriptionID); err != nil {
		respondError(w, http.StatusInternalServerError, "failed to cancel")
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (h *BillingHandler) Reactivate(w http.ResponseWriter, r *http.Request) {
	userID := userIDFromContext(r)
	sub, err := model.GetSubscriptionByUserID(r.Context(), h.DB, userID)
	if err != nil {
		respondError(w, http.StatusNotFound, "no subscription")
		return
	}
	if err := h.Billing.ReactivateSubscription(r.Context(), sub.StripeSubscriptionID); err != nil {
		respondError(w, http.StatusInternalServerError, "failed to reactivate")
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (h *BillingHandler) HandleWebhook(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		respondError(w, http.StatusBadRequest, "read body failed")
		return
	}

	webhookSecret := os.Getenv("STRIPE_WEBHOOK_SECRET")
	event, err := webhook.ConstructEvent(body, r.Header.Get("Stripe-Signature"), webhookSecret)
	if err != nil {
		respondError(w, http.StatusBadRequest, "invalid signature")
		return
	}

	switch event.Type {
	case "checkout.session.completed":
		h.handleCheckoutCompleted(r, event)
	case "customer.subscription.created", "customer.subscription.updated":
		h.handleSubscriptionChange(r, event)
	case "customer.subscription.deleted":
		h.handleSubscriptionDeleted(r, event)
	}

	w.WriteHeader(http.StatusOK)
}

func (h *BillingHandler) handleCheckoutCompleted(r *http.Request, event stripe.Event) {
	var sess stripe.CheckoutSession
	if err := json.Unmarshal(event.Data.Raw, &sess); err != nil {
		log.Printf("webhook: unmarshal checkout session: %v", err)
		return
	}

	if sess.Subscription == nil {
		return
	}

	// The subscription webhook will handle the actual upsert
	log.Printf("webhook: checkout completed for subscription %s", sess.Subscription.ID)
}

func (h *BillingHandler) handleSubscriptionChange(r *http.Request, event stripe.Event) {
	var sub stripe.Subscription
	if err := json.Unmarshal(event.Data.Raw, &sub); err != nil {
		log.Printf("webhook: unmarshal subscription: %v", err)
		return
	}

	// Resolve user from customer
	user, err := model.GetUserByStripeCustomerID(r.Context(), h.DB, sub.Customer.ID)
	if errors.Is(err, sql.ErrNoRows) {
		// Customer might have metadata
		userID := sub.Metadata["user_id"]
		if userID == "" {
			log.Printf("webhook: unknown customer %s, no user_id metadata", sub.Customer.ID)
			return
		}
		if err := model.SetUserStripeCustomerID(r.Context(), h.DB, userID, sub.Customer.ID); err != nil {
			log.Printf("webhook: set customer id: %v", err)
			return
		}
		user.ID = userID
	} else if err != nil {
		log.Printf("webhook: get user by customer: %v", err)
		return
	}

	dbSub := billing.SubscriptionFromStripe(&sub, user.ID)
	if _, err := model.UpsertSubscription(r.Context(), h.DB, dbSub); err != nil {
		log.Printf("webhook: upsert subscription: %v", err)
	}
}

func (h *BillingHandler) handleSubscriptionDeleted(r *http.Request, event stripe.Event) {
	var sub stripe.Subscription
	if err := json.Unmarshal(event.Data.Raw, &sub); err != nil {
		log.Printf("webhook: unmarshal subscription: %v", err)
		return
	}

	user, err := model.GetUserByStripeCustomerID(r.Context(), h.DB, sub.Customer.ID)
	if err != nil {
		log.Printf("webhook: get user for deleted sub: %v", err)
		return
	}

	dbSub := billing.SubscriptionFromStripe(&sub, user.ID)
	dbSub.Status = "canceled"
	if _, err := model.UpsertSubscription(r.Context(), h.DB, dbSub); err != nil {
		log.Printf("webhook: upsert deleted subscription: %v", err)
	}
}
