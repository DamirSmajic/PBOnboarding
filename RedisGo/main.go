package main

import (
	"context"
	"fmt"
	"log"

	"github.com/redis/go-redis/v9"
)

type Product struct {
	ID    string  `json:"id"`
	Name  string  `json:"name"`
	Price float64 `json:"price"`
}

var ctx = context.Background()

func main() {
	client := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})

	ping, err := client.Ping(ctx).Result()
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	fmt.Println(ping)

	// Test dodavanja entiteta u Redis
	product := Product{ID: "1", Name: "Product 1", Price: 19.99}
	addProduct(client, &product)

	// Test preuzimanje entiteta iz Redis
	retrievedProduct := getProduct(client, "1")
	fmt.Printf("Retrieved Product: %+v\n", retrievedProduct)
}

// Dodavanje entiteta u Redis
func addProduct(rdb *redis.Client, product *Product) {
	err := rdb.Set(ctx, product.ID, product.Name, 0).Err()
	if err != nil {
		log.Fatalf("Error adding product: %v", err)
	}
	fmt.Println("Product added:", product.Name)
}

// Preuzimanje entiteta iz Redisa
func getProduct(rdb *redis.Client, id string) *Product {
	name, err := rdb.Get(ctx, id).Result()
	if err != nil {
		log.Fatalf("Error retrieving product: %v", err)
	}
	return &Product{ID: id, Name: name}
}
