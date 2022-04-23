package models

import (
	"encoding/json"
	"time"
)

var (
	NEW        = "NEW"
	PROCESSING = "PROCESSING"
	INVALID    = "INVALID"
	PROCESSED  = "PROCESSED"
)

type Accrual struct {
	Order   string  `json:"order"`
	Status  string  `json:"status"`
	Accrual float64 `json:"accrual"`
}

type Credentials struct {
	Login    string
	Password string
}

type User struct {
	ID           uint64  `json:"-"`
	Login        string  `json:"-"`
	PasswordHash string  `json:"-"`
	Balance      float64 `json:"current"`
	Withdrawn    float64 `json:"withdrawn"`
}

type Order struct {
	ID         uint64    `json:"-"`
	Number     string    `json:"number"`
	Status     string    `json:"status"`
	Accrual    float64   `json:"accrual,omitempty"`
	UploadedAt time.Time `json:"uploaded_at"`
	UserID     uint64    `json:"-"`
}

type Withdrawal struct {
	ID          uint64    `json:"-"`
	Order       string    `json:"order"`
	Sum         float64   `json:"sum"`
	ProcessedAt time.Time `json:"processed_at"`
	UserID      uint64    `json:"-"`
}

func (w Withdrawal) MarshalJSON() ([]byte, error) {
	type WithdrawalAlias Withdrawal

	aliasValue := struct {
		WithdrawalAlias
		ProcessedAt string `json:"processed_at"`
	}{
		WithdrawalAlias: WithdrawalAlias(w),
		ProcessedAt:     w.ProcessedAt.Format(time.RFC3339),
	}

	return json.Marshal(aliasValue)
}

func (o Order) MarshalJSON() ([]byte, error) {
	type OrderAlias Order

	aliasValue := struct {
		OrderAlias
		UploadedAt string `json:"uploaded_at"`
	}{
		OrderAlias: OrderAlias(o),
		UploadedAt: o.UploadedAt.Format(time.RFC3339),
	}

	return json.Marshal(aliasValue)
}
