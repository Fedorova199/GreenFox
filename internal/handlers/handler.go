package handlers

import (
	"database/sql"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"sort"
	"strconv"
	"time"

	"github.com/Fedorova199/GreenFox/internal/models"
	"github.com/Fedorova199/GreenFox/internal/service"
	"github.com/theplant/luhn"
)

func (h *Handler) Register(w http.ResponseWriter, r *http.Request) {
	b, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	credentials := models.Credentials{}
	if err := json.Unmarshal(b, &credentials); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	_, err = h.user.GetByLogin(r.Context(), credentials.Login)
	if err == nil {
		http.Error(w, "login has already been taken", http.StatusConflict)
		return
	}

	newUser := models.User{
		Login:        credentials.Login,
		PasswordHash: service.Hash(credentials.Password),
	}

	err = h.user.Create(r.Context(), newUser)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (h *Handler) Login(w http.ResponseWriter, r *http.Request) {
	b, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	credentials := models.Credentials{}
	if err := json.Unmarshal(b, &credentials); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	user, err := h.user.GetByLogin(r.Context(), credentials.Login)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	if service.Hash(credentials.Password) != user.PasswordHash {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	err = h.cookieAuthenticator.SetCookie(w, credentials.Login)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (h *Handler) CreateOrder(w http.ResponseWriter, r *http.Request) {
	login, err := h.cookieAuthenticator.GetLogin(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	user, err := h.user.GetByLogin(r.Context(), login)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	b, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	numberInt, err := strconv.Atoi(string(b))
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if !luhn.Valid(numberInt) {
		http.Error(w, "invalid order number", http.StatusUnprocessableEntity)
		return
	}

	number := strconv.Itoa(numberInt)
	newOrder := models.Order{
		Number:     number,
		Status:     models.NEW,
		UploadedAt: time.Now(),
		UserID:     user.ID,
	}

	order, err := h.order.GetByNumber(r.Context(), number)
	if err != nil {
		if errors.As(err, &sql.ErrNoRows) {
			err = h.order.Create(r.Context(), newOrder)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			h.pointAccrualService.Accrue(newOrder.Number)

			w.WriteHeader(http.StatusAccepted)
			return
		} else {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}

	if order.UserID != user.ID {
		w.WriteHeader(http.StatusConflict)
		return
	}
	w.WriteHeader(http.StatusOK)
	return
}

func (h *Handler) GetOrders(w http.ResponseWriter, r *http.Request) {
	login, err := h.cookieAuthenticator.GetLogin(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	user, err := h.user.GetByLogin(r.Context(), login)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	orders, err := h.order.GetByUserID(r.Context(), user.ID)
	if err != nil {
		if errors.As(err, &sql.ErrNoRows) {
			w.WriteHeader(http.StatusNoContent)
			return
		} else {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}

	sort.Slice(orders, func(i, j int) bool {
		return orders[i].UploadedAt.Before(orders[j].UploadedAt)
	})

	res, err := json.Marshal(orders)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(res)
}

func (h *Handler) GetBalance(w http.ResponseWriter, r *http.Request) {
	login, err := h.cookieAuthenticator.GetLogin(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	user, err := h.user.GetByLogin(r.Context(), login)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	res, err := json.Marshal(user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(res)
}

func (h *Handler) GetWithdrawals(w http.ResponseWriter, r *http.Request) {
	login, err := h.cookieAuthenticator.GetLogin(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	user, err := h.user.GetByLogin(r.Context(), login)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	withdrawals, err := h.withdrawal.GetByUserID(r.Context(), user.ID)
	if err != nil {
		if errors.As(err, &sql.ErrNoRows) {
			w.WriteHeader(http.StatusNoContent)
			return
		} else {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}

	sort.Slice(withdrawals, func(i, j int) bool {
		return withdrawals[i].ProcessedAt.Before(withdrawals[j].ProcessedAt)
	})

	res, err := json.Marshal(withdrawals)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(res)
}

func (h *Handler) Withdraw(w http.ResponseWriter, r *http.Request) {
	login, err := h.cookieAuthenticator.GetLogin(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	user, err := h.user.GetByLogin(r.Context(), login)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	b, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var withdrawal models.Withdrawal
	err = json.Unmarshal(b, &withdrawal)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	orderInt, err := strconv.Atoi(withdrawal.Order)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if !luhn.Valid(orderInt) {
		http.Error(w, "invalid order number", http.StatusUnprocessableEntity)
		return
	}

	if user.Balance < withdrawal.Sum {
		http.Error(w, "insufficient balance", http.StatusPaymentRequired)
		return
	}

	err = h.user.DecreaseBalanceByUserID(r.Context(), user.ID, withdrawal.Sum)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	withdrawal.ProcessedAt = time.Now()
	withdrawal.UserID = user.ID
	err = h.withdrawal.Create(r.Context(), withdrawal)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	return
}
