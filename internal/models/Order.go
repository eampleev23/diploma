package models

// Order - модель заказа.
type Order struct {
	Number string `json:"number"`
	ID     int    `json:"id"`
	UserID int    `json:"user_id"`
}
