package main

import (
	"fmt"
	"log"
	"os"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"github.com/joho/godotenv"

	"cards/handlers"
	"cards/internal"
)

func main() {
	// Load environment variables
	if err := godotenv.Load(".env"); err != nil {
		fmt.Println(err.Error())
		log.Println("No .env file found, using system environment variables")
	}

	// Initialize Redis client
	redisClient := redis.NewClient(&redis.Options{
		Addr:     os.Getenv("REDIS_ADDR"),
		Password: os.Getenv("REDIS_PASSWORD"),
		DB:       0,
	})

	// Test Redis connection
	if err := redisClient.Ping(redisClient.Context()).Err(); err != nil {
		log.Fatal("Failed to connect to Redis:", err)
	}
	log.Println("Redis client succesful ping")

	// Initialize PostgreSQL service with auto-migration
	postgresService := internal.NewPostgresService()
	log.Println("Database schema migrated successfully")

	// Initialize services
	redisService := internal.NewRedisService(redisClient)

	// Initialize handlers
	registerHandler := handlers.NewRegisterHandler(redisService, postgresService)
	issueHandler := handlers.NewIssueHandler(redisService, postgresService)
	webhookHandler := handlers.NewWebhookHandler(redisService, postgresService)
	cardsHandler := handlers.NewCardsHandler(postgresService)

	// Setup router
	router := gin.Default()

	// Configure CORS to allow all connections
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"*"},
		ExposeHeaders:    []string{"*"},
		AllowCredentials: false,
	}))

	v1 := router.Group("/v1")

	// Register routes
	v1.POST("/register", registerHandler.Register)
	v1.POST("/issue", issueHandler.Issue)
	v1.POST("/webhook", webhookHandler.Webhook)
	v1.GET("/:citizen_id/cards", cardsHandler.GetCardsByCitizenID)

	// Health check endpoint
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "healthy"})
	})

	// Start server
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Starting Service A on port %s", port)
	if err := router.Run(":" + port); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}
