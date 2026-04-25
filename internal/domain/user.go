package domain

import (
	"errors"

	"github.com/google/uuid"
)

var ErrUserNotFound = errors.New("user not found")
var ErrInvalidCredentials = errors.New("invalid user credentials")
var ErrInvalidToken = errors.New("invalid token")
var ErrUserAlreadyExists = errors.New("user already exists")
var ErrUserNotCreated = errors.New("user not created")

type User struct {
	ID           uuid.UUID
	Email        string
	PasswordHash string
	Role         string
}

const (
	AdminRole     = "admin"
	OrganizerRole = "organizer"
	MemberRole    = "member"
	ViewerRole    = "viewer"
)
