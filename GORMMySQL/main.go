package main

import (
	"fmt"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// Definiši modele
type Customer struct {
	gorm.Model
	Name string
}

type Product struct {
	gorm.Model
	Name  string
	Price float64
}

type Order struct {
	gorm.Model
	CustomerID uint
	ProductID  uint
	Customer   Customer
	Product    Product
}

type OrderDetails struct {
	gorm.Model
	OrderID   uint
	Quantity  int
	TotalCost float64
	Order     Order
}

func main() {
	// Definisanje DSN za GORM
	dsn := "root:password@tcp(127.0.0.1:3306)/pbdatabase?charset=utf8mb4&parseTime=True&loc=Local"

	// Otvaranje konekcije s MySQL bazom koristeći GORM
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	fmt.Println("Successfully connected to MySQL database using GORM!")

	// Automatski migrira tj. mapira na tabelu u bazi podataka
	db.AutoMigrate(&Customer{}, &Product{}, &Order{}, &OrderDetails{})

	addCustomer(db, "Gormy Gorm")
}

func addCustomer(db *gorm.DB, name string) {
	customer := Customer{Name: name}
	result := db.Create(&customer)
	if result.Error != nil {
		panic(result.Error)
	}

	fmt.Printf("Customer added with ID %d\n", customer.ID)
}
