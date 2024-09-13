package main

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/go-sql-driver/mysql"
)

// Definiši strukture podataka
type Customer struct {
	ID   int
	Name string
}

type Product struct {
	ID    int
	Name  string
	Price float64
}

type Order struct {
	ID         int
	CustomerID int
	ProductID  int
}

type OrderDetails struct {
	ID        int
	OrderID   int
	Quantity  int
	TotalCost float64
}

func main() {
	// Definisanje DSN-a za MySQL
	dsn := "root:password@tcp(127.0.0.1:3306)/pbdatabase"

	// Otvaranje konekcije s MySQL bazom
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		log.Fatalf("Error connecting to the database: %v", err)
	}
	defer db.Close()

	// Testiranje konekcije
	err = db.Ping()
	if err != nil {
		log.Fatalf("Error pinging the database: %v", err)
	}
	fmt.Println("Successfully connected to MySQL database!")

	createTables(db)

	addCustomer(db, "John Doe")
}

func createTables(db *sql.DB) {
	queries := []string{
		`CREATE TABLE IF NOT EXISTS customers (
            id INT AUTO_INCREMENT PRIMARY KEY,
            name VARCHAR(100) NOT NULL
        );`,
		`CREATE TABLE IF NOT EXISTS products (
            id INT AUTO_INCREMENT PRIMARY KEY,
            name VARCHAR(100) NOT NULL,
            price DECIMAL(10, 2) NOT NULL
        );`,
		`CREATE TABLE IF NOT EXISTS orders (
            id INT AUTO_INCREMENT PRIMARY KEY,
            customer_id INT,
            product_id INT,
            FOREIGN KEY (customer_id) REFERENCES customers(id),
            FOREIGN KEY (product_id) REFERENCES products(id)
        );`,
		`CREATE TABLE IF NOT EXISTS order_details (
            id INT AUTO_INCREMENT PRIMARY KEY,
            order_id INT,
            quantity INT,
            total_cost DECIMAL(10, 2),
            FOREIGN KEY (order_id) REFERENCES orders(id)
        );`,
	}

	for _, query := range queries {
		_, err := db.Exec(query)
		if err != nil {
			log.Fatalf("Error creating table: %v", err)
		}
	}

	fmt.Println("Tables created successfully!")
}

func addCustomer(db *sql.DB, name string) {
	query := `INSERT INTO customers (name) VALUES (?)`
	result, err := db.Exec(query, name)
	if err != nil {
		log.Fatalf("Error inserting customer: %v", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		log.Fatalf("Error getting last insert id: %v", err)
	}

	fmt.Printf("Customer added with ID %d\n", id)
}
