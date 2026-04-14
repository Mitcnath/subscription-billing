CREATE TYPE invoice_status AS ENUM ('paid', 'open', 'void', 'uncollectible');
CREATE TABLE invoices (
    id              BIGSERIAL           PRIMARY KEY,
    user_account_id UUID                NOT NULL,
    subscription_id BIGINT              NOT NULL,
    status          invoice_status NOT NULL,
    amount          BIGINT         NOT NULL CHECK (amount >= 0),
    currency        VARCHAR        NOT NULL,
    pdf_url         VARCHAR        NOT NULL,
    created_at      TIMESTAMPTZ         NOT NULL DEFAULT NOW(),
    updated_at      TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
    FOREIGN KEY (user_account_id) REFERENCES user_accounts(id) ON DELETE CASCADE,
    FOREIGN KEY (subscription_id) REFERENCES subscriptions(id) ON DELETE CASCADE
);