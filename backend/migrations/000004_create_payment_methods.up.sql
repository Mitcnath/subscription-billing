CREATE TABLE payment_methods (
    id              BIGSERIAL           PRIMARY KEY,
    user_account_id UUID                NOT NULL,
    external_id     VARCHAR             NOT NULL,
    brand           VARCHAR             NOT NULL,
    last_four       VARCHAR(4)          NOT NULL,
    exp_month       SMALLINT            NOT NULL CHECK (exp_month >= 1 AND exp_month <= 12),
    -- For exp_year, we check that it's greater than or equal to the current year. This allows for future years while preventing obviously invalid years.
    -- EXTRACT(YEAR FROM CURRENT_DATE) is used to get the current year in PostgreSQL.
    exp_year        SMALLINT            NOT NULL CHECK (exp_year >= EXTRACT(YEAR FROM CURRENT_DATE)),
    is_default      BOOLEAN             NOT NULL DEFAULT FALSE,
    created_at      TIMESTAMPTZ         NOT NULL DEFAULT NOW(),
    updated_at      TIMESTAMPTZ         NOT NULL DEFAULT NOW(),
    FOREIGN KEY (user_account_id) REFERENCES user_accounts(id) ON DELETE CASCADE
);

CREATE UNIQUE INDEX one_default_per_user ON payment_methods (user_account_id) WHERE is_default = TRUE;