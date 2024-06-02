package repository

import (
	"context"
	"database/sql"

	"referral-system/model"

	"github.com/sirupsen/logrus"
)

type ContributionRepository struct {
	DB *sql.DB
}

// NewContributionRepository creates a new instance of ContributionRepository
func NewContributionRepository(db *sql.DB) *ContributionRepository {
	return &ContributionRepository{
		DB: db,
	}
}

func (c *ContributionRepository) CreateContribution(ctx context.Context, tx *sql.Tx, contribution *model.Contribution) error {
	log := logrus.WithFields(logrus.Fields{
		"trace":        "repository.CreateContribution",
		"ctx":          ctx,
		"contribution": contribution,
	})

	query := `
	INSERT INTO contributions (id, referral_link_id, contributor_id, accessed_at, created_at, updated_at)
	VALUES ($1, $2, $3, $4, $5, $6)
`
	_, err := tx.Exec(query, contribution.ID, contribution.ReferralLinkID, contribution.ContributorID, contribution.AccessedAt, contribution.CreatedAt, contribution.UpdatedAt)
	if err != nil {
		log.Error(err)
	}
	return err
}

func (c *ContributionRepository) GetReferralLinkByCode(ctx context.Context, code string) (*model.ReferralLink, error) {
	log := logrus.WithFields(logrus.Fields{
		"trace": "repository.GetReferralLinkByCode",
		"ctx":   ctx,
		"code":  code,
	})

	query := `
        SELECT id, code, created_at, updated_at
        FROM referral_links
        WHERE code = $1 AND deleted_at IS NULL
    `
	row := c.DB.QueryRowContext(ctx, query, code)
	var referralLink model.ReferralLink
	err := row.Scan(
		&referralLink.ID,
		&referralLink.Code,
		&referralLink.CreatedAt,
		&referralLink.UpdatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil // Referral link not found
		}
		log.Error(err)
		return nil, err
	}
	return &referralLink, nil
}

func (c *ContributionRepository) GetContributionByEmailAndReferralCode(ctx context.Context, email string, referralCode string) (*model.Contribution, error) {
	log := logrus.WithFields(logrus.Fields{
		"trace":        "repository.GetContributionByEmailAndReferralCode",
		"ctx":          ctx,
		"email":        email,
		"referralCode": referralCode,
	})

	query := `
        SELECT c.id, c.referral_link_id, c.contributor_id, c.accessed_at, c.created_at, c.updated_at
        FROM contributions c
        JOIN referral_links rl ON c.referral_link_id = rl.id AND rl.code = $2
        JOIN users u ON c.contributor_id = u.id AND u.email = $1 AND c.deleted_at IS NULL
    `
	var contribution model.Contribution
	err := c.DB.QueryRowContext(ctx, query, email, referralCode).
		Scan(&contribution.ID, &contribution.ReferralLinkID, &contribution.ContributorID, &contribution.AccessedAt, &contribution.CreatedAt, &contribution.UpdatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil // No contribution found
		}
		log.Error(err)
		return nil, err
	}
	return &contribution, nil
}
