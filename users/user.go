package users

import (
	"backendAPI/db"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"time"

	"github.com/gorilla/mux"
)

var authorized = NewAuthorizedUsers()

// ---Работа с аккаунтами---

func isValidEmail(email string) bool {
	re := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	return re.MatchString(email)
}

func SignUp(w http.ResponseWriter, r *http.Request) {
	var req RegistrationData
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, `{"message": "Internal Server Error"}`, http.StatusInternalServerError)
		return
	}

	errors := make(map[string]string)
	if req.Fio == "" {
		errors["fio"] = "FIO is a required field"
	}
	if !isValidEmail(req.Email) {
		errors["email"] = "Invalid email format"
	}
	if len(req.Password) < 3 {
		errors["password"] = "Password must contain at least 3 characters"
	}

	if len(errors) > 0 {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnprocessableEntity)
		json.NewEncoder(w).Encode(ValidationErrorResponse{Message: "Validation error", Errors: errors})
		return
	}

	if req.Avatar == "" {
		req.Avatar = "avatars/default.jpeg"
	}

	db := db.ConnectDB()
	db.InsertUser(req.Fio, req.Email, req.Password, req.Avatar)
	idU := db.GetIdUser(req.Email)
	db.CloseDB()

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

	token := int(time.Now().UnixMicro())
	authorized.AddToken(idU, token)

	mes := fmt.Sprintf(`{"user_token": %d}`, token)
	w.Write([]byte(mes))
}

func SignIn(w http.ResponseWriter, r *http.Request) {
	var req SignInData
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, `{"message": "Internal Server Error"}`, http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")

	errors := make(map[string]string)

	if !isValidEmail(req.Email) {
		errors["email"] = "Invalid email format"
	}
	if len(req.Password) < 3 {
		errors["password"] = "Password must contain at least 3 characters"
	}
	if len(errors) > 0 {
		w.WriteHeader(http.StatusUnprocessableEntity)
		json.NewEncoder(w).Encode(ValidationErrorResponse{Message: "Validation error", Errors: errors})
		return
	}

	db := db.ConnectDB()
	idU, password := db.GetPassword(req.Email)
	db.CloseDB()

	if password != req.Password {
		http.Error(w, `{"message": "Login failed"}`, http.StatusUnauthorized)
		return
	}
	w.WriteHeader(http.StatusOK)

	token := int(time.Now().UnixMicro())
	authorized.AddToken(idU, token)

	mes := fmt.Sprintf(`{"user_token": %d}`, token)
	w.Write([]byte(mes))
}

func findKeyByValue(m map[int]int, targetValue int) (int, bool) {
	for key, value := range m {
		if value == targetValue {
			return key, true
		}
	}
	return 0, false
}

func getTokenFromHeader(w http.ResponseWriter, r *http.Request) int {
	token, err := strconv.Atoi(r.Header.Get("token"))
	if err != nil {
		http.Error(w, `{"error": "Invalid token"}`, http.StatusBadRequest)
		return -1
	}
	return token
}

