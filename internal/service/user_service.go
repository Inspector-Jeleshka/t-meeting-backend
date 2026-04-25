package service

import (
	"context"
	"errors"
	"fmt"
	"t-meeting-backend/internal/domain"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type userRepository interface {
	Create(ctx context.Context, user *domain.User) error
	GetByEmail(ctx context.Context, email string) (*domain.User, error)
	GetByID(ctx context.Context, id uuid.UUID) (*domain.User, error)
}

type UserService struct {
	ur userRepository
}

func NewUserService(ur userRepository) *UserService {
	return &UserService{ur: ur}
}

func (us *UserService) Register(ctx context.Context, email, password string) (*domain.User, error) {
	// Проверка: зарегистрирован ли уже пользователь с этим email
	_, err := us.ur.GetByEmail(ctx, email)
	if err == nil {
		return nil, domain.ErrUserAlreadyExists
	}

	if !errors.Is(err, domain.ErrUserNotFound) {
		return nil, fmt.Errorf("get user by email: %w", err)
	}

	passwordHash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		// bcrypt.GenerateFromPassword возвращает два типа ошибок:
		// Первые - клиентские ошибки (4xx) - ErrPasswordTooLong
		// Вторые - внутренние ошибки сервера (5xx) - ErrUnexpectedEOF, CorruptInputError, KeySizeError
		return nil, fmt.Errorf("generate password hash: %w", err)
	}

	id := uuid.New()
	role := domain.AdminRole
	user := &domain.User{
		ID:           id,
		Email:        email,
		PasswordHash: string(passwordHash),
		Role:         role,
	}
	if err = us.ur.Create(ctx, user); err != nil {
		return nil, fmt.Errorf("create user: %w", err)
	}
	return user, nil
}

func (us *UserService) GetByID(ctx context.Context, id uuid.UUID) (*domain.User, error) {
	return us.ur.GetByID(ctx, id)
}

func (us *UserService) Login(ctx context.Context, email, password string) (*domain.User, error) {
	user, err := us.ur.GetByEmail(ctx, email)
	if err != nil {
		if errors.Is(err, domain.ErrUserNotFound) {
			return nil, domain.ErrInvalidCredentials
		}
		return nil, fmt.Errorf("get user by email: %w", err)
	}

	if err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)); err != nil {
		return nil, domain.ErrInvalidCredentials
	}
	return user, nil
}
