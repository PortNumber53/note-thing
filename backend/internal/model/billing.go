package model

import (
	"context"
	"database/sql"
	"time"
)

type BillingProduct struct {
	ID              string    `json:"id"`
	StripeProductID string    `json:"stripeProductId"`
	Name            string    `json:"name"`
	Active          bool      `json:"active"`
	CreatedAt       time.Time `json:"createdAt"`
}

type BillingPrice struct {
	ID            string    `json:"id"`
	StripePriceID string    `json:"stripePriceId"`
	ProductID     string    `json:"productId"`
	AmountCents   int       `json:"amountCents"`
	Currency      string    `json:"currency"`
	Interval      string    `json:"interval"`
	Active        bool      `json:"active"`
	CreatedAt     time.Time `json:"createdAt"`
}

type BillingSubscription struct {
	ID                   string     `json:"id"`
	UserID               string     `json:"userId"`
	StripeSubscriptionID string     `json:"stripeSubscriptionId"`
	StripeCustomerID     string     `json:"stripeCustomerId"`
	StripePriceID        string     `json:"stripePriceId"`
	Status               string     `json:"status"`
	TrialStart           *time.Time `json:"trialStart"`
	TrialEnd             *time.Time `json:"trialEnd"`
	CurrentPeriodStart   *time.Time `json:"currentPeriodStart"`
	CurrentPeriodEnd     *time.Time `json:"currentPeriodEnd"`
	CancelAtPeriodEnd    bool       `json:"cancelAtPeriodEnd"`
	CanceledAt           *time.Time `json:"canceledAt"`
	CreatedAt            time.Time  `json:"createdAt"`
	UpdatedAt            time.Time  `json:"updatedAt"`
}

type BillingPriceMigration struct {
	ID              string     `json:"id"`
	OldPriceID      string     `json:"oldPriceId"`
	NewPriceID      string     `json:"newPriceId"`
	Status          string     `json:"status"`
	GracePeriodDays int        `json:"gracePeriodDays"`
	MigrateAfter    time.Time  `json:"migrateAfter"`
	TotalSubs       int        `json:"totalSubs"`
	MigratedSubs    int        `json:"migratedSubs"`
	FailedSubs      int        `json:"failedSubs"`
	ErrorMessage    *string    `json:"errorMessage"`
	StartedAt       *time.Time `json:"startedAt"`
	CompletedAt     *time.Time `json:"completedAt"`
	CreatedAt       time.Time  `json:"createdAt"`
}

// Products

func CreateBillingProduct(ctx context.Context, db *sql.DB, stripeProductID, name string) (BillingProduct, error) {
	var p BillingProduct
	err := db.QueryRowContext(ctx, `
		INSERT INTO billing_products (stripe_product_id, name)
		VALUES ($1, $2)
		ON CONFLICT (stripe_product_id) DO UPDATE SET name = EXCLUDED.name, active = true
		RETURNING id, stripe_product_id, name, active, created_at
	`, stripeProductID, name).Scan(&p.ID, &p.StripeProductID, &p.Name, &p.Active, &p.CreatedAt)
	return p, err
}

func GetActiveProduct(ctx context.Context, db *sql.DB) (BillingProduct, error) {
	var p BillingProduct
	err := db.QueryRowContext(ctx, `
		SELECT id, stripe_product_id, name, active, created_at
		FROM billing_products WHERE active = true
		ORDER BY created_at DESC LIMIT 1
	`).Scan(&p.ID, &p.StripeProductID, &p.Name, &p.Active, &p.CreatedAt)
	return p, err
}

// Prices

