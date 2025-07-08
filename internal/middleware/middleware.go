// Package middleware is used for middleware
package middleware

import (
	"context"
	"net/http"

	"github.com/kodega2016/femapi/internal/store"
)

type UserMiddleware struct {
	UserStore store.UserStore
}

type ContextKey string

const UserContextKey = ContextKey("user")

func SetUser(r *http.Request, user *store.User) *http.Request {
	ctx := context.WithValue(r.Context(), UserContextKey, user)
	return r.WithContext(ctx)
}

func GetUser(r *http.Request) *store.User {
	// retrieve the value from request context and assert the type to *store.User
	user, ok := r.Context().Value(UserContextKey).(*store.User)
	if !ok {
		panic("missing user in request")
	}
	return user
}
