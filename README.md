# Subscription Billing Service

A RESTful backend service built with Go and Gin for managing subscription-based billing. The service handles user accounts, subscription plans, payment methods, and billing histories — covering the full lifecycle from registration through recurring charges.

## Tech Stack

- **Go** — Gin (HTTP), GORM (ORM)
- **PostgreSQL** — primary database
- **golang-migrate** — schema migrations
- **Argon2id** — password hashing
- **JWT** — authentication (HS256, 1hr expiry)
- **Swagger/swaggo** — API documentation
- **Docker + Docker Compose** — containerised local environment

## Getting Started

### Local (without Docker)

```bash
# Copy and fill in environment variables
cp .env.example .env

# Run database migrations
make migrate-up

# Seed development data
make seed

# Start the server
make run
```

Swagger UI is available at `http://localhost:8080/swagger/index.html`.

### Docker

**Prerequisites:** Docker and Docker Compose installed and running.

```bash
# Copy and fill in environment variables
cp .env.example .env
```

| Command | Description |
|---|---|
| `make deploy` | Build images, run migrations, start all services |
| `make deploy-seeded` | Same as above, then seeds the database with development data |
| `make down` | Stop and remove all containers (data volume is preserved) |
| `docker compose down -v` | Stop containers and wipe the database volume |

Once running, services are available at:

| Service | URL |
|---|---|
| API | `http://localhost:8080` |
| Swagger UI | `http://localhost:8080/swagger/index.html` |
| Frontend | `http://localhost:80` |

### Docker Services

The Compose environment runs four services:

- **db** — PostgreSQL 16. Initialised from `DB_USER`, `DB_PASSWORD`, `DB_NAME` in `.env`. Data is persisted in a named Docker volume.
- **migrate** — Runs all pending migrations from `backend/migrations/` on startup, then exits.
- **backend** — The compiled Go API server.
- **frontend** — nginx serving the static HTML frontend.

## Database Schema

```mermaid
erDiagram
    subscription_plans {
        smallserial id PK
        varchar name
        bigint amount
        varchar currency
        text description
        billing_interval billing_interval
        plan_status status
        timestamptz created_at
        timestamptz updated_at
    }

    user_accounts {
        uuid id PK
        text username
        text email
        text password_hash
        timestamptz created_at
        timestamptz updated_at
        timestamptz deleted_at
    }

    subscriptions {
        bigserial id PK
        uuid user_account_id FK
        smallint subscription_plan_id FK
        subscription_status status
        timestamptz trial_ends_at
        timestamptz current_period_ends_at
        boolean cancel_at_period_end
        timestamptz cancelled_at
        timestamptz created_at
        timestamptz updated_at
    }

    payment_methods {
        bigserial id PK
        uuid user_account_id FK
        varchar external_id
        varchar brand
        varchar last_four
        smallint exp_month
        smallint exp_year
        boolean is_default
        timestamptz created_at
        timestamptz updated_at
    }

    invoices {
        bigserial id PK
        uuid user_account_id FK
        bigint subscription_id FK
        invoice_status status
        bigint amount
        varchar currency
        varchar pdf_url
        timestamptz created_at
        timestamptz updated_at
    }

    user_accounts ||--o{ subscriptions : "has"
    user_accounts ||--o{ payment_methods : "has"
    user_accounts ||--o{ invoices : "has"
    subscription_plans ||--o{ subscriptions : "subscribed to"
    subscriptions ||--o{ invoices : "generates"
```

## General Flow

```mermaid
flowchart TD
    start(["User signs up for Subscription Plan"])
    doesUserHaveAccount{"Is the user logged in?"}
    start --> doesUserHaveAccount
    doesUserHaveAccount --> |No| register["Create an Account"]
    doesUserHaveAccount --> |Yes| addPM["Add Payment Method"]
    register --> addPM
    addPM --> stripe["Tokenize payment method via Stripe"]
    stripe --> makePayment["Make Payment"]
    makePayment --> paymentIsSuccessful{"Payment Successful?"}
    paymentIsSuccessful --> |No| addPM
    paymentIsSuccessful --> |Yes| createNewSubscriptionStatus["Register for Subscription"]
    createNewSubscriptionStatus --> hasTrialPhase{"Has Trial Phase?"}
    hasTrialPhase --> |Yes| trialDuration["Free for 7 days"]
    hasTrialPhase --> |No| attemptFirstCharge["Attempt Charge"]
    trialDuration --> attemptFirstCharge
    attemptFirstCharge --> isChargeAttemptSuccessful{"Was Charge Attempt Successful?"}
    isChargeAttemptSuccessful --> |No| gracePeriod["Grace Period + Smart Retries"]
    isChargeAttemptSuccessful --> |Yes| newCycle["Record charge + update current_period_ends_at"]
    gracePeriod --> retryOutcome{"Retry Succeeds?"} --> |Yes| newCycle
    retryOutcome --> |No| markCancelled["Subscription Cancelled"]
    newCycle --> createInvoice["Generate Invoice"] --> attemptFirstCharge
```
