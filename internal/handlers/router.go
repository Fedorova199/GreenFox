package handlers

import (
	"errors"
	"net/http"

	"github.com/Fedorova199/GreenFox/internal/interfaces"
	"github.com/Fedorova199/GreenFox/internal/models"
	"github.com/Fedorova199/GreenFox/internal/service"
	"github.com/go-chi/chi/v5"
)

type Handler struct {
	*chi.Mux
	baseURL             string
	user                interfaces.User
	order               interfaces.Order
	withdrawal          interfaces.Withdrawal
	cookieAuthenticator interfaces.CookieAuthenticator
	pointAccrualService interfaces.PointAccrualService
	authenticator       interfaces.Middleware
}

func NewHandler(
	baseURL string,
	user interfaces.User,
	order interfaces.Order,
	withdrawal interfaces.Withdrawal,
	cookieAuthenticator interfaces.CookieAuthenticator,
	pointAccrualService interfaces.PointAccrualService,
	authenticator interfaces.Middleware,
	middlewares []interfaces.Middleware,
) *Handler {
	h := &Handler{
		Mux:                 chi.NewMux(),
		baseURL:             baseURL,
		user:                user,
		order:               order,
		withdrawal:          withdrawal,
		cookieAuthenticator: cookieAuthenticator,
		pointAccrualService: pointAccrualService,
	}

	h.Post("/api/user/register", Middlewares(h.Register, middlewares))
	h.Post("/api/user/login", Middlewares(h.Login, middlewares))

	h.Post("/api/user/orders", authenticator.Handle(Middlewares(h.CreateOrder, middlewares)))
	h.Get("/api/user/orders", authenticator.Handle(Middlewares(h.GetOrders, middlewares)))
	h.Get("/api/user/balance", authenticator.Handle(Middlewares(h.GetBalance, middlewares)))
	h.Post("/api/user/balance/withdraw", authenticator.Handle(Middlewares(h.Withdraw, middlewares)))
	h.Get("/api/user/balance/withdrawals", authenticator.Handle(Middlewares(h.GetWithdrawals, middlewares)))

	return h
}

func Middlewares(handler http.HandlerFunc, middlewares []interfaces.Middleware) http.HandlerFunc {
	for _, middleware := range middlewares {
		handler = middleware.Handle(handler)
	}

	return handler
}

func (h *Handler) getAuthUser(r *http.Request) (models.User, error) {
	login, ok := service.LoginFromContext(r.Context())
	if !ok {
		return models.User{}, errors.New("unauthorized")
	}

	user, err := h.user.GetByLogin(r.Context(), login)
	if err != nil {
		return models.User{}, err
	}

	return user, nil
}
