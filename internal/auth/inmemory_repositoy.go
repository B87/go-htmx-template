package auth

import (
	"sync"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type InMemmoryUserRepository struct {
	users map[uuid.UUID]User
	mu    sync.RWMutex
}

func NewInMemmoryUserRepository() *InMemmoryUserRepository {
	defPass, err := hashPassword("123")
	if err != nil {
		panic(err)
	}
	return &InMemmoryUserRepository{
		users: map[uuid.UUID]User{
			uuid.New(): {
				Email:    "test@gmail.com",
				Password: defPass,
			},
		},
		mu: sync.RWMutex{},
	}
}

func (r *InMemmoryUserRepository) FindByEmail(email string) (User, error) {
	for _, u := range r.users {
		if u.Email == email {
			return u, nil
		}
	}
	return User{}, ErrUserNotFound
}

func (r *InMemmoryUserRepository) FindByID(id uuid.UUID) (User, error) {
	return r.users[id], nil
}

func (r *InMemmoryUserRepository) Create(email string, password string, role Role) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	_, err := r.FindByEmail(email)
	if err == nil {
		return ErrUserAlreadyExists
	}
	id, err := uuid.NewV7()
	if err != nil {
		return err
	}
	hashedPassword, err := hashPassword(password)
	if err != nil {
		return err
	}
	user := User{
		ID:       id,
		Email:    email,
		Password: hashedPassword,
	}
	r.users[user.ID] = user
	return nil
}

func (r *InMemmoryUserRepository) Update(user User) error {
	r.users[user.ID] = user
	return nil
}

func (r *InMemmoryUserRepository) Delete(id uuid.UUID) error {
	delete(r.users, id)
	return nil
}

func (r *InMemmoryUserRepository) Authenticate(email string, password string) (User, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	user, err := r.FindByEmail(email)
	if err == ErrUserNotFound {
		return User{}, ErrUserNotFound
	}
	if err != nil {
		return User{}, err
	}
	if !checkPasswordHash(password, user.Password) {
		return User{}, ErrInvalidCredentials
	}
	return user, nil
}

func hashPassword(password string) (string, error) {
	// The cost parameter specifies the computational difficulty of the hashing:
	// higher values are more secure but more time-consuming.
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}

func checkPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
