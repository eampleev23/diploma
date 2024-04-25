package models

import "time"

// Order - модель заказа.
type Order struct {
	UploadedAt time.Time `json:"uploaded_at"`
	Number     string    `json:"number"`
	Status     string    `json:"status"`
	ID         int       `json:"id"`
	CustomerID int       `json:"customer_id"`
	Accrual    float64   `json:"accrual"`
}

// ResponseGetOwnerOrders описывает элемент ответа пользователю на получение всех его ссылок.
type ResponseGetOwnerOrders struct {
	UploadedAt time.Time `json:"uploadedAt"`
	Number     string    `json:"number"`
	Status     string    `json:"status"`
	Accrual    float64   `json:"accrual"`
}

func GetResponseGetOwnerOrders(source []Order) (result []ResponseGetOwnerOrders, err error) {
	result = make([]ResponseGetOwnerOrders, 0, len(source))
	for _, v := range source {
		result = append(result, ResponseGetOwnerOrders{
			Number:     v.Number,
			Status:     v.Status,
			Accrual:    v.Accrual,
			UploadedAt: v.UploadedAt,
		})
	}
	return result, nil
}

// ResponseGetBalance описывает элемент ответа пользователю на получение суммы его баллов и суммы списаний.
type ResponseGetBalance struct {
	Current   float64 `json:"current"`
	Withdrawn float64 `json:"withdrawn"`
}

func GetResponseBalance(current, withdrawn float64) (resp ResponseGetBalance, err error) {
	resp = ResponseGetBalance{
		Current:   current,
		Withdrawn: withdrawn,
	}
	return resp, nil
}
