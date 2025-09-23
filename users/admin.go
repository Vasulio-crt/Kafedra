package users

import (
	"backendAPI/db"
	"encoding/json"
	"fmt"
	"net/http"
)

var admins = NewAdmins()

// ---Работа с продуктами---

func AddProduct(w http.ResponseWriter, r *http.Request) {
	token := getTokenFromHeader(w, r)
	if token == -1 {
		return
	}
	id, ok := findKeyByValue(authorized.token, token)
	if !ok {
		http.Error(w, `{"error": "Invalid token"}`, http.StatusBadRequest)
		return
	}
	if !contains(admins.IDs, id) {
		http.Error(w, `{"error": "not Admin"}`, http.StatusForbidden)
		return
	}

	var req Product
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, `{"message": "Internal Server Error"}`, http.StatusInternalServerError)
		return
	}
	if req.Name == "" || req.Description == "" || req.Price == 0 {
		http.Error(w, `{"error": "Invalid product data"}`, http.StatusBadRequest)
	}

	db := db.ConnectDB()
	idP := db.AddProduct(req.Name, req.Description, req.Price)
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
