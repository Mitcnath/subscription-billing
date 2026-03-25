# Subscription Billing Service

A RESTful backend service built with Go and Gin for managing subscription-based billing. The service handles user accounts, subscription plans, payment methods, and billing histories — covering the full lifecycle from registration through recurring charges.

## Tech Stack

- **Go** — Gin (HTTP), GORM (ORM)
- **PostgreSQL** — primary database
- **golang-migrate** — schema migrations
- **Argon2id** — password hashing
- **JWT** — authentication (HS256, 1hr expiry)
- **Swagger/swaggo** — API documentation

## Getting Started

```bash
# Copy environment variables
cp .env.example .env

# Run database migrations
make migrate-up

# Seed development data
make seed

# Start the server
make run
```

Swagger UI is available at `http://localhost:8080/swagger/index.html`.

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
