package models

import (
	"time"
)

// Withdrawn - модель списания.
type Withdrawn struct {
	ProcessedAt time.Time `json:"processed_at"`
	Order       string    `json:"order"`
	ID          int       `json:"id"`
	Sum         float64   `json:"sum"`
	UserID      int       `json:"user_id"`
}

// ResponseGetOwnerWithdrawals описывает элемент ответа пользователю на получение всех его списаний.
type ResponseGetOwnerWithdrawals struct {
	ProcessedAt time.Time `json:"processed_at"`
	Order       string    `json:"order"`
	Sum         float64   `json:"sum"`
}

func GetResponseGetOwnerWithdrawals(source []Withdrawn) (result []ResponseGetOwnerWithdrawals, err error) {
	result = make([]ResponseGetOwnerWithdrawals, 0, len(source))
	for _, v := range source {
		result = append(result, ResponseGetOwnerWithdrawals{
			Order:       v.Order,
			Sum:         v.Sum,
			ProcessedAt: v.ProcessedAt,
		})
	}
	return result, nil
}
