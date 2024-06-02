package usecase

import (
	"context"
	"database/sql"
	"referral-system/mocks"
	"referral-system/model"
	"referral-system/util"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestRegisterUserGenerator(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// Create mock repositories
	mockDBTransaction := mocks.NewMockIDBTransactionRepository(ctrl)
	mockUserRepo := mocks.NewMockIUserRepository(ctrl)
	mockRoleRepo := mocks.NewMockIRoleRepository(ctrl)
	mockReferralRepo := mocks.NewMockIReferralLinkRepository(ctrl)
	mockContributionRepo := mocks.NewMockIContributionRepository(ctrl)

	// Create a new instance of UserUsecase
	jwtSecret := "test-secret"
	referralLinkExp := time.Now().Add(24 * time.Hour)
	usecase := NewUserUsecase(
		mockDBTransaction,
		mockUserRepo,
		mockRoleRepo,
		mockReferralRepo,
		mockContributionRepo,
		jwtSecret,
		referralLinkExp,
	)

	// Define the input for the test
	req := model.RegisterUserGeneratorInput{
		Email:    "test@example.com",
		Password: "password123",
	}

	ctx := context.TODO()
	tx := &sql.Tx{}

	// Setup expected calls and return values for mocks
	mockUserRepo.EXPECT().GetUserByEmailAndRole(ctx, req.Email, model.GeneratorRole).Return(nil, nil)
	mockRoleRepo.EXPECT().GetRoleByName(ctx, model.GeneratorRole).Return(&model.Role{ID: uuid.NewString(), Name: model.GeneratorRole}, nil)
	mockDBTransaction.EXPECT().BeginTx(ctx).Return(tx, nil)
	mockUserRepo.EXPECT().CreateUser(ctx, gomock.Any(), gomock.Any()).Return(nil)
	mockReferralRepo.EXPECT().CreateReferralLink(ctx, gomock.Any(), gomock.Any()).Return(nil)
	mockDBTransaction.EXPECT().Commit(ctx, tx).Return(nil)

	// Call the method under test
	resp, err := usecase.RegisterUserGenerator(ctx, req)

	// Assertions
	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, req.Email, resp.User.Email)
	assert.NotEmpty(t, resp.Token)
}

func TestRegisterUserContributor(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// Create mock repositories
	mockDBTransaction := mocks.NewMockIDBTransactionRepository(ctrl)
	mockUserRepo := mocks.NewMockIUserRepository(ctrl)
	mockRoleRepo := mocks.NewMockIRoleRepository(ctrl)
	mockReferralRepo := mocks.NewMockIReferralLinkRepository(ctrl)
	mockContributionRepo := mocks.NewMockIContributionRepository(ctrl)

	// Create a new instance of UserUsecase
	jwtSecret := "test-secret"
	referralLinkExp := time.Now().Add(24 * time.Hour)
	usecase := NewUserUsecase(
		mockDBTransaction,
		mockUserRepo,
		mockRoleRepo,
		mockReferralRepo,
		mockContributionRepo,
		jwtSecret,
		referralLinkExp,
	)

	// Define the input for the test
	req := model.RegisterUserContributorInput{
		Email:        "test@example.com",
		ReferralCode: "test-code",
	}

	ctx := context.TODO()
	tx := &sql.Tx{}

	// Setup expected calls and return values for mocks
	mockReferralRepo.EXPECT().GetReferralLinkByCode(ctx, req.ReferralCode).Return(&model.ReferralLink{ID: uuid.NewString()}, nil)
	mockReferralRepo.EXPECT().GetReferralLinkByEmail(ctx, req.Email).Return("", nil)
	mockContributionRepo.EXPECT().GetContributionByEmailAndReferralCode(ctx, req.Email, req.ReferralCode).Return(nil, nil)
	mockRoleRepo.EXPECT().GetRoleByName(ctx, model.ContributorRole).Return(&model.Role{ID: uuid.NewString(), Name: model.ContributorRole}, nil)
	mockUserRepo.EXPECT().GetUserByEmailAndRole(ctx, req.Email, model.ContributorRole).Return(nil, nil)
	mockDBTransaction.EXPECT().BeginTx(ctx).Return(tx, nil)
	mockUserRepo.EXPECT().CreateUser(ctx, gomock.Any(), gomock.Any()).Return(nil)
	mockContributionRepo.EXPECT().CreateContribution(ctx, gomock.Any(), gomock.Any()).Return(nil)
	mockDBTransaction.EXPECT().Commit(ctx, tx).Return(nil)

	// Call the method under test
	resp, err := usecase.RegisterUserContributor(ctx, req)

	// Assertions
	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, req.Email, resp.User.Email)
}

func TestLogin(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// Create mock repositories
	mockUserRepo := mocks.NewMockIUserRepository(ctrl)

	// Create a new instance of UserUsecase
	jwtSecret := "test-secret"
	referralLinkExp := time.Now().Add(24 * time.Hour)
	usecase := NewUserUsecase(
		nil,
		mockUserRepo,
		nil,
		nil,
		nil,
		jwtSecret,
		referralLinkExp,
	)

	// Define the input for the test
	email := "test@example.com"
	password := "password123"
	hashedPassword, _ := util.HashPassword(password)

	ctx := context.TODO()

	// Setup expected calls and return values for mocks
	mockUserRepo.EXPECT().GetUserByEmailAndRole(ctx, email, model.GeneratorRole).Return(&model.User{ID: uuid.NewString(), Email: email, Password: hashedPassword, RoleID: uuid.NewString()}, nil)

	// Call the method under test
	token, err := usecase.Login(ctx, email, password)

	// Assertions
	assert.NoError(t, err)
	assert.NotEmpty(t, token)
}
