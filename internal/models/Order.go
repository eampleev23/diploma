package models

// Order - модель заказа
type Order struct {
	ID     int    `json:"id"`
	Number string `json:"number"`
	UserID int    `json:"user_id"`
}
