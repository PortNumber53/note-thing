package billing

import (
	"context"
	"database/sql"
	"errors"
	"log"
	"time"

	"note-thing/backend/internal/model"

	"github.com/stripe/stripe-go/v82"
	stripesub "github.com/stripe/stripe-go/v82/subscription"
)

type PriceMigrationWorker struct {
	db     *sql.DB
	stopCh chan struct{}
}

func NewPriceMigrationWorker(db *sql.DB) *PriceMigrationWorker {
	return &PriceMigrationWorker{db: db, stopCh: make(chan struct{})}
}

func (w *PriceMigrationWorker) Start(ctx context.Context) {
	ticker := time.NewTicker(30 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-w.stopCh:
			return
		case <-ctx.Done():
			return
		case <-ticker.C:
			w.tick(ctx)
		}
	}
}

func (w *PriceMigrationWorker) Stop() {
	close(w.stopCh)
}

func (w *PriceMigrationWorker) tick(ctx context.Context) {
	migration, err := model.GetPendingPriceMigration(ctx, w.db)
	if errors.Is(err, sql.ErrNoRows) {
		return
	}
	if err != nil {
		log.Printf("billing worker: error getting migration: %v", err)
		return
	}

	if migration.Status == "pending" {
		if err := model.UpdatePriceMigrationStatus(ctx, w.db, migration.ID, "in_progress"); err != nil {
			log.Printf("billing worker: error starting migration %s: %v", migration.ID, err)
			return
		}
		migration.Status = "in_progress"
	}

	// Get the old price's Stripe ID
	oldPrice, err := model.GetActivePrice(ctx, w.db)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		log.Printf("billing worker: error getting old price: %v", err)
		return
	}

	// Get the old price by ID from migration
	var oldStripePriceID string
	row := w.db.QueryRowContext(ctx, `SELECT stripe_price_id FROM billing_prices WHERE id = $1`, migration.OldPriceID)
	if err := row.Scan(&oldStripePriceID); err != nil {
		log.Printf("billing worker: error getting old stripe price: %v", err)
		return
	}

	var newStripePriceID string
	row = w.db.QueryRowContext(ctx, `SELECT stripe_price_id FROM billing_prices WHERE id = $1`, migration.NewPriceID)
	if err := row.Scan(&newStripePriceID); err != nil {
		log.Printf("billing worker: error getting new stripe price: %v", err)
		return
	}

	// Fetch batch of subscriptions on old price
	subs, err := model.ListSubscriptionsByStripePriceID(ctx, w.db, oldStripePriceID, 25)
	if err != nil {
		log.Printf("billing worker: error listing subs: %v", err)
		return
	}

	if len(subs) == 0 {
		// All done — deactivate old price
		if err := model.DeactivatePrice(ctx, w.db, migration.OldPriceID); err != nil {
			log.Printf("billing worker: error deactivating old price: %v", err)
		}
		// Archive in Stripe
		_, _ = oldPrice, oldStripePriceID // price already fetched
		if err := model.CompletePriceMigration(ctx, w.db, migration.ID, "completed", nil); err != nil {
			log.Printf("billing worker: error completing migration: %v", err)
		}
		log.Printf("billing worker: migration %s completed", migration.ID)
		return
	}

	migrated := migration.MigratedSubs
	failed := migration.FailedSubs

	for _, sub := range subs {
		// Get the subscription's current item ID from Stripe
		stripeSub, err := stripesub.Get(sub.StripeSubscriptionID, nil)
		if err != nil {
			log.Printf("billing worker: error fetching sub %s: %v", sub.StripeSubscriptionID, err)
			failed++
			continue
		}

		if len(stripeSub.Items.Data) == 0 {
			failed++
			continue
		}

		itemID := stripeSub.Items.Data[0].ID

		// Update the subscription's price
		params := &stripe.SubscriptionParams{
			Items: []*stripe.SubscriptionItemsParams{
				{
					ID:    stripe.String(itemID),
					Price: stripe.String(newStripePriceID),
				},
			},
			ProrationBehavior: stripe.String("none"),
		}
		_, err = stripesub.Update(sub.StripeSubscriptionID, params)
		if err != nil {
			log.Printf("billing worker: error migrating sub %s: %v", sub.StripeSubscriptionID, err)
			failed++
			continue
		}

		// Update local DB
		sub.StripePriceID = newStripePriceID
		if _, err := model.UpsertSubscription(ctx, w.db, sub); err != nil {
			log.Printf("billing worker: error updating local sub %s: %v", sub.ID, err)
		}
		migrated++
	}

	if err := model.UpdatePriceMigrationProgress(ctx, w.db, migration.ID, migrated, failed); err != nil {
		log.Printf("billing worker: error updating progress: %v", err)
	}
}