func CreateBillingPrice(ctx context.Context, db *sql.DB, stripePriceID, productID string, amountCents int, currency, interval string) (BillingPrice, error) {
	var p BillingPrice
	err := db.QueryRowContext(ctx, `
		INSERT INTO billing_prices (stripe_price_id, product_id, amount_cents, currency, interval)
		VALUES ($1, $2, $3, $4, $5)
		ON CONFLICT (stripe_price_id) DO UPDATE SET
			amount_cents = EXCLUDED.amount_cents,
			active = true
		RETURNING id, stripe_price_id, product_id, amount_cents, currency, interval, active, created_at
	`, stripePriceID, productID, amountCents, currency, interval).Scan(
		&p.ID, &p.StripePriceID, &p.ProductID, &p.AmountCents, &p.Currency, &p.Interval, &p.Active, &p.CreatedAt,
	)
	return p, err
}

func GetActivePrice(ctx context.Context, db *sql.DB) (BillingPrice, error) {
	var p BillingPrice
	err := db.QueryRowContext(ctx, `
		SELECT id, stripe_price_id, product_id, amount_cents, currency, interval, active, created_at
		FROM billing_prices WHERE active = true
		ORDER BY created_at DESC LIMIT 1
	`).Scan(&p.ID, &p.StripePriceID, &p.ProductID, &p.AmountCents, &p.Currency, &p.Interval, &p.Active, &p.CreatedAt)
	return p, err
}

func DeactivatePrice(ctx context.Context, db *sql.DB, priceID string) error {
	_, err := db.ExecContext(ctx, `UPDATE billing_prices SET active = false WHERE id = $1`, priceID)
	return err
}

// Subscriptions

func UpsertSubscription(ctx context.Context, db *sql.DB, sub BillingSubscription) (BillingSubscription, error) {
	var s BillingSubscription
	err := db.QueryRowContext(ctx, `
		INSERT INTO billing_subscriptions (
			user_id, stripe_subscription_id, stripe_customer_id, stripe_price_id,
			status, trial_start, trial_end, current_period_start, current_period_end,
			cancel_at_period_end, canceled_at
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
		ON CONFLICT (stripe_subscription_id) DO UPDATE SET
			stripe_price_id = EXCLUDED.stripe_price_id,
			status = EXCLUDED.status,
			trial_start = EXCLUDED.trial_start,
			trial_end = EXCLUDED.trial_end,
			current_period_start = EXCLUDED.current_period_start,
			current_period_end = EXCLUDED.current_period_end,
			cancel_at_period_end = EXCLUDED.cancel_at_period_end,
			canceled_at = EXCLUDED.canceled_at,
			updated_at = now()
		RETURNING id, user_id, stripe_subscription_id, stripe_customer_id, stripe_price_id,
			status, trial_start, trial_end, current_period_start, current_period_end,
			cancel_at_period_end, canceled_at, created_at, updated_at
	`, sub.UserID, sub.StripeSubscriptionID, sub.StripeCustomerID, sub.StripePriceID,
		sub.Status, sub.TrialStart, sub.TrialEnd, sub.CurrentPeriodStart, sub.CurrentPeriodEnd,
		sub.CancelAtPeriodEnd, sub.CanceledAt,
	).Scan(
		&s.ID, &s.UserID, &s.StripeSubscriptionID, &s.StripeCustomerID, &s.StripePriceID,
		&s.Status, &s.TrialStart, &s.TrialEnd, &s.CurrentPeriodStart, &s.CurrentPeriodEnd,
		&s.CancelAtPeriodEnd, &s.CanceledAt, &s.CreatedAt, &s.UpdatedAt,
	)
	return s, err
}

func GetSubscriptionByUserID(ctx context.Context, db *sql.DB, userID string) (BillingSubscription, error) {
	var s BillingSubscription
	err := db.QueryRowContext(ctx, `
		SELECT id, user_id, stripe_subscription_id, stripe_customer_id, stripe_price_id,
			status, trial_start, trial_end, current_period_start, current_period_end,
			cancel_at_period_end, canceled_at, created_at, updated_at
		FROM billing_subscriptions WHERE user_id = $1
	`, userID).Scan(
		&s.ID, &s.UserID, &s.StripeSubscriptionID, &s.StripeCustomerID, &s.StripePriceID,
		&s.Status, &s.TrialStart, &s.TrialEnd, &s.CurrentPeriodStart, &s.CurrentPeriodEnd,
		&s.CancelAtPeriodEnd, &s.CanceledAt, &s.CreatedAt, &s.UpdatedAt,
	)
	return s, err
}

