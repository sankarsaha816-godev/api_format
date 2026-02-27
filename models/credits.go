package models

import "time"

// PlanTier defines credit limits per subscription tier
// Example IDs: individual, business, enterprise
// Credits are allocated per period.
type PlanTier struct {
	ID                 string    `gorm:"primaryKey;column:id" json:"id"`
	MonthlyCredits     int       `gorm:"column:monthly_credits" json:"monthly_credits"`
	BurstLimitPerMin   int       `gorm:"column:burst_limit_per_minute" json:"burst_limit_per_minute"`
	AllowOverage       bool      `gorm:"column:allow_overage" json:"allow_overage"`
	CreatedAt          time.Time `gorm:"column:created_at" json:"created_at"`
	UpdatedAt          time.Time `gorm:"column:updated_at" json:"updated_at"`
}

func (PlanTier) TableName() string {
	return "plan_tiers"
}

// TenantSubscription tracks credit usage for a tenant per period
// Status: active, suspended, cancelled, etc.
type TenantSubscription struct {
	ID               int64     `gorm:"primaryKey;column:id" json:"id"`
	SubscriptionID   string    `gorm:"column:subscription_id;index" json:"subscription_id"`
	TenantInfoID     int64     `gorm:"column:tenant_info_id;index" json:"tenant_info_id"`
	PlanTierID       string    `gorm:"column:plan_tier_id;index" json:"plan_tier_id"`
	Status           int    `gorm:"column:status;index" json:"status"`
	PeriodStart      time.Time `gorm:"column:period_start" json:"period_start"`
	PeriodEnd        time.Time `gorm:"column:period_end" json:"period_end"`
	CreditsTotal     int       `gorm:"column:credits_total" json:"credits_total"`
	CreditsUsed      int       `gorm:"column:credits_used" json:"credits_used"`
	CreditsRemaining int       `gorm:"column:credits_remaining" json:"credits_remaining"`
	CreatedAt        time.Time `gorm:"column:created_at" json:"created_at"`
	UpdatedAt        time.Time `gorm:"column:updated_at" json:"updated_at"`
}

func (TenantSubscription) TableName() string {
	return "tenant_subscriptions"
}

// CreditUsageLedger stores per-request credit charges and adjustments
// Status examples: charged, refunded, adjusted
// Reason is optional (admin notes, etc.).
type CreditUsageLedger struct {
	ID             int64     `gorm:"primaryKey;column:id" json:"id"`
	TenantInfoID   int64     `gorm:"column:tenant_info_id;index" json:"tenant_info_id"`
	UserID         *int64    `gorm:"column:user_id" json:"user_id"`
	Endpoint       string    `gorm:"column:endpoint" json:"endpoint"`
	CreditsCharged int       `gorm:"column:credits_charged" json:"credits_charged"`
	RequestID      string    `gorm:"column:request_id;uniqueIndex" json:"request_id"`
	Status         string    `gorm:"column:status" json:"status"`
	Reason         *string   `gorm:"column:reason" json:"reason"`
	CreatedAt      time.Time `gorm:"column:created_at" json:"created_at"`
}

func (CreditUsageLedger) TableName() string {
	return "credit_usage_ledger"
}

// APICreditLedger stores credit costs per API
type APICreditLedger struct {
	ID              int64     `gorm:"primaryKey;column:id" json:"id"`
	APIName         *string   `gorm:"column:api_name" json:"api_name"`
	APIURL          *string   `gorm:"column:api_url" json:"api_url"`
	CreditsRequired *int      `gorm:"column:credits_required" json:"credits_required"`
	Status          *int      `gorm:"column:status" json:"status"`
	CreatedAt       time.Time `gorm:"column:created_at" json:"created_at"`
	UpdatedAt       time.Time `gorm:"column:updated_at" json:"updated_at"`
}

func (APICreditLedger) TableName() string {
	return "api_credit_ledger"
}
