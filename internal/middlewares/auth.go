package middleware

import (
	"net/http"
)

type CookieAuthenticatorChecker interface {
	Check(r *http.Request) error
}

type Authenticator struct {
	cookieAuthenticator CookieAuthenticatorChecker
}

func NewAuthenticator(cookieAuthenticator CookieAuthenticatorChecker) *Authenticator {
	return &Authenticator{cookieAuthenticator: cookieAuthenticator}
}

func (a Authenticator) Handle(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		err := a.cookieAuthenticator.Check(r)
		if err != nil {
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}

		next.ServeHTTP(w, r)
	}
}
