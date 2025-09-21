package users

import "backendAPI/db"

// Структуры

type RegistrationData struct {
	Fio      string `json:"fio"`
	Email    string `json:"email"`
	Password string `json:"password"`
	Avatar   string `json:"avatar"`
}

type ValidationErrorResponse struct {
	Message string            `json:"message"`
	Errors  map[string]string `json:"errors"`
}

type SignInData struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type ProfileData struct {
	User db.UserData `json:"user"`
}

type Order struct {
	IdO int `json:"id"`
	Products []int `json:"products"`
	Order_price int `json:"order_price"`
}

// Чувствительные данные

type authorizedUsers struct {
	token map[int]int
}

func NewAuthorizedUsers() *authorizedUsers {
	return &authorizedUsers{
		token: make(map[int]int, 5),
	}
}

func (a *authorizedUsers) AddToken(id int, token int) {
	a.token[id] = token
}

func (a *authorizedUsers) RemoveToken(id int) {
	delete(a.token, id)
}
