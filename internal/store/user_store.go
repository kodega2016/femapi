package store

import (
	"database/sql"
	"errors"
	"time"

	"golang.org/x/crypto/bcrypt"
)

type password struct {
	plainText *string
	hash      []byte
}

func (p *password) Set(plaintestPassword string) error {
	hash, err := bcrypt.GenerateFromPassword([]byte(plaintestPassword), 12)
	if err != nil {
		return err
	}
	p.plainText = &plaintestPassword
	p.hash = hash
	return nil
}

func (p *password) Matches(plaintestPassword string) (bool, error) {
	err := bcrypt.CompareHashAndPassword(
		p.hash, []byte(plaintestPassword),
	)
	if err != nil {
		switch {
		case errors.Is(err, bcrypt.ErrMismatchedHashAndPassword):
			return false, nil
		default:
			return false, err // internal server error
		}
	}

	return true, nil
}

type User struct {
	ID           int         `json:"id"`
	Username     string      `json:"username"`
	Email        string      `json:"email"`
	PasswordHash password    `json:"-"`
	Bio          string      `json:"bio"`
	CreatedAt    time.Time   `json:"created_at"`
	UpdatedAt    time.Ticker `json:"updated_at"`
}

type PostgresUserStrore struct {
	db *sql.DB
}

func NewPostgresUserStore(db *sql.DB) *PostgresUserStrore {
	return &PostgresUserStrore{
		db: db,
	}
}

type UserStore interface {
	CreateUser(*User) error
	GetUserByUsername(username string) (*User, error)
	UpdateUser(*User) error
}

func (s *PostgresUserStrore) CreateUser(user *User) error {
	query := `
	INSERT INTO users(username,email,password_hash,bio)
	VALUES($1,$2,$3,$4)
	RETURNING id,created_at,updated_at
	`

	err := s.db.QueryRow(query, user.Username, user.Email, user.PasswordHash, user.Bio).Scan(&user.ID, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		return err
	}

	return nil
}

func (s *PostgresUserStrore) GetUserByUsername(username string) (*User, error) {
	user := &User{
		PasswordHash: password{},
	}

	query := `
	SELECT id,username,email,password_hash,bio,created_at,updated_at
	FROM users
	WHERE username=$1
	`

	err := s.db.QueryRow(query, username).Scan(&user.ID, &user.Username, &user.Email, &user.PasswordHash.hash, &user.Bio, &user.CreatedAt, &user.UpdatedAt)

	if err == sql.ErrNoRows {
		return nil, nil
	}

	if err != nil {
		return nil, err
	}

	return user, nil
}

func (s *PostgresUserStrore) UpdateUser(user *User) error {
	query := `
	UPDATE users
	SET username=$1, email=$2,bio=$3,updated_at=CURRENT_TIMESTAMP
	RETURNING updated_at
	`

	result, err := s.db.Exec(query, user.Username, user.Email, user.Bio)
	if err != nil {
		return err
	}

	rowsEffected, err := result.RowsAffected()
	if rowsEffected == 0 {
		return sql.ErrNoRows
	}

	return nil
}
