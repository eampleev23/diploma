package models

// Order - модель заказа
type Order struct {
	ID     int    `json:"id"`
	Number string `json:"number"`
	UserId int    `json:"user_id"`
}
