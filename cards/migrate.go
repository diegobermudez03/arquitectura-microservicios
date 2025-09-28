package main

import (
"log"
"os"

"cards/models"
"gorm.io/driver/postgres"
"gorm.io/gorm"
"gorm.io/gorm/logger"
)

// Example of how to run manual migrations if needed
// This is optional - GORM auto-migration handles most cases
func runMigrations() {
postgresURL := os.Getenv("POSTGRES_URL")
if postgresURL == "" {
log.Fatal("POSTGRES_URL environment variable is required")
}

db, err := gorm.Open(postgres.Open(postgresURL), &gorm.Config{
Logger: logger.Default.LogMode(logger.Info),
})
if err != nil {
log.Fatal("Failed to connect to PostgreSQL:", err)
}

// Auto-migrate all models
err = db.AutoMigrate(
&models.UserRecord{},
&models.IssuedCardRecord{},
&models.FailedAttemptRecord{},
)
if err != nil {
log.Fatal("Failed to migrate database:", err)
}

log.Println("Database migration completed successfully")
}

// Example of how to add custom indexes or constraints
func addCustomIndexes(db *gorm.DB) {
// Add custom indexes if needed
db.Exec("CREATE INDEX IF NOT EXISTS idx_issued_cards_status ON issued_cards(status)")
db.Exec("CREATE INDEX IF NOT EXISTS idx_failed_attempts_status ON failed_attempts(status)")

// Add custom constraints if needed
// db.Exec("ALTER TABLE issued_cards ADD CONSTRAINT chk_pan_length CHECK (LENGTH(pan) = 16)")
}

// Example of how to seed initial data
func seedInitialData(db *gorm.DB) {
// Example: Create admin user or initial data
// This is just an example - remove if not needed
log.Println("Seeding initial data...")
// Add your seed data here
}

// Uncomment the main function below to run migrations manually
// func main() {
// runMigrations()
// 
// // Get database connection
// postgresURL := os.Getenv("POSTGRES_URL")
// db, _ := gorm.Open(postgres.Open(postgresURL), &gorm.Config{})
// 
// addCustomIndexes(db)
// seedInitialData(db)
// }
