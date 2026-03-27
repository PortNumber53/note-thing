package billing

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"note-thing/backend/internal/model"

	"github.com/stripe/stripe-go/v82"
	"github.com/stripe/stripe-go/v82/billingportal/session"
	checkoutsession "github.com/stripe/stripe-go/v82/checkout/session"
	"github.com/stripe/stripe-go/v82/customer"
	"github.com/stripe/stripe-go/v82/price"
	"github.com/stripe/stripe-go/v82/product"
	stripesub "github.com/stripe/stripe-go/v82/subscription"
)

const (
	appMetadataKey   = "app"
	appMetadataValue = "note-thing"
	defaultTrialDays = 14
)

type Service struct {
	DB *sql.DB
}

func NewService(db *sql.DB, stripeKey string) *Service {
	stripe.Key = stripeKey
	return &Service{DB: db}
}

// Bootstrap ensures a product and price exist in both Stripe and our DB.
func (s *Service) Bootstrap(ctx context.Context) error {
	// Check if we already have an active price locally
	_, err := model.GetActivePrice(ctx, s.DB)
	if err == nil {
		return nil // already bootstrapped
	}
	if !errors.Is(err, sql.ErrNoRows) {
		return fmt.Errorf("check active price: %w", err)
	}

	// Search Stripe for existing product
	params := &stripe.ProductSearchParams{}
	params.Query = fmt.Sprintf("metadata['%s']:'%s'", appMetadataKey, appMetadataValue)
	iter := product.Search(params)
	if iter.Next() {
		return s.syncProductFromStripe(ctx, iter.Product())
	}

	// Create new product + price in Stripe
	prodParams := &stripe.ProductParams{
		Name: stripe.String("Note Thing Subscription"),
		Metadata: map[string]string{
			appMetadataKey: appMetadataValue,
		},
	}
	prod, err := product.New(prodParams)
	if err != nil {
		return fmt.Errorf("create stripe product: %w", err)
	}

	priceParams := &stripe.PriceParams{
		Product:    stripe.String(prod.ID),
		UnitAmount: stripe.Int64(1099),
		Currency:   stripe.String("usd"),
		Recurring: &stripe.PriceRecurringParams{
			Interval: stripe.String("month"),
		},
	}
	p, err := price.New(priceParams)
	if err != nil {
		return fmt.Errorf("create stripe price: %w", err)
	}

	// Save to DB
	dbProd, err := model.CreateBillingProduct(ctx, s.DB, prod.ID, prod.Name)
	if err != nil {
		return fmt.Errorf("save product: %w", err)
	}
	_, err = model.CreateBillingPrice(ctx, s.DB, p.ID, dbProd.ID, int(p.UnitAmount), string(p.Currency), string(p.Recurring.Interval))
	if err != nil {
		return fmt.Errorf("save price: %w", err)
	}

	log.Printf("billing: bootstrapped product %s with price %s ($%.2f/month)", prod.ID, p.ID, float64(p.UnitAmount)/100)
	return nil
}

func (s *Service) syncProductFromStripe(ctx context.Context, prod *stripe.Product) error {
	dbProd, err := model.CreateBillingProduct(ctx, s.DB, prod.ID, prod.Name)
	if err != nil {
		return fmt.Errorf("sync product: %w", err)
	}

	// Find active price for this product
	params := &stripe.PriceListParams{
		Active:  stripe.Bool(true),
		Product: stripe.String(prod.ID),
	}
	iter := price.List(params)
	for iter.Next() {
		p := iter.Price()
		_, err := model.CreateBillingPrice(ctx, s.DB, p.ID, dbProd.ID, int(p.UnitAmount), string(p.Currency), string(p.Recurring.Interval))
		if err != nil {
			return fmt.Errorf("sync price %s: %w", p.ID, err)
		}
	}
	return iter.Err()
}

