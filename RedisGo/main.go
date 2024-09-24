package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
)

type Product struct {
	ID    string  `json:"id"`
	Name  string  `json:"name"`
	Price float64 `json:"price"`
}

var ctx = context.Background()
var rdb *redis.Client

func main() {
	rdb = redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})

	r := gin.Default()

	// Define routes for CRUD operations
	r.POST("/product", createProduct)
	r.GET("/product/:id", getProduct)
	r.PUT("/product/:id", updateProduct)
	r.DELETE("/product/:id", deleteProduct)

	if err := r.Run(":8080"); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

func createProduct(c *gin.Context) {
	var product Product
	if err := c.ShouldBindJSON(&product); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	// Marshal product into JSON and store in Redis
	productJSON, err := json.Marshal(product)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error marshalling product"})
		return
	}

	err = rdb.Set(ctx, product.ID, productJSON, 0).Err()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error saving product to Redis"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Product added", "product": product})
}

// Get product by ID
func getProduct(c *gin.Context) {
	id := c.Param("id")

	productJSON, err := rdb.Get(ctx, id).Result()
	if err == redis.Nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Product not found"})
		return
	} else if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error retrieving product"})
		return
	}

	var product Product
	err = json.Unmarshal([]byte(productJSON), &product)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error unmarshalling product"})
		return
	}

	c.JSON(http.StatusOK, product)
}

// Update existing product
func updateProduct(c *gin.Context) {
	id := c.Param("id")

	var product Product
	if err := c.ShouldBindJSON(&product); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	product.ID = id

	productJSON, err := json.Marshal(product)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error marshalling product"})
		return
	}

	err = rdb.Set(ctx, product.ID, productJSON, 0).Err()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error updating product in Redis"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Product updated", "product": product})
}

// Delete product by ID
func deleteProduct(c *gin.Context) {
	id := c.Param("id")

	err := rdb.Del(ctx, id).Err()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error deleting product"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Product deleted"})
}
