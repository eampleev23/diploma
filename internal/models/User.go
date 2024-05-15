package models

// User - модель пользователя.
type User struct {
	Login    string `json:"login"`
	Password string `json:"password"`
	ID       int    `json:"id"`
}

// UserRegReq - модель запроса на регистрацию.
type UserRegReq struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

// UserLoginReq - модель запроса на авторизацию.
type UserLoginReq struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}
