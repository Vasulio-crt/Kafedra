package users

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

var admins = NewAdmins()

// func getIdAdmin(){}

// ---Работа с продуктами---

func AddProductAdmin(w http.ResponseWriter, r *http.Request) {
	id, ok := whatIdUser(r)
	if !ok {
		errorJSON(w, "Login failed", 403)
		return
	}
	if !contains(admins.IDs, id) {
		errorJSON(w, "not Admin", 403)
		return
	}

	var req Product
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		panic(err)
	}
	if req.Name == "" || req.Description == "" || req.Price == 0 {
		errorJSON(w, "Invalid product data", 400)
	}

	idP := DB.AddProductAdmin(req.Name, req.Description, req.Price)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

	mes := fmt.Sprintf(`{"product_id": %d, "message": "Product added"}`, idP)
	w.Write([]byte(mes))
}

func contains(i []int, id int) bool {
	for _, v := range i {
		if v == id {
			return true
		}
	}
	return false
}

func DeleteProductAdmin(w http.ResponseWriter, r *http.Request) {
	id, ok := whatIdUser(r)
	if !ok {
		errorJSON(w, "Login failed", 403)
		return
	}
	if !contains(admins.IDs, id) {
		errorJSON(w, "not Admin", 403)
		return
	}

	idP, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		errorJSON(w, "Forbidden for you", 403)
		return
	}

	DB.DeleteProductAdmin(idP)

	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(`{"message": "Product removed"}`))
}

func EditProductAdmin(w http.ResponseWriter, r *http.Request) {
	id, ok := whatIdUser(r)
	if !ok {
		errorJSON(w, "Login failed", 403)
		return
	}
	if !contains(admins.IDs, id) {
		errorJSON(w, "not Admin", 403)
		return
	}
	idP, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		errorJSON(w, "Forbidden for you", 403)
		return
	}

	var req Product
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		errorJSON(w, "Invalid product data", 400)
		return
	}

	if req.Name != "" {
		if _, err := DB.Exec("UPDATE products SET name = ? WHERE idP = ?", req.Name, idP); err != nil {
			panic(err)
		}
	}
	if req.Description != "" {
		if _, err := DB.Exec("UPDATE products SET description = ? WHERE idP = ?", req.Description, idP); err != nil {
			panic(err)
		}
	}
	if req.Price != 0 {
		if _, err := DB.Exec("UPDATE products SET price = ? WHERE idP = ?", req.Price, idP); err != nil {
			panic(err)
		}
	}

	product := DB.SelectProduct(idP)
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(product); err != nil {
		panic(err)
	}
}
