package auth

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

var ErrUserNotFound = errors.New("user not found")
var ErrUserAlreadyExists = errors.New("user already exists")

var ErrInvalidCredentials = errors.New("invalid credentials")

type Role string

const (
	RoleAdmin Role = "admin"
	RoleUser  Role = "user"
)

type User struct {
	ID            uuid.UUID `db:"id"`
	Email         string    `db:"email"`
	EmailVerified bool      `db:"email_validated"`
	Password      string    `db:"password"`
	Role          Role      `db:"role"`
	CreatedAt     time.Time `db:"created_at"`
	UpdatedAt     time.Time `db:"updated_at"`
}

// UserRepository is the interface that wraps the basic user repository methods
type UserRepository interface {
	// FindByEmail will find user by email
	FindByEmail(email string) (User, error)
	// FindByID will find user by id
	FindByID(id uuid.UUID) (User, error)
	// Create will create a new user
	Create(email string, password string, role Role) error
	// Update will update user
	Update(user User) error
	// Delete will delete user
	Delete(id uuid.UUID) error

	//Authenticate will authenticate user
	Authenticate(email string, password string) (User, error)
}
