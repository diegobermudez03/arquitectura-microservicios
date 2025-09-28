package internal

import (
	"fmt"
	"os"

	"cards/models"

	"github.com/google/uuid"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type PostgresService struct {
	db *gorm.DB
}

func NewPostgresService() *PostgresService {
	postgresURL := os.Getenv("POSTGRES_URL")
	if postgresURL == "" {
		panic("POSTGRES_URL environment variable is required")
	}

	db, err := gorm.Open(postgres.Open(postgresURL), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		panic("Failed to connect to PostgreSQL: " + err.Error())
	}

	// Auto-migrate the schema
	err = db.AutoMigrate(
		&models.UserRecord{},
		&models.IssuedCardRecord{},
		&models.FailedAttemptRecord{},
	)
	if err != nil {
		panic("Failed to migrate database: " + err.Error())
	}

	return &PostgresService{
		db: db,
	}
}

// StoreUser stores a user in the database
func (p *PostgresService) StoreUser(userToken string, user models.User) (*models.UserRecord, error) {
	userRecord := &models.UserRecord{
		ID:          uuid.New().String(),
		UserToken:   userToken,
		Name:        user.Name,
		Lastname:    user.Lastname,
		BirthDate:   user.BirthDate,
		CountryCode: user.CountryCode,
	}

	result := p.db.Create(userRecord)
	if result.Error != nil {
		return nil, result.Error
	}

	return userRecord, nil
}

// GetUserByToken retrieves a user by token
func (p *PostgresService) GetUserByToken(userToken string) (*models.UserRecord, error) {
	var user models.UserRecord
	result := p.db.Where("user_token = ?", userToken).First(&user)
	if result.Error != nil {
		return nil, result.Error
	}

	return &user, nil
}

// StoreIssuedCard stores an issued card in the database
func (p *PostgresService) StoreIssuedCard(record models.IssuedCardRecord) error {
	result := p.db.Create(&record)
	return result.Error
}

// StoreFailedAttempt stores a failed attempt in the database
func (p *PostgresService) StoreFailedAttempt(record models.FailedAttemptRecord) error {
	result := p.db.Create(&record)
	return result.Error
}

// CreateIssuedCardRecord creates a record for successful card issuance
func (p *PostgresService) CreateIssuedCardRecord(
	userID string,
	userToken string,
	response models.IssuerResponse,
) models.IssuedCardRecord {
	return models.IssuedCardRecord{
		ID:         uuid.New().String(),
		UserID:     userID,
		UserToken:  userToken,
		PAN:        response.IssuedCard.PAN,
		CVV:        response.IssuedCard.CVV,
		ExpiryDate: response.IssuedCard.ExpiryDate,
		CardType:   response.IssuedCard.CardType,
		Status:     response.Status,
	}
}

// CreateFailedAttemptRecord creates a record for failed card issuance
func (p *PostgresService) CreateFailedAttemptRecord(
	userID string,
	userToken string,
	cardType string,
	response models.IssuerResponse,
) models.FailedAttemptRecord {
	return models.FailedAttemptRecord{
		ID:            uuid.New().String(),
		UserID:        userID,
		UserToken:     userToken,
		CardType:      cardType,
		DeclineReason: response.DeclineReason.Reason,
		Status:        response.Status,
	}
}

// SendNotification sends a notification to the notifications service
func (p *PostgresService) SendNotification(userToken string, response models.IssuerResponse) error {
	fmt.Printf("Sending notification for user %s with status %s\n", userToken, response.Status)
	return nil
}

// GetDB returns the GORM database instance for advanced queries if needed
func (p *PostgresService) GetDB() *gorm.DB {
	return p.db
}
