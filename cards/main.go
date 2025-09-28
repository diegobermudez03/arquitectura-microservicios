package main

import (
"log"
"os"

"github.com/gin-gonic/gin"
"github.com/go-redis/redis/v8"
"github.com/joho/godotenv"

"cards/handlers"
"cards/internal"
)

func main() {
// Load environment variables
if err := godotenv.Load(); err != nil {
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

// Initialize PostgreSQL service with auto-migration
postgresService := internal.NewPostgresService()
log.Println("Database schema migrated successfully")

// Initialize services
redisService := internal.NewRedisService(redisClient)

// Initialize handlers
registerHandler := handlers.NewRegisterHandler(redisService, postgresService)
issueHandler := handlers.NewIssueHandler(redisService, postgresService)
webhookHandler := handlers.NewWebhookHandler(redisService, postgresService)

// Setup router
router := gin.Default()

// Register routes
router.POST("/register", registerHandler.Register)
router.POST("/issue", issueHandler.Issue)
router.POST("/webhook", webhookHandler.Webhook)

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
