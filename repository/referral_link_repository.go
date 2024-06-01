package repository

import (
	"context"
	"database/sql"
	"time"

	"referral-system/model"

	"github.com/sirupsen/logrus"
)

type ReferralLinkRepository struct {
	DB *sql.DB
}

// NewReferralLinkRepository creates a new instance of ReferralLinkRepository
func NewReferralLinkRepository(db *sql.DB) *ReferralLinkRepository {
	return &ReferralLinkRepository{
		DB: db,
	}
}

func (r *ReferralLinkRepository) CreateReferralLink(ctx context.Context, tx *sql.Tx, referralLink *model.ReferralLink) error {
	log := logrus.WithFields(logrus.Fields{
		"trace":        "repository.CreateReferralLink",
		"ctx":          ctx,
		"referralLink": referralLink,
	})

	query := `
		INSERT INTO referral_links (id, generator_id, code, expired_at, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6)
	`
	_, err := tx.Exec(query, referralLink.ID, referralLink.GeneratorID, referralLink.Code, referralLink.ExpiredAt, time.Now(), time.Now())
	if err != nil {
		log.Error(err)
	}
	return err
}

func (r *ReferralLinkRepository) GetReferralLinkByCode(ctx context.Context, code string) (*model.ReferralLink, error) {
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
	row := r.DB.QueryRowContext(ctx, query, code)
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
