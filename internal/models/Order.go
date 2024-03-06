package models

import "time"

// Order - модель заказа.
type Order struct {
	Number     string    `json:"number"`
	ID         int       `json:"id"`
	CustomerID int       `json:"customer_id"`
	Status     string    `json:"status"`
	Accrual    int       `json:"accrual"`
	UploadedAt time.Time `json:"uploaded_at"`
}

// ResponseGetOwnerOrders описывает элемент ответа пользователю на получение всех его ссылок.
type ResponseGetOwnerOrders struct {
	Number     string    `json:"number"`
	Status     string    `json:"status"`
	Accrual    int       `json:"accrual"`
	UploadedAt time.Time `json:"uploadedAt"`
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
	Current   int `json:"current"`
	Withdrawn int `json:"withdrawn"`
}

func GetResponseBalance(current, withdrawn int) (resp ResponseGetBalance, err error) {
	resp = ResponseGetBalance{
		Current:   current,
		Withdrawn: withdrawn,
	}
	return resp, nil
}
