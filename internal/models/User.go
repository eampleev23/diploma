package models

// User - модель пользователя
type User struct {
	ID       int    `json:"id"`
	Login    string `json:"login"`
	Password string `json:"password"`
}

// UserRegReq - модель запроса на регистрацию
type UserRegReq struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

// UserLoginReq - модель запроса на авторизацию
type UserLoginReq struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}
