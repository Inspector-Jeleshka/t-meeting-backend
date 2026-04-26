package repository

import (
	"context"
	"errors"
	"fmt"
	"t-meeting-backend/internal/domain"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type UserRepository struct {
	db *pgxpool.Pool
}

func NewUserRepository(db *pgxpool.Pool) (*UserRepository, error) {
	if db == nil {
		return nil, errors.New("db is nil")
	}
	return &UserRepository{db: db}, nil
}

func (ur *UserRepository) Create(ctx context.Context, user *domain.User) error {
	commandTag, err := ur.db.Exec(ctx, qUserCreate, user.ID, user.Email, user.PasswordHash, user.Role)
	if err != nil {
		return fmt.Errorf("create user: %w", err)
	}

	if commandTag.RowsAffected() != 1 {
		return fmt.Errorf("%w: rows affected = %d", domain.ErrUserNotCreated, commandTag.RowsAffected())
	}
	return nil
}

func (ur *UserRepository) GetByEmail(ctx context.Context, email string) (*domain.User, error) {
	var user domain.User
	err := ur.db.QueryRow(ctx, qUserGetByEmail, email).Scan(
		&user.ID,
		&user.Email,
		&user.PasswordHash,
		&user.Role,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain.ErrUserNotFound
		}
		return nil, err
	}
	return &user, nil
}

func (ur *UserRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.User, error) {
	var user domain.User

	err := ur.db.QueryRow(ctx, qUserGetByID, id).Scan(
		&user.ID,
		&user.Email,
		&user.PasswordHash,
		&user.Role)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain.ErrUserNotFound
		}
		return nil, err
	}
	return &user, nil
}
