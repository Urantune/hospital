package models

import "time"

type CancellationPolicy struct {
	ID               int        `json:"id" db:"id"`
	PolicyName       string     `json:"policy_name" db:"policy_name"`
	HoursBefore      int        `json:"hours_before" db:"hours_before"`
	RefundPercentage float64    `json:"refund_percentage" db:"refund_percentage"`
	Description      string     `json:"description" db:"description"`
	IsActive         bool       `json:"is_active" db:"is_active"`
	Version          int        `json:"version" db:"version"`
	EffectiveFrom    time.Time  `json:"effective_from" db:"effective_from"`
	EffectiveTo      *time.Time `json:"effective_to" db:"effective_to"`
	CreatedAt        time.Time  `json:"created_at" db:"created_at"`
}
