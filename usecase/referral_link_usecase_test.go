package usecase

import (
	"context"
	"database/sql"
	"referral-system/mocks"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestGenerateReferralLink(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// Create mock repositories
	mockDBTransaction := mocks.NewMockIDBTransactionRepository(ctrl)
	mockReferralRepo := mocks.NewMockIReferralLinkRepository(ctrl)

	// Create a new instance of ReferralLinkUsecase
	referralLinkUsecase := NewReferralLinkUsecase(mockDBTransaction, mockReferralRepo)

	// Define the input for the test
	userID := "test-user-id"
	ctx := context.TODO()

	// Create a mock transaction
	mockTx := &sql.Tx{}

	// Setup expected calls and return values for mocks
	mockDBTransaction.EXPECT().BeginTx(ctx).Return(mockTx, nil)
	mockReferralRepo.EXPECT().DeleteReferralLinkByUserID(ctx, mockTx, userID).Return(nil)
	mockReferralRepo.EXPECT().CreateReferralLink(ctx, mockTx, gomock.Any()).Return(nil)
	mockDBTransaction.EXPECT().Commit(ctx, mockTx).Return(nil)

	// Call the method under test
	referralLink, err := referralLinkUsecase.GenerateReferralLink(ctx, userID)

	// Assertions
	assert.NoError(t, err)
	assert.NotNil(t, referralLink)
	assert.Equal(t, userID, referralLink.GeneratorID)
	assert.NotEmpty(t, referralLink.Code)
	assert.WithinDuration(t, time.Now().Add(7*24*time.Hour), referralLink.ExpiredAt, time.Second)
}
