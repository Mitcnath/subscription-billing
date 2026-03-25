CREATE TYPE billing_status AS ENUM ('paid', 'open', 'void', 'uncollectible');
CREATE TABLE billing_histories (
    id              BIGSERIAL           PRIMARY KEY,
    user_account_id UUID                NOT NULL,
    subscription_id BIGINT              NOT NULL,
    status          billing_status     NOT NULL,
    amount_paid      BIGINT              NOT NULL CHECK (amount_paid >= 0),
    pdf_url         VARCHAR             NOT NULL,
    created_at      TIMESTAMPTZ         NOT NULL DEFAULT NOW(),
    updated_at      TIMESTAMPTZ         NOT NULL DEFAULT NOW(),
    FOREIGN KEY (user_account_id) REFERENCES user_accounts(id) ON DELETE CASCADE,
    FOREIGN KEY (subscription_id) REFERENCES subscriptions(id) ON DELETE CASCADE
);