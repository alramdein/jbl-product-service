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

type UserUsecase struct {
	DbTransaction    repository.IDBTransactionRepository
	UserRepo         repository.IUserRepository
	RoleRepo         repository.IRoleRepository
	ReferralRepo     repository.IReferralLinkRepository
	ContributionRepo repository.IContributionRepository
	JWTSecret        string
	ReferralLinkExp  time.Time
}

// NewUserRepository creates a new instance of userRepository
func NewUserUsecase(
	DbTransaction repository.IDBTransactionRepository,
	UserRepo repository.IUserRepository,
	RoleRepo repository.IRoleRepository,
	ReferralRepo repository.IReferralLinkRepository,
	ContributionRepo repository.IContributionRepository,
	JWTSecret string,
	ReferralLinkExp time.Time,
) *UserUsecase {
	return &UserUsecase{
		DbTransaction:    DbTransaction,
		UserRepo:         UserRepo,
		RoleRepo:         RoleRepo,
		ReferralRepo:     ReferralRepo,
		ContributionRepo: ContributionRepo,
		JWTSecret:        JWTSecret,
		ReferralLinkExp:  ReferralLinkExp,
	}
}

func (u *UserUsecase) RegisterUserGenerator(ctx context.Context, req model.RegisterUserGeneratorInput) (*model.RegisterUserGeneratorResponse, error) {
	log := logrus.WithFields(logrus.Fields{
		"trace":   "usecase.RegisterUserGenerator",
		"ctx":     ctx,
		"payload": req,
	})

	err := u.validateRegisterUserGeneratorInput(&req)
	if err != nil {
		return nil, err
	}

	user, err := u.UserRepo.GetUserByEmailAndRole(ctx, req.Email, model.GeneratorRole)
	if err != nil {
		log.Error(err)
		return nil, err
	}
	if user != nil {
		return nil, ErrEmailAlreadyExist
	}

	hashedPassword, err := util.HashPassword(req.Password)
	if err != nil {
		log.Error(err)
		return nil, err
	}

	role, err := u.RoleRepo.GetRoleByName(ctx, model.GeneratorRole)
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

	referralLink := &model.ReferralLink{
		ID:          uuid.NewString(),
		GeneratorID: newUser.ID,
		Code:        util.GenerateUniqueCode(),
		ExpiredAt:   u.ReferralLinkExp,
		CreatedAt:   now,
		UpdatedAt:   now,
	}
	err = u.ReferralRepo.CreateReferralLink(ctx, tx, referralLink)
	if err != nil {
		log.Error(err)
		tx.Rollback()
		return nil, err
	}

	err = tx.Commit()
	if err != nil {
		log.Error(err)
		return nil, err
	}

	token, err := util.GenerateJWT(newUser.ID, newUser.RoleID, u.JWTSecret)
	if err != nil {
		log.Error(err)
		return nil, err
	}

	// only for information
	newUser.Role = role

	return &model.RegisterUserGeneratorResponse{
		User:         *newUser,
		ReferralLink: referralLink,
		Token:        token,
	}, nil
}

func (u *UserUsecase) RegisterUserContributor(ctx context.Context, req model.RegisterUserContributorInput) (*model.RegisterUserContributorResponse, error) {
	log := logrus.WithFields(logrus.Fields{
		"trace":   "usecase.RegisterUserContributor",
		"ctx":     ctx,
		"payload": req,
	})

	err := u.validateRegisterUserContributorInput(&req)
	if err != nil {
		return nil, err
	}

	// validate the existing of referral code
	ref, err := u.ReferralRepo.GetReferralLinkByCode(ctx, req.ReferralCode)
	if err != nil {
		return nil, err
	}
	if ref == nil {
		return nil, ErrReferralCodeIsNotExisit
	}

	// check if it join it owns referral code
	code, err := u.ReferralRepo.GetReferralLinkByEmail(ctx, req.Email)
	if err != nil {
		return nil, err
	}
	if code == req.ReferralCode {
		return nil, ErrCantReferralToOwnCode
	}

	// check if contributor already join the referral code
	contribution, err := u.ContributionRepo.GetContributionByEmailAndReferralCode(ctx, req.Email, req.ReferralCode)
	if err != nil {
		return nil, err
	}
	if contribution != nil {
		return nil, ErrCantMultipleSubmitReferral
	}

	// get role metadata for insertion
	role, err := u.RoleRepo.GetRoleByName(ctx, model.ContributorRole)
	if err != nil {
		log.Error(err)
		return nil, err
	}
	if role == nil {
		return nil, ErrInvalidRole
	}

	user, err := u.UserRepo.GetUserByEmailAndRole(ctx, req.Email, model.ContributorRole)
	if err != nil {
		log.Error(err)
		return nil, err
	}

	tx, err := u.DbTransaction.BeginTx(ctx)
	if err != nil {
		log.Error(err)
		return nil, err
	}

	if user == nil {
		now := time.Now().UTC()
		newUser := &model.User{
			ID:        uuid.NewString(),
			Email:     req.Email,
			RoleID:    role.ID,
			CreatedAt: now,
			UpdatedAt: now,
		}

		err = u.UserRepo.CreateUser(ctx, tx, newUser)
		if err != nil {
			log.Error(err)
			tx.Rollback()
			return nil, err
		}

		user = newUser
	}

	now := time.Now().UTC()
	referralLink := &model.Contribution{
		ID:             uuid.NewString(),
		ReferralLinkID: ref.ID,
		ContributorID:  user.ID,
		AccessedAt:     now,
		CreatedAt:      now,
		UpdatedAt:      now,
	}
	err = u.ContributionRepo.CreateContribution(ctx, tx, referralLink)
	if err != nil {
		return nil, err
	}

	err = tx.Commit()
	if err != nil {
		log.Error(err)
		return nil, err
	}

	return &model.RegisterUserContributorResponse{
		User: *user,
	}, nil
}

func (u *UserUsecase) Login(ctx context.Context, email, password string) (string, error) {
	log := logrus.WithFields(logrus.Fields{
		"trace": "usecase.Login",
		"ctx":   ctx,
		"email": email,
	})

	user, err := u.UserRepo.GetUserByEmailAndRole(ctx, email, model.GeneratorRole)
	if err != nil {
		log.Error(err)
		return "", err
	}
	if user == nil {
		return "", ErrInvalidCredentials
	}

	if !util.CheckPasswordHash(password, user.Password) {
		return "", ErrInvalidCredentials
	}

	token, err := util.GenerateJWT(user.ID, user.RoleID, u.JWTSecret)
	if err != nil {
		log.Error(err)
		return "", err
	}

	return token, nil
}

// validateRegisterUserGeneratorInput validates the validateRegisterUserGeneratorInput struct
func (u *UserUsecase) validateRegisterUserGeneratorInput(req *model.RegisterUserGeneratorInput) error {
	if req.Email == "" {
		return ErrEmailRequired
	}
	if req.Password == "" {
		return ErrPasswordRequired
	}
	if !model.EmailRegex.MatchString(req.Email) {
		return ErrInvalidEmail
	}
	return nil
}

// RegisterUserContributorInput validates the RegisterUserContributorInput struct
func (u *UserUsecase) validateRegisterUserContributorInput(req *model.RegisterUserContributorInput) error {
	if req.Email == "" {
		return ErrEmailRequired
	}
	if req.ReferralCode == "" {
		return ErrReferalCodeRequired
	}
	if !model.EmailRegex.MatchString(req.Email) {
		return ErrInvalidEmail
	}
	return nil
}
