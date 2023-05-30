package middleware

import (
	"context"
	"net/http"
)

type ContextKey string

const ContextLoginKey ContextKey = "loginKey"

type CookieAuthenticatorChecker interface {
	GetLogin(r *http.Request) (string, error)
}

type Authenticator struct {
	cookieAuthenticator CookieAuthenticatorChecker
}

func NewAuthenticator(cookieAuthenticator CookieAuthenticatorChecker) *Authenticator {
	return &Authenticator{cookieAuthenticator: cookieAuthenticator}
}

func (a Authenticator) Handle(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		login, err := a.cookieAuthenticator.GetLogin(r)
		if err != nil {
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}

		ctx := context.WithValue(r.Context(), ContextLoginKey, login)
		next.ServeHTTP(w, r.WithContext(ctx))
	}
}
