package repository

import (
	"context"
	"database/sql"
	"referral-system/model"

	"github.com/sirupsen/logrus"
)

type RoleRepository struct {
	DB *sql.DB
}

// NewRoleRepository creates a new instance of NewRoleRepository
func NewRoleRepository(db *sql.DB) *RoleRepository {
	return &RoleRepository{
		DB: db,
	}
}

func (r *RoleRepository) GetRoleByName(ctx context.Context, name string) (*model.Role, error) {
	log := logrus.WithFields(logrus.Fields{
		"trace": "repository.GetRoleByName",
		"ctx":   ctx,
		"name":  name,
	})

	query := "SELECT id, name, created_at, updated_at FROM roles WHERE name = $1 AND deleted_at IS NULL"
	var role model.Role
	err := r.DB.QueryRowContext(ctx, query, name).Scan(&role.ID, &role.Name, &role.CreatedAt, &role.UpdatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil // No user found
		}
		log.Error(err)
		return nil, err
	}
	return &role, nil
}
