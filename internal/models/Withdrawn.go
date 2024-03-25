package models

// Withdrawn - модель списания.
type Withdrawn struct {
	Order  string `json:"order"`
	ID     int    `json:"id"`
	Sum    int    `json:"sum"`
	UserID int    `json:"user_id"`
}

// ResponseGetOwnerWithdrawals описывает элемент ответа пользователю на получение всех его списаний.
type ResponseGetOwnerWithdrawals struct {
	Order       string `json:"order"`
	Sum         string `json:"sum"`
	ProcessedAt int    `json:"processed_at"`
}

func GetResponseGetOwnerWithdrawals(source []Withdrawn) (result []ResponseGetOwnerWithdrawals, err error) {
	result = make([]ResponseGetOwnerWithdrawals, 0, len(source))
	//for _, v := range source {
	//result = append(result, ResponseGetOwnerOrders{
	//	Number:     v.Number,
	//	Status:     v.Status,
	//	Accrual:    v.Accrual,
	//	UploadedAt: v.UploadedAt,
	//})
	//}
	return result, nil
}
