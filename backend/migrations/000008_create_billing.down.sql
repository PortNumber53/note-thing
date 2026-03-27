DROP TABLE IF EXISTS billing_price_migrations;
DROP TABLE IF EXISTS billing_subscriptions;
DROP TABLE IF EXISTS billing_prices;
DROP TABLE IF EXISTS billing_products;
DROP INDEX IF EXISTS idx_users_stripe_customer;
ALTER TABLE users DROP COLUMN IF EXISTS stripe_customer_id;
