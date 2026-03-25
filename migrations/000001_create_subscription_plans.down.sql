-- Migration: 000001_create_subscription_plans.down.sql
-- SAM1-15: Design subscription plan data model

DROP TABLE IF EXISTS subscription_plans;
DROP TYPE  IF EXISTS plan_status;
DROP TYPE  IF EXISTS billing_interval;
