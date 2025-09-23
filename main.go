package main

import (
	"backendAPI/db"
	"backendAPI/users"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
)

func main() {
	db.CreateDatabase()
	router := mux.NewRouter()

	router.Path("/signup").Methods("POST").HandlerFunc(users.SignUp)
	router.Path("/signin").Methods("POST").HandlerFunc(users.SignIn)
	router.Path("/logout").Methods("GET").HandlerFunc(users.Logout)
	router.Path("/products").Methods("GET").HandlerFunc(users.ViewProduct)
	router.Path("/profile").Methods("GET").HandlerFunc(users.ViewProfile)
	router.Path("/cart/{product_id}").Methods("POST").HandlerFunc(users.AddingProduct)
	router.Path("/cart").Methods("GET").HandlerFunc(users.ViewCart)
	router.Path("/cart/{idC}").Methods("DELETE").HandlerFunc(users.DeleteProduct)
	router.Path("/order").Methods("POST").HandlerFunc(users.PlacingOrder)
	router.Path("/order").Methods("GET").HandlerFunc(users.ViewOrder)
	router.Path("/profile").Methods("PATCH").HandlerFunc(users.EditProfile)
	
	router.Path("/product").Methods("POST").HandlerFunc(users.AddProductAdmin)
	router.Path("/product/{id}").Methods("DELETE").HandlerFunc(users.AddProductAdmin)
	router.Path("/product/{id}").Methods("PATCH").HandlerFunc(users.EditProductAdmin)

	router.Path("/test").Methods("OPTIONS").HandlerFunc(users.TEST)
	if err := http.ListenAndServe(":8080", router); err != nil {
		fmt.Println("Error:", err.Error())
	}
}
