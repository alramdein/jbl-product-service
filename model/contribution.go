package model

import "time"

type Contribution struct {
	ID             string    `json:"id"`
	ReferralLinkID string    `json:"referral_link_id"`
	ContributorID  string    `json:"contributor_id"`
	AccessedAt     time.Time `json:"accessed_at"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
	DeletedAt      time.Time `json:"deleted_at"`
}
