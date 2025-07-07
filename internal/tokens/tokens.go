// Package tokens is used to generate the auth token
package tokens

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/base32"
	"time"
)

const (
	ScopeAuth = "authentication"
)

type Token struct {
	Plaintext string    `json:"token"`
	Hash      []byte    `json:"-"`
	UserID    int       `json:"-"`
	Expiry    time.Time `json:"expiry"`
	Scope     string    `json:"-"`
}

func GenerateToken(UserID int, ttl time.Duration, scope string) (*Token, error) {
	token := &Token{
		UserID: UserID,
		Expiry: time.Now().Add(ttl),
	}

	// make bytes address to hold the hash
	emptyBytes := make([]byte, 32)

	// fills the bytes with the secure random bytes
	_, err := rand.Read(emptyBytes)
	if err != nil {
		return nil, err
	}

	token.Plaintext = base32.StdEncoding.WithPadding(base32.NoPadding).EncodeToString(emptyBytes)
	hash := sha256.Sum224([]byte(token.Plaintext))
	// convert this into the flexible slice type
	token.Hash = hash[:]
	return token, nil
}
