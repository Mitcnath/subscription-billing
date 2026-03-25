CREATE TYPE subscription_status AS ENUM ('trial', 'active', 'cancelled', 'past_due');
CREATE TABLE subscriptions (
    id                                 BIGSERIAL                           PRIMARY KEY,
    user_account_id                    UUID                                NOT NULL,
    subscription_plan_id               SMALLINT                            NOT NULL,
    trial_ends_at                      TIMESTAMPTZ,
    current_period_ends_at             TIMESTAMPTZ                         NOT NULL,
    cancel_at_period_end               BOOLEAN                             NOT NULL,
    status                             subscription_status                 NOT NULL,
    created_at                         TIMESTAMPTZ                         NOT NULL DEFAULT NOW(),
    updated_at                         TIMESTAMPTZ                         NOT NULL DEFAULT NOW(),
    cancelled_at                       TIMESTAMPTZ,
    FOREIGN KEY (user_account_id)      REFERENCES user_accounts(id)        ON DELETE CASCADE,
    FOREIGN KEY (subscription_plan_id) REFERENCES subscription_plans(id)
);