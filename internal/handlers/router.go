package handlers

import (
	"net/http"

	"github.com/Fedorova199/GreenFox/internal/interfaces"
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
		authenticator:       authenticator,
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
