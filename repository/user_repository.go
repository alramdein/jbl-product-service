package repository

import (
	"context"
	"database/sql"

	"referral-system/model"

	"github.com/sirupsen/logrus"
)

type UserRepository struct {
	DB *sql.DB
}

// NewUserRepository creates a new instance of userRepository
func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{
		DB: db,
	}
}

func (u *UserRepository) CreateUser(ctx context.Context, tx *sql.Tx, user *model.User) error {
	log := logrus.WithFields(logrus.Fields{
		"trace": "repository.CreateUser",
		"ctx":   ctx,
		"user":  user,
	})

	query := `
		INSERT INTO users (id, email, password, role_id, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6)
	`
	_, err := tx.Exec(query, user.ID, user.Email, user.Password, user.RoleID, user.CreatedAt, user.UpdatedAt)
	if err != nil {
		log.Error(err)
	}
	return err
}

func (u *UserRepository) GetUserByEmailAndRole(ctx context.Context, email string, roleName string) (*model.User, error) {
	log := logrus.WithFields(logrus.Fields{
		"trace":    "repository.GetUserByEmailAndRole",
		"ctx":      ctx,
		"email":    email,
		"roleName": roleName,
	})

	query := `
        SELECT u.id, u.password, u.email, r.id, r.name
        FROM users u
        JOIN roles r ON u.role_id = r.id AND u.deleted_at IS NULL
        WHERE u.email = $1 AND r.name = $2 AND u.deleted_at IS NULL
    `
	var user model.User
	user.Role = &model.Role{}
	err := u.DB.QueryRowContext(ctx, query, email, roleName).
		Scan(&user.ID, &user.Password, &user.Email, &user.Role.ID, &user.Role.Name)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil // No user found
		}
		log.Error(err)
		return nil, err
	}
	return &user, nil
}
