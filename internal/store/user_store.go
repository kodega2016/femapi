package store

import (
	"database/sql"
	"time"
)

type password struct {
	plainText *string
	hash      []byte
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