func GetSubscriptionByStripeID(ctx context.Context, db *sql.DB, stripeSubID string) (BillingSubscription, error) {
	var s BillingSubscription
	err := db.QueryRowContext(ctx, `
		SELECT id, user_id, stripe_subscription_id, stripe_customer_id, stripe_price_id,
			status, trial_start, trial_end, current_period_start, current_period_end,
			cancel_at_period_end, canceled_at, created_at, updated_at
		FROM billing_subscriptions WHERE stripe_subscription_id = $1
	`, stripeSubID).Scan(
		&s.ID, &s.UserID, &s.StripeSubscriptionID, &s.StripeCustomerID, &s.StripePriceID,
		&s.Status, &s.TrialStart, &s.TrialEnd, &s.CurrentPeriodStart, &s.CurrentPeriodEnd,
		&s.CancelAtPeriodEnd, &s.CanceledAt, &s.CreatedAt, &s.UpdatedAt,
	)
	return s, err
}

func ListSubscriptionsByStripePriceID(ctx context.Context, db *sql.DB, stripePriceID string, limit int) ([]BillingSubscription, error) {
	rows, err := db.QueryContext(ctx, `
		SELECT id, user_id, stripe_subscription_id, stripe_customer_id, stripe_price_id,
			status, trial_start, trial_end, current_period_start, current_period_end,
			cancel_at_period_end, canceled_at, created_at, updated_at
		FROM billing_subscriptions
		WHERE stripe_price_id = $1 AND status IN ('active', 'trialing')
		LIMIT $2
	`, stripePriceID, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var subs []BillingSubscription
	for rows.Next() {
		var s BillingSubscription
		if err := rows.Scan(
			&s.ID, &s.UserID, &s.StripeSubscriptionID, &s.StripeCustomerID, &s.StripePriceID,
			&s.Status, &s.TrialStart, &s.TrialEnd, &s.CurrentPeriodStart, &s.CurrentPeriodEnd,
			&s.CancelAtPeriodEnd, &s.CanceledAt, &s.CreatedAt, &s.UpdatedAt,
		); err != nil {
			return nil, err
		}
		subs = append(subs, s)
	}
	return subs, rows.Err()
}

func UserHasActiveSubscription(ctx context.Context, db *sql.DB, userID string) (bool, error) {
	var exists bool
	err := db.QueryRowContext(ctx, `
		SELECT EXISTS(
			SELECT 1 FROM billing_subscriptions
			WHERE user_id = $1 AND status IN ('active', 'trialing')
		)
	`, userID).Scan(&exists)
	return exists, err
}

func UserHasHadSubscription(ctx context.Context, db *sql.DB, userID string) (bool, error) {
	var exists bool
	err := db.QueryRowContext(ctx, `
		SELECT EXISTS(SELECT 1 FROM billing_subscriptions WHERE user_id = $1)
	`, userID).Scan(&exists)
	return exists, err
}

// Stripe customer on users

func SetUserStripeCustomerID(ctx context.Context, db *sql.DB, userID, customerID string) error {
	_, err := db.ExecContext(ctx, `
		UPDATE users SET stripe_customer_id = $1, updated_at = now() WHERE id = $2
	`, customerID, userID)
	return err
}

func GetUserStripeCustomerID(ctx context.Context, db *sql.DB, userID string) (string, error) {
	var id sql.NullString
	err := db.QueryRowContext(ctx, `SELECT stripe_customer_id FROM users WHERE id = $1`, userID).Scan(&id)
	if err != nil {
		return "", err
	}
	return id.String, nil
}

func GetUserByStripeCustomerID(ctx context.Context, db *sql.DB, customerID string) (User, error) {
	var u User
	err := db.QueryRowContext(ctx, `
		SELECT id, google_id, email, name, avatar_url, created_at, updated_at
		FROM users WHERE stripe_customer_id = $1
	`, customerID).Scan(
		&u.ID, &u.GoogleID, &u.Email, &u.Name, &u.AvatarURL, &u.CreatedAt, &u.UpdatedAt,
	)
	return u, err
}

// Price migrations

func CreatePriceMigration(ctx context.Context, db *sql.DB, oldPriceID, newPriceID string, graceDays int, migrateAfter time.Time, totalSubs int) (BillingPriceMigration, error) {
	var m BillingPriceMigration
	err := db.QueryRowContext(ctx, `
		INSERT INTO billing_price_migrations (old_price_id, new_price_id, grace_period_days, migrate_after, total_subs)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id, old_price_id, new_price_id, status, grace_period_days, migrate_after,
			total_subs, migrated_subs, failed_subs, error_message, started_at, completed_at, created_at
	`, oldPriceID, newPriceID, graceDays, migrateAfter, totalSubs).Scan(
		&m.ID, &m.OldPriceID, &m.NewPriceID, &m.Status, &m.GracePeriodDays, &m.MigrateAfter,
		&m.TotalSubs, &m.MigratedSubs, &m.FailedSubs, &m.ErrorMessage, &m.StartedAt, &m.CompletedAt, &m.CreatedAt,
	)
	return m, err
}

func GetPendingPriceMigration(ctx context.Context, db *sql.DB) (BillingPriceMigration, error) {
	var m BillingPriceMigration
	err := db.QueryRowContext(ctx, `
		SELECT id, old_price_id, new_price_id, status, grace_period_days, migrate_after,
			total_subs, migrated_subs, failed_subs, error_message, started_at, completed_at, created_at
		FROM billing_price_migrations
		WHERE status IN ('pending', 'in_progress') AND migrate_after <= now()
		ORDER BY created_at ASC LIMIT 1
	`).Scan(
		&m.ID, &m.OldPriceID, &m.NewPriceID, &m.Status, &m.GracePeriodDays, &m.MigrateAfter,
		&m.TotalSubs, &m.MigratedSubs, &m.FailedSubs, &m.ErrorMessage, &m.StartedAt, &m.CompletedAt, &m.CreatedAt,
	)
	return m, err
}

func UpdatePriceMigrationStatus(ctx context.Context, db *sql.DB, id, status string) error {
	_, err := db.ExecContext(ctx, `
		UPDATE billing_price_migrations SET status = $1, started_at = COALESCE(started_at, now()), updated_at = now()
		WHERE id = $2
	`, status, id)
	return err
}

func UpdatePriceMigrationProgress(ctx context.Context, db *sql.DB, id string, migrated, failed int) error {
	_, err := db.ExecContext(ctx, `
		UPDATE billing_price_migrations SET migrated_subs = $1, failed_subs = $2, updated_at = now()
		WHERE id = $3
	`, migrated, failed, id)
	return err
}

func CompletePriceMigration(ctx context.Context, db *sql.DB, id, status string, errMsg *string) error {
	_, err := db.ExecContext(ctx, `
		UPDATE billing_price_migrations SET status = $1, error_message = $2, completed_at = now(), updated_at = now()
		WHERE id = $3
	`, status, errMsg, id)
	return err
}

func GetLatestPriceMigration(ctx context.Context, db *sql.DB) (BillingPriceMigration, error) {
	var m BillingPriceMigration
	err := db.QueryRowContext(ctx, `
		SELECT id, old_price_id, new_price_id, status, grace_period_days, migrate_after,
			total_subs, migrated_subs, failed_subs, error_message, started_at, completed_at, created_at
		FROM billing_price_migrations
		ORDER BY created_at DESC LIMIT 1
	`).Scan(
		&m.ID, &m.OldPriceID, &m.NewPriceID, &m.Status, &m.GracePeriodDays, &m.MigrateAfter,
		&m.TotalSubs, &m.MigratedSubs, &m.FailedSubs, &m.ErrorMessage, &m.StartedAt, &m.CompletedAt, &m.CreatedAt,
	)
	return m, err
}
