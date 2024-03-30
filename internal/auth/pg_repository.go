package auth

import (
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

type PostgresUserRepository struct {
	db *sqlx.DB
}

func NewPostgresUserRepository(db *sqlx.DB) *PostgresUserRepository {
	return &PostgresUserRepository{
		db: db,
	}
}

func (r *PostgresUserRepository) FindByEmail(email string) (User, error) {
	var user User
	err := r.db.Get(&user, "SELECT * FROM users WHERE email = $1", email)
	return user, err
}

func (r *PostgresUserRepository) FindByID(id uuid.UUID) (User, error) {
	var user User
	err := r.db.Get(&user, "SELECT * FROM users WHERE id = $1", id)
	return user, err
}

func (r *PostgresUserRepository) Create(email string, password string, role Role) error {
	hashedPwd, err := hashPassword(password)
	if err != nil {
		return err
	}
	_, err = r.db.Exec("INSERT INTO users (email, password, role) VALUES ($1, $2, $3)", email, hashedPwd, role)
	return err
}

func (r *PostgresUserRepository) Update(user User) error {
	_, err := r.db.Exec("UPDATE users SET email = $1, password = $2, role = $3 WHERE id = $4", user.Email, user.Password, user.Role, user.ID)
	return err
}

func (r *PostgresUserRepository) Delete(id uuid.UUID) error {
	_, err := r.db.Exec("DELETE FROM users WHERE id = $1", id)
	return err
}

func (r *PostgresUserRepository) Authenticate(email string, password string) (User, error) {
	var user User
	err := r.db.Get(&user, "SELECT * FROM users WHERE email = $1", email)

	if !checkPasswordHash(password, user.Password) {
		return User{}, ErrInvalidCredentials
	}

	return user, err
}
