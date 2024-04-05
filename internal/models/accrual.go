package models

import "time"

// MyApiError — описание ошибки при неверном запросе.
type MyApiError struct {
	Code      int       `json:"code"`
	Message   string    `json:"message"`
	Timestamp time.Time `json:"timestamp"`
}

// OrderAccrual - модель заказа.
type OrderAccrual struct {
	Order string `json:"order"`
	Goods []Good
}

// Good - модель товара.
type Good struct {
	Description string `json:"description"`
	Price       int    `json:"price"`
}

// OrderAccrualResp - модель ответа на запрос.
type OrderAccrualResp struct {
	Order   string `json:"order"`
	Status  string `json:"status"`
	Accrual int    `json:"accrual"`
}
