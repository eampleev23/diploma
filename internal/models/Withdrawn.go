package models

// Withdrawn - модель списания.
type Withdrawn struct {
	Order  string `json:"order"`
	ID     int    `json:"id"`
	Sum    int    `json:"sum"`
	UserID int    `json:"user_id"`
}
