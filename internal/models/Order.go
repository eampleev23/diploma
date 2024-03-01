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
