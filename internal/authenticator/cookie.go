package authenticator

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"crypto/sha512"
	"encoding/hex"
	"errors"
	"fmt"
	"net/http"
	"strconv"

	middleware "github.com/Fedorova199/GreenFox/internal/middlewares"
	"github.com/theplant/luhn"
)

type CookieAuth interface {
	SetCookie(w http.ResponseWriter, login string) error
}

type CookieAuthenticator struct {
	secret []byte
}

func NewCookieAuthenticator(secret []byte) *CookieAuthenticator {
	return &CookieAuthenticator{secret: secret}
}

func (a *CookieAuthenticator) GetLogin(r *http.Request) (string, error) {
	userCookie, err := r.Cookie("user_id")
	if err != nil {
		return "", err
	}

	signCookie, err := r.Cookie("sign")
	if err != nil {
		return "", err
	}

	h := hmac.New(sha256.New, a.secret)
	h.Write([]byte(userCookie.Value))
	calculatedSign := h.Sum(nil)
	sign, err := hex.DecodeString(signCookie.Value)
	if err != nil {
		return "", err
	}

	if !hmac.Equal(calculatedSign, sign) {
		return "", fmt.Errorf("wrong sign")
	}

	return userCookie.Value, nil
}

func (a *CookieAuthenticator) SetCookie(w http.ResponseWriter, login string) error {
	h := hmac.New(sha256.New, a.secret)
	_, err := h.Write([]byte(login))
	if err != nil {
		return err
	}

	sign := hex.EncodeToString(h.Sum(nil))
	userIDCookie := &http.Cookie{
		Name:  "user_id",
		Value: login,
	}
	signCookie := &http.Cookie{
		Name:  "sign",
		Value: sign,
	}

	http.SetCookie(w, userIDCookie)
	http.SetCookie(w, signCookie)

	return nil
}

func Hash(password string) string {
	s := sha512.New()
	s.Write([]byte(password))
	return hex.EncodeToString(s.Sum(nil))
}

func CheckOrderNumber(number string) error {
	orderInt, err := strconv.Atoi(number)
	if err != nil {
		return err
	}

	if !luhn.Valid(orderInt) {
		return errors.New("invalid order number")
	}

	return nil
}

func LoginFromContext(ctx context.Context) (string, bool) {
	u, ok := ctx.Value(middleware.ContextLoginKey).(string)
	return u, ok
}
