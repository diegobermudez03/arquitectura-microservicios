package main

import (
	"log"
	"os"

	"webhook/handlers"
	"webhook/internal"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"github.com/joho/godotenv"
)

func main() {
	// Load environment variables
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using system environment variables")
	}

	// Initialize Redis client
	redisClient := redis.NewClient(&redis.Options{
		Addr:     getEnv("REDIS_ADDR", "localhost:6379"),
		Password: getEnv("REDIS_PASSWORD", ""),
		DB:       0,
	})

	// Test Redis connection
	ctx := redisClient.Context()
	if err := redisClient.Ping(ctx).Err(); err != nil {
		log.Fatalf("Failed to connect to Redis: %v", err)
	}

	// Initialize Redis service
	redisService := internal.NewRedisService(redisClient)

	// Initialize handlers
	suscribeHandler := handlers.NewSuscribeHandler(redisService)
	forwardRequestHandler := handlers.NewForwardRequestHandler()
	forwardResponseHandler := handlers.NewForwardResponseHandler(redisService)

	// Setup Gin router
	router := gin.Default()

	// Register routes
	router.POST("/suscribe", suscribeHandler.HandleSuscribe)
	router.POST("/request", forwardRequestHandler.HandleForwardRequest)
	router.POST("/response", forwardResponseHandler.HandleForwardResponse)

	// Start server
	port := getEnv("PORT", "8080")
	log.Printf("Webhook Service starting on port %s", port)
	if err := router.Run(":" + port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
