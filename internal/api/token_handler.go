package api

import (
	"log"

	"github.com/kodega2016/femapi/internal/store"
)

type TokenHandler struct {
	tokenStore store.TokenStore
	logger     *log.Logger
}

func NewTokenHandler(tokenStore store.TokenStore, logger *log.Logger) *TokenHandler {
	return &TokenHandler{
		tokenStore: tokenStore,
		logger:     logger,
	}
}
