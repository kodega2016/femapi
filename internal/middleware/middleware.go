// Package middleware is used for middleware
package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/kodega2016/femapi/internal/store"
	"github.com/kodega2016/femapi/internal/tokens"
	"github.com/kodega2016/femapi/internal/utils"
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

func (um *UserMiddleware) Authenticate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Vary", "Authorization")
		authHeader := r.Header.Get("Authorization")

		if authHeader == "" {
			r = SetUser(r, store.AnonymousUser)
			next.ServeHTTP(w, r)
			return
		}

		headerParts := strings.Split(authHeader, " ")
		if len(headerParts) != 2 || headerParts[0] != "Bearer" {
			utils.WriteJSON(w, http.StatusUnauthorized, utils.Envelope{
				"error": "invalid authorization",
			})
			return
		}

		token := headerParts[1]
		user, err := um.UserStore.GetUserToken(tokens.ScopeAuth, token)
		if err != nil {
			utils.WriteJSON(w, http.StatusUnauthorized, utils.Envelope{
				"error": "invalid authorization",
			})
			return
		}
		if user == nil {
			utils.WriteJSON(w, http.StatusUnauthorized, utils.Envelope{
				"error": "token expired or invalid",
			})
			return
		}
		r = SetUser(r, user)
		next.ServeHTTP(w, r)
	})
}

func (um *UserMiddleware) RequireUser(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user := GetUser(r)

		if user.IsAnonymous() {
			utils.WriteJSON(w, http.StatusUnauthorized, utils.Envelope{
				"error": "you must be logged in to access this route",
			})
			return
		}
		next.ServeHTTP(w, r)
	})
}
