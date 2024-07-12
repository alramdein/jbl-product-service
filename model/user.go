package model

import (
	"regexp"
	"time"
)

var (
	EmailRegex = regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
)

// User :nodoc:
type User struct {
	ID        string    `json:"id"`
	Email     string    `json:"email"`
	Password  string    `json:"password"`
	RoleID    string    `json:"role_id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	DeletedAt time.Time `json:"deleted_at"`

	Role *Role `json:"role,omitempty" pq:"-"`
}

// RegisterUserRequest payload for handler layer
type RegisterUserRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// RegisterUserInput payload for usecase layer
type RegisterUserGeneratorInput struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type RegisterUserContributorInput struct {
	Email        string `json:"email"`
	ReferralCode string `json:"referral_code"`
}

// RegisterUserResponse :nodoc:
type RegisterUserGeneratorResponse struct {
	User         User          `json:"user"`
	ReferralLink *ReferralLink `json:"referral_link,omitempty"`
	Token        string        `json:"token"`
}

type RegisterUserContributorResponse struct {
	User         User          `json:"user"`
	ReferralLink *ReferralLink `json:"referral_link,omitempty"`
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}
