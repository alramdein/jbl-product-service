package usecase

import (
	"context"
	"database/sql"
	"referral-system/model"
	"referral-system/repository"
	"referral-system/util"
	"regexp"
	"time"

	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

type UserUsecase struct {
	DbTransaction    repository.IDBTransactionRepository
	UserRepo         repository.IUserRepository
	RoleRepo         repository.IRoleRepository
	ReferralRepo     repository.IReferralLinkRepository
	ContributionRepo repository.IContributionRepository
}

// NewUserRepository creates a new instance of userRepository
func NewUserUsecase(
	DbTransaction repository.IDBTransactionRepository,
	UserRepo repository.IUserRepository,
	RoleRepo repository.IRoleRepository,
	ReferralRepo repository.IReferralLinkRepository,
	ContributionRepo repository.IContributionRepository,
) *UserUsecase {
	return &UserUsecase{
		DbTransaction:    DbTransaction,
		UserRepo:         UserRepo,
		RoleRepo:         RoleRepo,
		ReferralRepo:     ReferralRepo,
		ContributionRepo: ContributionRepo,
	}
}

func (u *UserUsecase) RegisterUser(ctx context.Context, req model.RegisterUserInput) (*model.RegisterUserResponse, error) {
	log := logrus.WithFields(logrus.Fields{
		"trace":   "usecase.RegisterUser",
		"ctx":     ctx,
		"payload": req,
	})

	err := u.validateRegisterUserInput(&req)
	if err != nil {
		return nil, err
	}

	hashedPassword, err := util.HashPassword(req.Password)
	if err != nil {
		log.Error(err)
		return nil, err
	}

	user, err := u.UserRepo.GetUserByEmailAndRole(ctx, req.Email, req.Role)
	if err != nil {
		log.Error(err)
		return nil, err
	}
	if user != nil {
		return nil, ErrEmailAlreadyExist
	}

	role, err := u.RoleRepo.GetRoleByName(ctx, req.Role)
	if err != nil {
		log.Error(err)
		return nil, err
	}
	if role == nil {
		return nil, ErrInvalidRole
	}

	now := time.Now().UTC()
	newUser := &model.User{
		ID:        uuid.NewString(),
		Email:     req.Email,
		Password:  hashedPassword,
		RoleID:    role.ID,
		CreatedAt: now,
		UpdatedAt: now,
	}

	tx, err := u.DbTransaction.BeginTx(ctx)
	if err != nil {
		log.Error(err)
		return nil, err
	}

	err = u.UserRepo.CreateUser(ctx, tx, newUser)
	if err != nil {
		log.Error(err)
		tx.Rollback()
		return nil, err
	}

	var referralLink *model.ReferralLink
	switch req.Role {
	case model.GeneratorRole:
		referralLink, err = u.handleCreateReferralLink(ctx, tx, newUser.ID)
		if err != nil {
			log.Error(err)
			tx.Rollback()
			return nil, err
		}
	case model.ContributorRole:
		_, err = u.handleAddContribution(ctx, tx, newUser.ID, req.ReferralCode)
		if err != nil {
			log.Error(err)
			tx.Rollback()
			return nil, err
		}
	default:
		log.Error(err)
		tx.Rollback()
		return nil, ErrInvalidRole
	}

	err = tx.Commit()
	if err != nil {
		log.Error(err)
		return nil, err
	}

	// only for information
	newUser.Role = role

	return &model.RegisterUserResponse{
		User:         *newUser,
		ReferralLink: referralLink,
	}, nil
}

func (u *UserUsecase) handleCreateReferralLink(ctx context.Context, tx *sql.Tx, userID string) (*model.ReferralLink, error) {
	now := time.Now().UTC()
	referralLink := &model.ReferralLink{
		ID:          uuid.NewString(),
		GeneratorID: userID,
		Code:        util.GenerateUniqueCode(),
		ExpiredAt:   time.Now().Add(7 * 24 * time.Hour),
		CreatedAt:   now,
		UpdatedAt:   now,
	}
	err := u.ReferralRepo.CreateReferralLink(ctx, tx, referralLink)
	if err != nil {
		return nil, err
	}

	return referralLink, nil
}

func (u *UserUsecase) handleAddContribution(ctx context.Context, tx *sql.Tx, userID, referralCode string) (*model.Contribution, error) {
	ref, err := u.ReferralRepo.GetReferralLinkByCode(ctx, referralCode)
	if err != nil {
		return nil, err
	}
	if ref == nil {
		return nil, ErrReferralCodeIsNotExisit
	}

	now := time.Now().UTC()
	referralLink := &model.Contribution{
		ID:             uuid.NewString(),
		ReferralLinkID: ref.ID,
		ContributorID:  userID,
		AccessedAt:     now,
		CreatedAt:      now,
		UpdatedAt:      now,
	}
	err = u.ContributionRepo.CreateContribution(ctx, tx, referralLink)
	if err != nil {
		return nil, err
	}

	return referralLink, nil
}

// validateRegisterUserInput validates the RegisterUserInput struct
func (u *UserUsecase) validateRegisterUserInput(req *model.RegisterUserInput) error {
	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	if req.Email == "" {
		return ErrEmailRequired
	}
	if req.Password == "" {
		return ErrPasswordRequired
	}
	if req.Role == "" {
		return ErrRoleRequired
	}
	if !emailRegex.MatchString(req.Email) {
		return ErrInvalidEmail
	}
	if req.Role != model.GeneratorRole && req.Role != model.ContributorRole {
		return ErrInvalidRole
	}
	if req.Role == model.ContributorRole && req.ReferralCode == "" {
		return ErrReferalCodeRequired
	}
	return nil
}
