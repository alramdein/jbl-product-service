package repository

import (
	"context"
	"database/sql"

	"github.com/sirupsen/logrus"
)

type DBTransactionRepository struct {
	DB *sql.DB
}

// NewDBTransactionRepository creates a new instance of DBTransactionRepository
func NewDBTransactionRepository(db *sql.DB) *DBTransactionRepository {
	return &DBTransactionRepository{
		DB: db,
	}
}

func (d *DBTransactionRepository) BeginTx(ctx context.Context) (*sql.Tx, error) {
	log := logrus.WithFields(logrus.Fields{
		"trace": "repository.BeginTx",
		"ctx":   ctx,
	})

	tx, err := d.DB.BeginTx(ctx, nil)
	if err != nil {
		log.Error(err)
		return nil, err
	}
	return tx, nil
}
