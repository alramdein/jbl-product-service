package usecase

import (
	"context"
	"referral-system/model"
)

type IUserUsecase interface {
	RegisterUser(ctx context.Context, req model.RegisterUserInput) (*model.RegisterUserResponse, error)
}