// SyncFromStripe fetches all active subscriptions from Stripe and upserts them.
func (s *Service) SyncFromStripe(ctx context.Context) error {
	params := &stripe.SubscriptionListParams{}
	params.Filters.AddFilter("status", "", "all")
	iter := stripesub.List(params)

	synced := 0
	for iter.Next() {
		sub := iter.Subscription()
		user, err := model.GetUserByStripeCustomerID(ctx, s.DB, sub.Customer.ID)
		if errors.Is(err, sql.ErrNoRows) {
			continue // unknown customer, skip
		}
		if err != nil {
			log.Printf("billing sync: error looking up customer %s: %v", sub.Customer.ID, err)
			continue
		}

		if _, err := model.UpsertSubscription(ctx, s.DB, subscriptionFromStripe(sub, user.ID)); err != nil {
			log.Printf("billing sync: error upserting sub %s: %v", sub.ID, err)
			continue
		}
		synced++
	}
	if err := iter.Err(); err != nil {
		return fmt.Errorf("list subscriptions: %w", err)
	}
	if synced > 0 {
		log.Printf("billing: synced %d subscriptions from Stripe", synced)
	}
	return nil
}

// GetOrCreateCustomer ensures a Stripe customer exists for the user.
func (s *Service) GetOrCreateCustomer(ctx context.Context, userID, email, name string) (string, error) {
	existing, err := model.GetUserStripeCustomerID(ctx, s.DB, userID)
	if err != nil {
		return "", err
	}
	if existing != "" {
		return existing, nil
	}

	params := &stripe.CustomerParams{
		Email: stripe.String(email),
		Name:  stripe.String(name),
		Metadata: map[string]string{
			"user_id": userID,
		},
	}
	cust, err := customer.New(params)
	if err != nil {
		return "", fmt.Errorf("create stripe customer: %w", err)
	}

	if err := model.SetUserStripeCustomerID(ctx, s.DB, userID, cust.ID); err != nil {
		return "", fmt.Errorf("save customer id: %w", err)
	}

	return cust.ID, nil
}

// CreateCheckoutSession creates a Stripe Checkout session for subscription.
func (s *Service) CreateCheckoutSession(ctx context.Context, userID, customerID, successURL, cancelURL string) (string, error) {
	activePrice, err := model.GetActivePrice(ctx, s.DB)
	if err != nil {
		return "", fmt.Errorf("get active price: %w", err)
	}

	params := &stripe.CheckoutSessionParams{
		Customer: stripe.String(customerID),
		Mode:     stripe.String(string(stripe.CheckoutSessionModeSubscription)),
		LineItems: []*stripe.CheckoutSessionLineItemParams{
			{
				Price:    stripe.String(activePrice.StripePriceID),
				Quantity: stripe.Int64(1),
			},
		},
		SuccessURL: stripe.String(successURL),
		CancelURL:  stripe.String(cancelURL),
		SubscriptionData: &stripe.CheckoutSessionSubscriptionDataParams{
			Metadata: map[string]string{
				"user_id": userID,
			},
		},
	}

	// Only offer trial if user has never had a subscription
	hadSub, err := model.UserHasHadSubscription(ctx, s.DB, userID)
	if err != nil {
		return "", err
	}
	if !hadSub {
		params.SubscriptionData.TrialPeriodDays = stripe.Int64(defaultTrialDays)
	}

	sess, err := checkoutsession.New(params)
	if err != nil {
		return "", fmt.Errorf("create checkout session: %w", err)
	}
	return sess.URL, nil
}

// CreatePortalSession creates a Stripe customer portal session.
func (s *Service) CreatePortalSession(ctx context.Context, customerID, returnURL string) (string, error) {
	params := &stripe.BillingPortalSessionParams{
		Customer:  stripe.String(customerID),
		ReturnURL: stripe.String(returnURL),
	}
	sess, err := session.New(params)
	if err != nil {
		return "", fmt.Errorf("create portal session: %w", err)
	}
	return sess.URL, nil
}

// CancelSubscription sets cancel_at_period_end on the Stripe subscription.
func (s *Service) CancelSubscription(ctx context.Context, stripeSubID string) error {
	params := &stripe.SubscriptionParams{
		CancelAtPeriodEnd: stripe.Bool(true),
	}
	_, err := stripesub.Update(stripeSubID, params)
	return err
}

// ReactivateSubscription undoes cancel_at_period_end.
func (s *Service) ReactivateSubscription(ctx context.Context, stripeSubID string) error {
	params := &stripe.SubscriptionParams{
		CancelAtPeriodEnd: stripe.Bool(false),
	}
	_, err := stripesub.Update(stripeSubID, params)
	return err
}

