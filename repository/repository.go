package repository

import (
	"context"
	"database/sql"
	"referral-system/model"
)

type IUserRepository interface {
	CreateUser(ctx context.Context, tx *sql.Tx, user *model.User) error
	GetUserByEmailAndRole(ctx context.Context, email string, roleName string) (*model.User, error)
}

type IReferralLinkRepository interface {
	CreateReferralLink(ctx context.Context, tx *sql.Tx, user *model.ReferralLink) error
	GetReferralLinkByCode(ctx context.Context, code string) (*model.ReferralLink, error)
	GetReferralLinkByEmail(ctx context.Context, email string) (string, error)
	DeleteReferralLinkByUserID(ctx context.Context, tx *sql.Tx, userID string) error
}

type IRoleRepository interface {
	GetRoleByName(ctx context.Context, name string) (*model.Role, error)
}

type IContributionRepository interface {
	CreateContribution(ctx context.Context, tx *sql.Tx, contribution *model.Contribution) error
	GetContributionByEmailAndReferralCode(ctx context.Context, email string, referralCode string) (*model.Contribution, error)
}

type IDBTransactionRepository interface {
	BeginTx(ctx context.Context) (*sql.Tx, error)
}
