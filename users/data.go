package users

import (
	"backendAPI/db"
	"encoding/json"
	"os"
)

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
	IdO         int   `json:"id"`
	Products    []int `json:"products"`
	Order_price int   `json:"order_price"`
}

type Product struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Price       int    `json:"price"`
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

type admin struct {
	IDs []int
}

func NewAdmins() *admin {
	idAdmins := make([]int, 2)

	file, err := os.Open("users/admins.json")
	if err != nil {
		panic(err)
	}
	defer file.Close()

	if err := json.NewDecoder(file).Decode(&idAdmins); err != nil {
		panic(err)
	}

	return &admin{
		IDs: idAdmins,
	}
}