// ChangePrice creates a new Stripe product+price, then queues a migration.
func (s *Service) ChangePrice(ctx context.Context, newAmountCents int) (model.BillingPriceMigration, error) {
	oldPrice, err := model.GetActivePrice(ctx, s.DB)
	if err != nil {
		return model.BillingPriceMigration{}, fmt.Errorf("get current price: %w", err)
	}

	// Create new Stripe product + price
	prodParams := &stripe.ProductParams{
		Name: stripe.String(fmt.Sprintf("Note Thing Subscription (%s)", time.Now().Format("2006-01-02"))),
		Metadata: map[string]string{
			appMetadataKey: appMetadataValue,
		},
	}
	prod, err := product.New(prodParams)
	if err != nil {
		return model.BillingPriceMigration{}, fmt.Errorf("create stripe product: %w", err)
	}

	priceParams := &stripe.PriceParams{
		Product:    stripe.String(prod.ID),
		UnitAmount: stripe.Int64(int64(newAmountCents)),
		Currency:   stripe.String("usd"),
		Recurring: &stripe.PriceRecurringParams{
			Interval: stripe.String("month"),
		},
	}
	p, err := price.New(priceParams)
	if err != nil {
		return model.BillingPriceMigration{}, fmt.Errorf("create stripe price: %w", err)
	}

	// Save to DB
	dbProd, err := model.CreateBillingProduct(ctx, s.DB, prod.ID, prod.Name)
	if err != nil {
		return model.BillingPriceMigration{}, fmt.Errorf("save product: %w", err)
	}
	newPrice, err := model.CreateBillingPrice(ctx, s.DB, p.ID, dbProd.ID, int(p.UnitAmount), string(p.Currency), string(p.Recurring.Interval))
	if err != nil {
		return model.BillingPriceMigration{}, fmt.Errorf("save price: %w", err)
	}

	// Count subs to migrate
	subs, err := model.ListSubscriptionsByStripePriceID(ctx, s.DB, oldPrice.StripePriceID, 100000)
	if err != nil {
		return model.BillingPriceMigration{}, err
	}

	graceDays := 0
	if g := os.Getenv("STRIPE_PRICE_MIGRATION_GRACE_DAYS"); g != "" {
		if n, err := strconv.Atoi(g); err == nil && n >= 0 {
			graceDays = n
		}
	}

	migrateAfter := time.Now().Add(time.Duration(graceDays) * 24 * time.Hour)
	migration, err := model.CreatePriceMigration(ctx, s.DB, oldPrice.ID, newPrice.ID, graceDays, migrateAfter, len(subs))
	if err != nil {
		return model.BillingPriceMigration{}, fmt.Errorf("create migration: %w", err)
	}

	log.Printf("billing: price change queued: $%.2f -> $%.2f, %d subs to migrate (grace: %d days)",
		float64(oldPrice.AmountCents)/100, float64(newAmountCents)/100, len(subs), graceDays)

	return migration, nil
}

func SubscriptionFromStripe(sub *stripe.Subscription, userID string) model.BillingSubscription {
	return subscriptionFromStripe(sub, userID)
}

func subscriptionFromStripe(sub *stripe.Subscription, userID string) model.BillingSubscription {
	s := model.BillingSubscription{
		UserID:               userID,
		StripeSubscriptionID: sub.ID,
		StripeCustomerID:     sub.Customer.ID,
		Status:               string(sub.Status),
		CancelAtPeriodEnd:    sub.CancelAtPeriodEnd,
	}

	if len(sub.Items.Data) > 0 {
		s.StripePriceID = sub.Items.Data[0].Price.ID
	}

	if sub.TrialStart > 0 {
		t := time.Unix(sub.TrialStart, 0)
		s.TrialStart = &t
	}
	if sub.TrialEnd > 0 {
		t := time.Unix(sub.TrialEnd, 0)
		s.TrialEnd = &t
	}
	// CurrentPeriodStart/End are available via the latest invoice in v82+
	// We populate them from webhook invoice data when available
	if sub.CanceledAt > 0 {
		t := time.Unix(sub.CanceledAt, 0)
		s.CanceledAt = &t
	}

	return s
}
