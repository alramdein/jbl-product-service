package usecase

import (
	"context"
	"referral-system/model"
)

type IUserUsecase interface {
	RegisterUserGenerator(ctx context.Context, req model.RegisterUserGeneratorInput) (*model.RegisterUserGeneratorResponse, error)
	RegisterUserContributor(ctx context.Context, req model.RegisterUserContributorInput) (*model.RegisterUserContributorResponse, error)
	Login(ctx context.Context, email, password string) (string, error)
}
