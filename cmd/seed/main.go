package main

import (
	"log"
	"log/slog"
	"os"
	"time"

	"github.com/joho/godotenv"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"billingService/backend/internal/accounts"
	"billingService/backend/internal/payment"
	"billingService/backend/internal/plans"
	"billingService/backend/internal/statuses"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading environment variables.")
	}

	db, err := gorm.Open(postgres.Open(os.Getenv("DSN")), &gorm.Config{TranslateError: true})
	if err != nil {
		log.Fatal(err)
	}

	// Truncate all tables in reverse FK dependency order before seeding
	if err := db.Exec("TRUNCATE billing_histories, subscriptions, payment_methods, user_accounts, subscription_plans RESTART IDENTITY CASCADE").Error; err != nil {
		log.Fatal("failed to truncate tables: ", err)
	}
	slog.Info("Truncated all tables")

	hash, err := bcrypt.GenerateFromPassword([]byte("seed_password_123"), bcrypt.DefaultCost)
	if err != nil {
		log.Fatal("failed to hash seed password: ", err)
	}
	passwordHash := string(hash)

	now := time.Now()
	trialEnd := now.AddDate(0, 0, 14)
	monthEnd := now.AddDate(0, 1, 0)
	quarterEnd := now.AddDate(0, 3, 0)
	yearEnd := now.AddDate(1, 0, 0)
	pastPeriodEnd := now.AddDate(0, -1, 0)
	cancelledAt := now.AddDate(0, 0, -7)

	// --- Seed Plans ---
	slog.Info("Seeding subscription plans...")

	seedPlans := []plans.SubscriptionPlans{
		{
			Name:            "Starter",
			Description:     "Basic plan for individuals and hobbyists",
			Amount:          999,
			Currency:        "CAD",
			BillingInterval: plans.BillingIntervalMonthly,
		},
		{
			Name:            "Pro",
			Description:     "Advanced features for professionals",
			Amount:          2999,
			Currency:        "CAD",
			BillingInterval: plans.BillingIntervalMonthly,
		},
		{
			Name:            "Business",
			Description:     "Team-oriented plan with priority support",
			Amount:          9999,
			Currency:        "CAD",
			BillingInterval: plans.BillingIntervalQuarterly,
		},
		{
			Name:            "Enterprise",
			Description:     "Full-featured plan for large organizations",
			Amount:          29999,
			Currency:        "CAD",
			BillingInterval: plans.BillingIntervalAnnual,
		},
	}

	for i := range seedPlans {
		if err := db.Create(&seedPlans[i]).Error; err != nil {
			log.Fatalf("failed to seed plan %s: %v", seedPlans[i].Name, err)
		}
		slog.Info("Seeded plan", "name", seedPlans[i].Name, "id", seedPlans[i].ID)
	}

	// --- Seed Users ---
	slog.Info("Seeding user accounts...")

	seedUsers := []accounts.UserAccounts{
		{Email: "alice.martin@example.com", Username: "alice_martin", PasswordHash: passwordHash},
		{Email: "bob.chen@example.com", Username: "bob_chen", PasswordHash: passwordHash},
		{Email: "carlos.reyes@example.com", Username: "carlos_reyes", PasswordHash: passwordHash},
		{Email: "diana.patel@example.com", Username: "diana_patel", PasswordHash: passwordHash},
		{Email: "ethan.kowalski@example.com", Username: "ethan_kowalski", PasswordHash: passwordHash},
		{Email: "fatima.hassan@example.com", Username: "fatima_hassan", PasswordHash: passwordHash},
		{Email: "grace.liu@example.com", Username: "grace_liu", PasswordHash: passwordHash},
		{Email: "henry.okafor@example.com", Username: "henry_okafor", PasswordHash: passwordHash},
		{Email: "isabelle.tremblay@example.com", Username: "isabelle_tremblay", PasswordHash: passwordHash},
		{Email: "james.fitzgerald@example.com", Username: "james_fitzgerald", PasswordHash: passwordHash},
		{Email: "kira.yamamoto@example.com", Username: "kira_yamamoto", PasswordHash: passwordHash},
		{Email: "liam.oconnor@example.com", Username: "liam_oconnor", PasswordHash: passwordHash},
		{Email: "maya.singh@example.com", Username: "maya_singh", PasswordHash: passwordHash},
		{Email: "noah.petrov@example.com", Username: "noah_petrov", PasswordHash: passwordHash},
		{Email: "olivia.santos@example.com", Username: "olivia_santos", PasswordHash: passwordHash},
		{Email: "pascal.dubois@example.com", Username: "pascal_dubois", PasswordHash: passwordHash},
		{Email: "quinn.nakamura@example.com", Username: "quinn_nakamura", PasswordHash: passwordHash},
		{Email: "rachel.ibrahim@example.com", Username: "rachel_ibrahim", PasswordHash: passwordHash},
		{Email: "samuel.kim@example.com", Username: "samuel_kim", PasswordHash: passwordHash},
		{Email: "tara.morrison@example.com", Username: "tara_morrison", PasswordHash: passwordHash},
	}

	for i := range seedUsers {
		if err := db.Create(&seedUsers[i]).Error; err != nil {
			log.Fatalf("failed to seed user %s: %v", seedUsers[i].Username, err)
		}
		slog.Info("Seeded user", "username", seedUsers[i].Username, "id", seedUsers[i].ID)
	}

	// --- Seed Subscriptions ---
	slog.Info("Seeding subscriptions...")

	// Plan indices: 0=Starter, 1=Pro, 2=Business, 3=Enterprise
	// Status distribution: 4 trial, 10 active, 3 past_due, 3 cancelled
	type subConfig struct {
		userIdx           int
		planIdx           int
		status            statuses.SubscriptionStatus
		trialEndsAt       *time.Time
		currentPeriodEnds time.Time
		cancelAtPeriodEnd bool
		cancelledAt       *time.Time
	}

	configs := []subConfig{
		// Trial (4) — trialEndsAt is set to a future date
		{0, 0, statuses.SubscriptionStatusTrial, &trialEnd, monthEnd, false, nil},
		{1, 1, statuses.SubscriptionStatusTrial, &trialEnd, monthEnd, false, nil},
		{2, 0, statuses.SubscriptionStatusTrial, &trialEnd, monthEnd, false, nil},
		{3, 2, statuses.SubscriptionStatusTrial, &trialEnd, quarterEnd, false, nil},
		// Active (10) — trial has already ended, trialEndsAt is nil
		{4, 0, statuses.SubscriptionStatusActive, nil, monthEnd, false, nil},
		{5, 1, statuses.SubscriptionStatusActive, nil, monthEnd, false, nil},
		{6, 2, statuses.SubscriptionStatusActive, nil, quarterEnd, false, nil},
		{7, 3, statuses.SubscriptionStatusActive, nil, yearEnd, false, nil},
		{8, 1, statuses.SubscriptionStatusActive, nil, monthEnd, false, nil},
		{9, 0, statuses.SubscriptionStatusActive, nil, monthEnd, true, nil},
		{10, 2, statuses.SubscriptionStatusActive, nil, quarterEnd, false, nil},
		{11, 3, statuses.SubscriptionStatusActive, nil, yearEnd, false, nil},
		{12, 1, statuses.SubscriptionStatusActive, nil, monthEnd, false, nil},
		{13, 0, statuses.SubscriptionStatusActive, nil, monthEnd, true, nil},
		// Past Due (3)
		{14, 0, statuses.SubscriptionStatusPastDue, nil, pastPeriodEnd, false, nil},
		{15, 1, statuses.SubscriptionStatusPastDue, nil, pastPeriodEnd, false, nil},
		{16, 2, statuses.SubscriptionStatusPastDue, nil, pastPeriodEnd, false, nil},
		// Cancelled (3)
		{17, 0, statuses.SubscriptionStatusCanceled, nil, pastPeriodEnd, true, &cancelledAt},
		{18, 1, statuses.SubscriptionStatusCanceled, nil, pastPeriodEnd, true, &cancelledAt},
		{19, 3, statuses.SubscriptionStatusCanceled, nil, yearEnd, true, &cancelledAt},
	}

	for _, cfg := range configs {
		sub := statuses.Subscriptions{
			UserAccountID:       seedUsers[cfg.userIdx].ID,
			SubscriptionPlanID:  seedPlans[cfg.planIdx].ID,
			Status:              cfg.status,
			TrialEndsAt:         cfg.trialEndsAt,
			CurrentPeriodEndsAt: cfg.currentPeriodEnds,
			CancelAtPeriodEnd:   cfg.cancelAtPeriodEnd,
			CancelledAt:         cfg.cancelledAt,
		}

		if err := db.Create(&sub).Error; err != nil {
			slog.Error("Failed to seed subscription", "user", seedUsers[cfg.userIdx].Username, "error", err)
			continue
		}
		slog.Info("Seeded subscription",
			"id", sub.ID,
			"user", seedUsers[cfg.userIdx].Username,
			"plan", seedPlans[cfg.planIdx].Name,
			"status", sub.Status,
		)
	}

	// --- Seed Payment Methods ---
	slog.Info("Seeding payment methods...")

	type pmConfig struct {
		userIdx    int
		externalID string
		brand      string
		lastFour   string
		expMonth   int16
		expYear    int16
	}

	pmConfigs := []pmConfig{
		{0, "pm_1A2B3C4D5E6F7G8H", "Visa", "4242", 12, 2027},
		{1, "pm_2B3C4D5E6F7G8H9I", "Mastercard", "5555", 8, 2028},
		{2, "pm_3C4D5E6F7G8H9I0J", "Visa", "1234", 3, 2026},
		{3, "pm_4D5E6F7G8H9I0J1K", "Amex", "0005", 11, 2029},
		{4, "pm_5E6F7G8H9I0J1K2L", "Mastercard", "9999", 6, 2027},
		{5, "pm_6F7G8H9I0J1K2L3M", "Visa", "3782", 1, 2028},
		{6, "pm_7G8H9I0J1K2L3M4N", "Discover", "6011", 9, 2026},
		{7, "pm_8H9I0J1K2L3M4N5O", "Visa", "4111", 4, 2030},
		{8, "pm_9I0J1K2L3M4N5O6P", "Mastercard", "2223", 7, 2027},
		{9, "pm_0J1K2L3M4N5O6P7Q", "Amex", "3714", 2, 2028},
		{10, "pm_1K2L3M4N5O6P7Q8R", "Visa", "4000", 10, 2026},
		{11, "pm_2L3M4N5O6P7Q8R9S", "Mastercard", "5105", 5, 2029},
		{12, "pm_3M4N5O6P7Q8R9S0T", "Visa", "4012", 12, 2027},
		{13, "pm_4N5O6P7Q8R9S0T1U", "Discover", "6512", 3, 2028},
		{14, "pm_5O6P7Q8R9S0T1U2V", "Visa", "4222", 8, 2026},
		{15, "pm_6P7Q8R9S0T1U2V3W", "Mastercard", "5200", 11, 2030},
		{16, "pm_7Q8R9S0T1U2V3W4X", "Amex", "3787", 6, 2027},
		{17, "pm_8R9S0T1U2V3W4X5Y", "Visa", "4916", 1, 2028},
		{18, "pm_9S0T1U2V3W4X5Y6Z", "Mastercard", "5425", 9, 2029},
		{19, "pm_0T1U2V3W4X5Y6Z7A", "Visa", "4539", 4, 2027},
	}

	for _, cfg := range pmConfigs {
		pm := payment.PaymentMethod{
			UserAccountID: seedUsers[cfg.userIdx].ID,
			ExternalID:    cfg.externalID,
			Brand:         cfg.brand,
			LastFour:      cfg.lastFour,
			ExpMonth:      cfg.expMonth,
			ExpYear:       cfg.expYear,
			IsDefault:     true,
		}

		if err := db.Create(&pm).Error; err != nil {
			slog.Error("Failed to seed payment method", "user", seedUsers[cfg.userIdx].Username, "error", err)
			continue
		}
		slog.Info("Seeded payment method",
			"id", pm.ID,
			"user", seedUsers[cfg.userIdx].Username,
			"brand", pm.Brand,
			"last_four", pm.LastFour,
		)
	}

	slog.Info("Seeding complete.")
}
