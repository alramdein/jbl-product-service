package usecase

import "errors"

var (
	ErrEmailRequired       = errors.New("email is required")
	ErrPasswordRequired    = errors.New("password is required")
	ErrRoleRequired        = errors.New("role is required")
	ErrReferalCodeRequired = errors.New("referral code required")

	ErrEmailAlreadyExist  = errors.New("email already exist")
	ErrInvalidEmail       = errors.New("invalid email format")
	ErrInvalidRole        = errors.New("invalid role")
	ErrInvalidCredentials = errors.New("invalid credentials")

	ErrRoleNotFound               = errors.New("role not found")
	ErrReferralCodeIsNotExisit    = errors.New("referral code not found")
	ErrCantReferralToOwnCode      = errors.New("can't submit to your own referral")
	ErrCantMultipleSubmitReferral = errors.New("can't submit the same referral more than once")
)
