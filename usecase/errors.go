package usecase

import "errors"

var (
	ErrEmailRequired    = errors.New("email is required")
	ErrPasswordRequired = errors.New("password is required")
	ErrRoleRequired     = errors.New("role is required")

	ErrEmailAlreadyExist   = errors.New("email already exist")
	ErrInvalidEmail        = errors.New("invalid email format")
	ErrInvalidRole         = errors.New("invalid role")
	ErrReferalCodeRequired = errors.New("referral code required")

	ErrRoleNotFound            = errors.New("role not found")
	ErrReferralCodeIsNotExisit = errors.New("referral code not found")
)
