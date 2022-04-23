package service

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/Fedorova199/GreenFox/internal/models"
)

type OrderRepository interface {
	UpdateAccrualStatus(ctx context.Context, accrual models.Accrual) error
}

type PointAccrualService struct {
	orders               chan string
	accrualSystemAddress string
	orderRepository      OrderRepository
}

func NewPointAccrualService(accrualSystemAddress string, orderRepository OrderRepository) *PointAccrualService {
	return &PointAccrualService{
		orders:               make(chan string, 100),
		accrualSystemAddress: accrualSystemAddress,
		orderRepository:      orderRepository,
	}
}

func (s *PointAccrualService) Start() {
	go func() {
		for order := range s.orders {
			err := s.handleOrder(order)
			if err != nil {
				s.Accrue(order)
			}
		}
	}()
}

func (s *PointAccrualService) handleOrder(order string) error {
	url := fmt.Sprintf("%s/api/orders/%s", s.accrualSystemAddress, order)
	response, err := http.Get(url)
	if err != nil {
		return err
	}

	switch response.StatusCode {
	case http.StatusOK:
		defer response.Body.Close()
		payload, err := io.ReadAll(response.Body)
		if err != nil {
			return err
		}

		accrual := models.Accrual{}
		if err := json.Unmarshal(payload, &accrual); err != nil {
			return err
		}

		err = s.orderRepository.UpdateAccrualStatus(context.Background(), accrual)
		if err != nil {
			return err
		}
	case http.StatusTooManyRequests:
		s.Accrue(order)
	case http.StatusInternalServerError:
		s.Accrue(order)
	}

	return nil
}

func (s *PointAccrualService) Stop() {
	close(s.orders)
}

func (s *PointAccrualService) Accrue(order string) {
	s.orders <- order
}
