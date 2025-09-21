package db

import (
	"database/sql"
	"fmt"
	"os"

	_ "github.com/mattn/go-sqlite3"
)

const storage = "db/main.db"

type Product struct {
	IdProduct   int
	Name        string
	Description string
	Price       int
}

func CreateDatabase() {
	if _, err := os.Stat(storage); err != nil {
		db, err := sql.Open("sqlite3", storage)
		if err != nil {
			panic(err)
		}
		defer db.Close()

		sqlBytes, err := os.ReadFile("db/createTables.sql")
		if err != nil {
			panic(err)
		}

		if _, err = db.Exec(string(sqlBytes)); err != nil {
			os.Remove(storage)
			panic(err)
		}
	} else {
		fmt.Println("Файл main.db существует")
	}
}

type cdb struct {
	db *sql.DB
}

func ConnectDB() *cdb {
	db, err := sql.Open("sqlite3", storage)
	if err != nil {
		panic(err)
	}
	return &cdb{db: db}
}

func (c *cdb) CloseDB() {
	c.db.Close()
}

func (c *cdb) Exec(query string, args ... any) (sql.Result, error) {
	return c.db.Exec(query, args...)
}

func (c *cdb) GetIdUser(email string) int {
	rows, err := c.db.Query("SELECT id FROM users WHERE email = ?", email)
	if err != nil {
		panic(err)
	}
	defer rows.Close()
	id := 1

	for rows.Next() {
		err := rows.Scan(&id)
		if err != nil {
			fmt.Println(err)
			break
		}
	}
	return id
}

func (c *cdb) GetPassword(email string) (int, string) {
	rows, err := c.db.Query("SELECT id, password FROM users WHERE email = ?", email)
	if err != nil {
		panic(err)
	}
	defer rows.Close()
	id := 1
	password := ""

	for rows.Next() {
		err := rows.Scan(&id, &password)
		if err != nil {
			fmt.Println(err)
			break
		}
	}
	return id, password
}

func (c *cdb) InsertUser(fio string, email string, password string, avatar string) {
	query := "INSERT INTO users (fio, email, password, avatar) VALUES (?, ?, ?, ?)"

	_, err := c.db.Exec(query, fio, email, password, avatar)
	if err != nil {
		panic(err)
	}
}

func (c *cdb) GetProduct() []Product {
	rows, err := c.db.Query("select * from products")
	if err != nil {
		panic(err)
	}
	defer rows.Close()
	products := []Product{}

	for rows.Next() {
		p := Product{}
		err := rows.Scan(&p.IdProduct, &p.Name, &p.Description, &p.Price)
		if err != nil {
			fmt.Println(err)
			continue
		}
		products = append(products, p)
	}
	return products
}

type UserData struct {
	Id     int    `json:"id"`
	Fio    string `json:"fio"`
	Email  string `json:"email"`
	Avatar string `json:"avatar"`
}

func (c *cdb) GetUser(id int) UserData {
	rows, err := c.db.Query("SELECT id, fio, email, avatar FROM users WHERE id = ?", id)
	if err != nil {
		panic(err)
	}
	defer rows.Close()
	user := UserData{}

	for rows.Next() {
		err := rows.Scan(&user.Id, &user.Fio, &user.Email, &user.Avatar)
		if err != nil {
			fmt.Println(err)
			break
		}
	}

	return user
}

func (c *cdb) GetProductById(id int) Product {
	rows, err := c.db.Query("SELECT * FROM products WHERE idP = ?", id)
	if err != nil {
		panic(err)
	}
	defer rows.Close()
	product := Product{}

	for rows.Next() {
		err := rows.Scan(&product.IdProduct, &product.Name, &product.Description, &product.Price)
		if err != nil {
			fmt.Println(err)
			continue
		}
	}
	return product
}

func (c *cdb) AddToCart(id int, product_id int) {
	_, err := c.db.Exec("INSERT INTO cart (id, idP) VALUES (?, ?)", id, product_id)
	if err != nil {
		panic(err)
	}
}

type Cart struct {
	IdC         int    `json:"id"`
	IdProduct   int    `json:"product_id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Price       int    `json:"price"`
}

func (c *cdb) ViewCart(idU int) []Cart {
	query := `SELECT cart.idC, products.idP, products.name, products.description, products.price
	FROM cart JOIN products ON cart.idP = products.idP WHERE cart.id = ?`
	rows, err := c.db.Query(query, idU)
	if err != nil {
		panic(err)
	}
	cart := []Cart{}

	for rows.Next() {
		c := Cart{}
		err := rows.Scan(&c.IdC, &c.IdProduct, &c.Name, &c.Description, &c.Price)
		if err != nil {
			fmt.Println(err)
			continue
		}
		cart = append(cart, c)
	}
	return cart
}

func (c *cdb) DeleteProduct(idC int, id int) {
	_, err := c.db.Exec("DELETE FROM cart WHERE idC = ? AND id = ?", idC, id)
	if err != nil {
		panic(err)
	}
}

func (c *cdb) DeleteCart(id int) {
	_, err := c.db.Exec("DELETE FROM cart WHERE id = ?", id)
	if err != nil {
		panic(err)
	}
}

func (c *cdb) PlacingOrder(id int) ([]int, int) {
	rows, err := c.db.Query("SELECT cart.idP, products.price FROM cart JOIN products ON cart.idP = products.idP WHERE id = ?", id)
	if err != nil {
		panic(err)
	}
	defer rows.Close()
	products := []int{}
	priceTotal := 0
	for rows.Next() {
		var prod, price int
		if err := rows.Scan(&prod, &price); err != nil {
			fmt.Println(err)
			continue
		}
		products = append(products, prod)
		priceTotal += price
	}
	return products, priceTotal
}
