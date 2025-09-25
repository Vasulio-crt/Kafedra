package users

import (
	"backendAPI/db"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/gorilla/mux"
)

var authorized = NewAuthorizedUsers()
var DB = db.ConnectDB()

// ---Работа с аккаунтами---

func SignUp(w http.ResponseWriter, r *http.Request) {
	var req RegistrationData
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		panic(err)
	}

	if req.Fio == "" || !isValidEmail(req.Email) || len(req.Password) < 3 {
		errorJSON(w, "Validation error", 422)
	}
	if req.Avatar == "" {
		req.Avatar = "avatars/default.jpeg"
	}

	DB.InsertUser(req.Fio, req.Email, req.Password, req.Avatar)
	idU := DB.GetIdUser(req.Email)

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
		panic(err)
	}

	if !isValidEmail(req.Email) || len(req.Password) < 3 {
		errorJSON(w, "Validation error", 422)
	}

	idU, password := DB.GetPassword(req.Email)

	if password != req.Password {
		errorJSON(w, "Login failed", 403)
		return
	}

	token := int(time.Now().UnixMicro())
	authorized.AddToken(idU, token)

	w.Header().Set("Content-Type", "application/json")

	mes := fmt.Sprintf(`{"user_token": %d}`, token)
	w.Write([]byte(mes))
}

func Logout(w http.ResponseWriter, r *http.Request) {
	id, ok := whatIdUser(r)

	if ok {
		authorized.RemoveToken(id)
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"message": "logout"}`))
	} else {
		errorJSON(w, "Login failed", 403)
		return
	}
}

func EditProfile(w http.ResponseWriter, r *http.Request) {
	id, ok := whatIdUser(r)
	if !ok {
		errorJSON(w, "Login failed", 403)
	}

	var req RegistrationData
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		panic(err)
	}

	if req.Avatar != "" {
		if _, err := DB.Exec("UPDATE users SET avatar = ? WHERE id = ?", req.Avatar, id); err != nil {
			panic(err)
		}
	}
	if req.Fio != "" {
		if _, err := DB.Exec("UPDATE users SET fio = ? WHERE id = ?", req.Fio, id); err != nil {
			panic(err)
		}
	}
	if req.Password != "" {
		if _, err := DB.Exec("UPDATE users SET password = ? WHERE id = ?", req.Password, id); err != nil {
			panic(err)
		}
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(`{"message": "data updated successfully"}`))
}

// ---Просмотр---

func ViewProduct(w http.ResponseWriter, r *http.Request) {
	products := DB.GetProduct()

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(products); err != nil {
		panic(err)
	}
}

func ViewProfile(w http.ResponseWriter, r *http.Request) {
	id, ok := whatIdUser(r)
	if !ok {
		errorJSON(w, "Login failed", 403)
		return
	}

	user := DB.GetUser(id)

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(ProfileData{User: user}); err != nil {
		panic(err)
	}
}

// ---Корзина---

func AddingProductCart(w http.ResponseWriter, r *http.Request) {
	product_id, err := strconv.Atoi(mux.Vars(r)["product_id"])
	if err != nil {
		errorJSON(w, "Invalid product ID", 400)
		return
	}

	id, ok := whatIdUser(r)
	if !ok {
		errorJSON(w, "Login failed", 403)
		return
	}

	product := DB.GetProductById(product_id)

	if product.IdProduct == 0 {
		errorJSON(w, "Invalid product ID", 400)
		return
	}

	DB.AddToCart(id, product_id)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte(`{"message": "Product add to card"}`))
}

func ViewCart(w http.ResponseWriter, r *http.Request) {
	id, ok := whatIdUser(r)
	if !ok {
		errorJSON(w, "Login failed", 403)
		return
	}

	cart := DB.ViewCart(id)

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(cart); err != nil {
		panic(err)
	}
}

func DeleteProduct(w http.ResponseWriter, r *http.Request) {
	idC, err := strconv.Atoi(mux.Vars(r)["idC"])
	if err != nil {
		errorJSON(w, "Forbidden for you", 403)
		return
	}

	id, ok := whatIdUser(r)
	if !ok {
		errorJSON(w, "Login failed", 403)
		return
	}

	DB.DeleteProduct(idC, id)

	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(`{"message": "Item removed from cart"}`))
}

// ---Заказ---

func PlacingOrder(w http.ResponseWriter, r *http.Request) {
	id, ok := whatIdUser(r)
	if !ok {
		errorJSON(w, "Login failed", 403)
		return
	}

	products, order_price := DB.PlacingOrder(id)

	if len(products) == 0 {
		errorJSON(w, "Cart is empty", 400)
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
		DB.DeleteCart(id)

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
		DB.DeleteCart(id)

		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(fmt.Sprintf(`{"order_id": %d, "message": "Order is processed"}`, newOrder.IdO)))
	}
}

func ViewOrder(w http.ResponseWriter, r *http.Request) {
	id, ok := whatIdUser(r)
	if !ok {
		errorJSON(w, "Login failed", 403)
		return
	}

	storage := fmt.Sprintf("db/orderStorage/%d.json", id)
	if _, err := os.Stat(storage); err != nil {
		errorJSON(w, "Orders not found", 404)
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
