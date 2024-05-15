package models

import (
	"time"
)

// MyAPIError — описание ошибки при неверном запросе.
type MyAPIError struct {
	Timestamp time.Time `json:"timestamp"`
	Message   string    `json:"message"`
	Code      int       `json:"code"`
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
	Order   string  `json:"order"`
	Status  string  `json:"status"`
	Accrual float64 `json:"accrual"`
}
