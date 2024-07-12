package usecase

import (
	"context"
	"referral-system/model"
	"referral-system/repository"
	"referral-system/util"
	"time"

	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

type ReferralLinkUsecase interface {
	GenerateReferralLink(ctx context.Context, userID string) (*model.ReferralLink, error)
}

type referralLinkUsecase struct {
	DbTransaction repository.IDBTransactionRepository
	ReferralRepo  repository.IReferralLinkRepository
}

func NewReferralLinkUsecase(
	DbTransaction repository.IDBTransactionRepository,
	referralRepo repository.IReferralLinkRepository,
) ReferralLinkUsecase {
	return &referralLinkUsecase{
		DbTransaction: DbTransaction,
		ReferralRepo:  referralRepo,
	}
}

func (u *referralLinkUsecase) GenerateReferralLink(ctx context.Context, userID string) (*model.ReferralLink, error) {
	log := logrus.WithFields(logrus.Fields{
		"trace":  "usecase.GenerateReferralLink",
		"ctx":    ctx,
		"userID": userID,
	})

	tx, err := u.DbTransaction.BeginTx(ctx)
	if err != nil {
		log.Error(err)
		return nil, err
	}

	err = u.ReferralRepo.DeleteReferralLinkByUserID(ctx, tx, userID)
	if err != nil {
		log.Error(err)
		tx.Rollback()
		return nil, err
	}

	now := time.Now().UTC()
	newReferralLink := &model.ReferralLink{
		ID:          uuid.NewString(),
		GeneratorID: userID,
		Code:        util.GenerateUniqueCode(),
		ExpiredAt:   now.Add(7 * 24 * time.Hour),
		CreatedAt:   now,
		UpdatedAt:   now,
	}
	err = u.ReferralRepo.CreateReferralLink(ctx, tx, newReferralLink)
	if err != nil {
		log.Error(err)
		tx.Rollback()
		return nil, err
	}

	err = u.DbTransaction.Commit(ctx, tx)
	if err != nil {
		log.Error(err)
		return nil, err
	}

	return newReferralLink, nil
}