func Logout(w http.ResponseWriter, r *http.Request) {
	token := getTokenFromHeader(w, r)
	if token == -1 {
		return
	}
	w.Header().Set("Content-Type", "application/json")

	id, ok := findKeyByValue(authorized.token, token)
	if ok {
		authorized.RemoveToken(id)
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"message": "logout"}`))
	} else {
		http.Error(w, `{"error": "Invalid token"}`, http.StatusBadRequest)
		return
	}
}

func EditProfile(w http.ResponseWriter, r *http.Request) {
	token := getTokenFromHeader(w, r)
	if token == -1 {
		return
	}
	id, ok := findKeyByValue(authorized.token, token)
	if !ok {
		http.Error(w, `{"error": "Invalid token"}`, http.StatusBadRequest)
		return
	}
	var req RegistrationData
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, `{"message": "Internal Server Error"}`, http.StatusInternalServerError)
		return
	}
	db := db.ConnectDB()
	defer db.CloseDB()
	if req.Avatar != "" {
		if _, err := db.Exec("UPDATE users SET avatar = ? WHERE id = ?", req.Avatar, id); err != nil {
			panic(err)
		}
	}
	if req.Fio != "" {
		if _, err := db.Exec("UPDATE users SET fio = ? WHERE id = ?", req.Fio, id); err != nil {
			panic(err)
		}
	}
	if req.Password != "" {
		if _, err := db.Exec("UPDATE users SET password = ? WHERE id = ?", req.Password, id); err != nil {
			panic(err)
		}
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(`{"message": "data updated successfully"}`))
}

// ---Просмотр---

func ViewProduct(w http.ResponseWriter, r *http.Request) {
	db := db.ConnectDB()
	products := db.GetProduct()
	db.CloseDB()

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(products); err != nil {
		http.Error(w, `{"message": "Internal Server Error"}`, http.StatusInternalServerError)
		return
	}
}

func ViewProfile(w http.ResponseWriter, r *http.Request) {
	token := getTokenFromHeader(w, r)
	if token == -1 {
		return
	}
	w.Header().Set("Content-Type", "application/json")

	id, ok := findKeyByValue(authorized.token, token)
	if !ok {
		http.Error(w, `{"error": "Invalid token"}`, http.StatusBadRequest)
		return
	}

	db := db.ConnectDB()
	user := db.GetUser(id)
	db.CloseDB()

	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(ProfileData{User: user}); err != nil {
		http.Error(w, `{"message": "Internal Server Error"}`, http.StatusInternalServerError)
		return
	}
}

// ---Корзина---

func AddingProduct(w http.ResponseWriter, r *http.Request) {
	product_id, err := strconv.Atoi(mux.Vars(r)["product_id"])
	if err != nil {
		http.Error(w, `{"error": "Invalid product ID"}`, http.StatusBadRequest)
		return
	}
	token := getTokenFromHeader(w, r)
	if token == -1 {
		return
	}

	db := db.ConnectDB()
	product := db.GetProductById(product_id)
	defer db.CloseDB()

	if product.IdProduct == 0 {
		http.Error(w, `{"error": "Invalid product ID"}`, http.StatusBadRequest)
		return
	}

	id, ok := findKeyByValue(authorized.token, token)
	if !ok {
		http.Error(w, `{"error": "Invalid token"}`, http.StatusBadRequest)
		return
	}

	db.AddToCart(id, product_id)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte(`{"message": "Product add to card"}`))
}

func ViewCart(w http.ResponseWriter, r *http.Request) {
	token := getTokenFromHeader(w, r)
	if token == -1 {
		return
	}
	id, ok := findKeyByValue(authorized.token, token)
	if !ok {
		http.Error(w, `{"error": "Invalid token"}`, http.StatusBadRequest)
		return
	}
	db := db.ConnectDB()
	cart := db.ViewCart(id)
	db.CloseDB()

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(cart); err != nil {
		http.Error(w, `{"message": "Internal Server Error"}`, http.StatusInternalServerError)
		return
	}
}

func DeleteProduct(w http.ResponseWriter, r *http.Request) {
	idC, err := strconv.Atoi(mux.Vars(r)["idC"])
	if err != nil {
		http.Error(w, `{"error": "Forbidden for you"}`, http.StatusForbidden)
		return
	}
	token := getTokenFromHeader(w, r)
	if token == -1 {
		return
	}
	id, ok := findKeyByValue(authorized.token, token)
	if !ok {
		http.Error(w, `{"error": "Invalid token"}`, http.StatusBadRequest)
		return
	}
	db := db.ConnectDB()
	db.DeleteProduct(idC, id)
	db.CloseDB()

	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(`{"message": "Item removed from cart"}`))
}

// ---Заказ---

func PlacingOrder(w http.ResponseWriter, r *http.Request) {
	token := getTokenFromHeader(w, r)
	if token == -1 {
		return
	}
	id, ok := findKeyByValue(authorized.token, token)
	if !ok {
		http.Error(w, `{"error": "Invalid token"}`, http.StatusBadRequest)
		return
	}
	db := db.ConnectDB()
	products, order_price := db.PlacingOrder(id)
	defer db.CloseDB()
	if len(products) == 0 {
		http.Error(w, `{"error": "Cart is empty"}`, http.StatusUnprocessableEntity)
		return
	}
	storage := fmt.Sprintf("db/orderStorage/%d.json", id)
	if _, err := os.Stat(storage); err != nil {
		order := [1]Order{{IdO: 1, Products: products, Order_price: order_price}}
		file, err := os.Create(storage)
		if err != nil {
			panic(err)
		}
		defer file.Close()

		encoder := json.NewEncoder(file)
		if err = encoder.Encode(order); err != nil {
			panic(err)
		}
		db.DeleteCart(id)

		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"order_id": 1,	"message": "Order is processed"}`))
	} else {
		var orders []Order
		file, err := os.Open(storage)
		if err != nil {
			panic(err)
		}
		decoder := json.NewDecoder(file)
		if err = decoder.Decode(&orders); err != nil {
			panic(err)
		}
		file.Close()

		newOrder := Order{IdO: len(orders) + 1, Products: products, Order_price: order_price}
		orders = append(orders, newOrder)

		file, err = os.Create(storage)
		if err != nil {
			panic(err)
		}
		defer file.Close()

		encoder := json.NewEncoder(file)
		if err = encoder.Encode(orders); err != nil {
			panic(err)
		}
		db.DeleteCart(id)

		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(fmt.Sprintf(`{"order_id": %d, "message": "Order is processed"}`, newOrder.IdO)))
	}
}

func ViewOrder(w http.ResponseWriter, r *http.Request) {
	token := getTokenFromHeader(w, r)
	if token == -1 {
		return
	}
	id, ok := findKeyByValue(authorized.token, token)
	if !ok {
		http.Error(w, `{"error": "Invalid token"}`, http.StatusBadRequest)
		return
	}
	storage := fmt.Sprintf("db/orderStorage/%d.json", id)
	if _, err := os.Stat(storage); err != nil {
		http.Error(w, `{"error": "Orders not found"}`, http.StatusNotFound)
		return
	} else {
		file, err := os.Open(storage)
		if err != nil {
			panic(err)
		}
		defer file.Close()

		var orders []Order
		decoder := json.NewDecoder(file)
		if err = decoder.Decode(&orders); err != nil {
			panic(err)
		}

		w.Header().Set("Content-Type", "application/json")
		if err = json.NewEncoder(w).Encode(orders); err != nil {
			panic(err)
		}
	}
}

func TEST(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	fmt.Println(authorized.token)
	w.Write([]byte(fmt.Sprintf(`{"authorized.token": "%v"}`, authorized.token)))
}
