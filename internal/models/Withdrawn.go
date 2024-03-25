package models

import (
	"log"
	"time"
)

// Withdrawn - модель списания.
type Withdrawn struct {
	Order       string    `json:"order"`
	ID          int       `json:"id"`
	Sum         int       `json:"sum"`
	UserID      int       `json:"user_id"`
	ProcessedAt time.Time `json:"processed_at"`
}

// ResponseGetOwnerWithdrawals описывает элемент ответа пользователю на получение всех его списаний.
type ResponseGetOwnerWithdrawals struct {
	Order       string    `json:"order"`
	Sum         int       `json:"sum"`
	ProcessedAt time.Time `json:"processed_at"`
}

func GetResponseGetOwnerWithdrawals(source []Withdrawn) (result []ResponseGetOwnerWithdrawals, err error) {
	result = make([]ResponseGetOwnerWithdrawals, 0, len(source))
	for _, v := range source {
		log.Println("v.ProcessedAt", v.ProcessedAt)
		result = append(result, ResponseGetOwnerWithdrawals{
			Order:       v.Order,
			Sum:         v.Sum,
			ProcessedAt: v.ProcessedAt,
		})
	}
	return result, nil
}
