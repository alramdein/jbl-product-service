package model

import "time"

type ReferralLink struct {
	ID          string    `json:"id"`
	GeneratorID string    `json:"generator_id"`
	Code        string    `json:"code"`
	ExpiredAt   time.Time `json:"expired_at"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	DeletedAt   time.Time `json:"deleted_at"`
}
