package users

import (
	"backendAPI/db"
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
		http.Error(w, `{"error": "not Admin"}`, http.StatusForbidden)
		return
	}

	var req Product
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, `{"message": "Invalid product data"}`, http.StatusBadRequest)
		return
	}
	if req.Name == "" || req.Description == "" || req.Price == 0 {
		http.Error(w, `{"error": "Invalid product data"}`, http.StatusBadRequest)
	}

	db := db.ConnectDB()
	idP := db.AddProductAdmin(req.Name, req.Description, req.Price)
	db.CloseDB()

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
		http.Error(w, `{"error": "not Admin"}`, http.StatusForbidden)
		return
	}

	idP, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		http.Error(w, `{"error": "Forbidden for you"}`, http.StatusForbidden)
		return
	}
	db := db.ConnectDB()
	db.DeleteProductAdmin(idP)
	db.CloseDB()

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
		http.Error(w, `{"error": "not Admin"}`, http.StatusForbidden)
		return
	}
	idP, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		http.Error(w, `{"error": "Forbidden for you"}`, http.StatusForbidden)
		return
	}

	var req Product
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, `{"message": "Invalid product data"}`, http.StatusBadRequest)
		return
	}

	db := db.ConnectDB()
	defer db.CloseDB()
	if req.Name != "" {
		if _, err := db.Exec("UPDATE products SET name = ? WHERE idP = ?", req.Name, idP); err != nil {
			panic(err)
		}
	}
	if req.Description != "" {
		if _, err := db.Exec("UPDATE products SET description = ? WHERE idP = ?", req.Description, idP); err != nil {
			panic(err)
		}
	}
	if req.Price != 0 {
		if _, err := db.Exec("UPDATE products SET price = ? WHERE idP = ?", req.Price, idP); err != nil {
			panic(err)
		}
	}

	w.Header().Set("Content-Type", "application/json")
	product := db.SelectProduct(idP)
	if err := json.NewEncoder(w).Encode(product); err != nil {
		panic(err)
	}
}
