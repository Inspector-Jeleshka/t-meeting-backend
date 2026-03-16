package domain

import (
	"errors"

	"github.com/google/uuid"
)

var ErrUserNotFound = errors.New("user not found")
var ErrInvalidCredentials = errors.New("invalid user credentials")

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
