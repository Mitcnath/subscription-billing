CREATE TYPE billing_interval AS ENUM (
    'daily',
    'weekly',
    'bi_weekly',
    'monthly',
    'quarterly',
    'semi_annual',
    'annual'
);
CREATE TYPE plan_status AS ENUM ('active', 'deprecated');
CREATE TABLE subscription_plans (
    id               SMALLSERIAL  		PRIMARY KEY,
    name             VARCHAR      		NOT NULL UNIQUE,
    amount           BIGINT       		NOT NULL CHECK (amount > 0),
    currency         VARCHAR      		NOT NULL,
    description      TEXT         		NOT NULL DEFAULT '',
    billing_interval billing_interval   NOT NULL,
    status           plan_status  		NOT NULL DEFAULT 'active',
    created_at       TIMESTAMPTZ 		NOT NULL DEFAULT NOW(),
    updated_at       TIMESTAMPTZ  		NOT NULL DEFAULT NOW()
);