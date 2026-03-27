CREATE TABLE billing_products (
    id                UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    stripe_product_id TEXT NOT NULL UNIQUE,
    name              TEXT NOT NULL,
    active            BOOLEAN NOT NULL DEFAULT true,
    created_at        TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE billing_prices (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    stripe_price_id TEXT NOT NULL UNIQUE,
    product_id      UUID NOT NULL REFERENCES billing_products(id),
    amount_cents    INTEGER NOT NULL,
    currency        TEXT NOT NULL DEFAULT 'usd',
    interval        TEXT NOT NULL DEFAULT 'month',
    active          BOOLEAN NOT NULL DEFAULT true,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX idx_billing_prices_active ON billing_prices(active) WHERE active = true;

CREATE TABLE billing_subscriptions (
    id                     UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id                UUID NOT NULL UNIQUE REFERENCES users(id) ON DELETE CASCADE,
    stripe_subscription_id TEXT NOT NULL UNIQUE,
    stripe_customer_id     TEXT NOT NULL,
    stripe_price_id        TEXT NOT NULL,
    status                 TEXT NOT NULL,
    trial_start            TIMESTAMPTZ,
    trial_end              TIMESTAMPTZ,
    current_period_start   TIMESTAMPTZ,
    current_period_end     TIMESTAMPTZ,
    cancel_at_period_end   BOOLEAN NOT NULL DEFAULT false,
    canceled_at            TIMESTAMPTZ,
    created_at             TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at             TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX idx_billing_subs_status ON billing_subscriptions(status);
CREATE INDEX idx_billing_subs_price ON billing_subscriptions(stripe_price_id);

CREATE TABLE billing_price_migrations (
    id                UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    old_price_id      UUID NOT NULL REFERENCES billing_prices(id),
    new_price_id      UUID NOT NULL REFERENCES billing_prices(id),
    status            TEXT NOT NULL DEFAULT 'pending',
    grace_period_days INTEGER NOT NULL DEFAULT 0,
    migrate_after     TIMESTAMPTZ NOT NULL DEFAULT now(),
    total_subs        INTEGER NOT NULL DEFAULT 0,
    migrated_subs     INTEGER NOT NULL DEFAULT 0,
    failed_subs       INTEGER NOT NULL DEFAULT 0,
    error_message     TEXT,
    started_at        TIMESTAMPTZ,
    completed_at      TIMESTAMPTZ,
    created_at        TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at        TIMESTAMPTZ NOT NULL DEFAULT now()
);

ALTER TABLE users ADD COLUMN stripe_customer_id TEXT;
CREATE UNIQUE INDEX idx_users_stripe_customer ON users(stripe_customer_id) WHERE stripe_customer_id IS NOT NULL;
