package model

import "time"

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
type RegisterUserInput struct {
	Email        string `json:"email"`
	Password     string `json:"password"`
	Role         string `json:"role"`
	ReferralCode string `json:"referral_code"`
}

// RegisterUserResponse :nodoc:
type RegisterUserResponse struct {
	User         User          `json:"user"`
	ReferralLink *ReferralLink `json:"referral_link,omitempty"`
}
