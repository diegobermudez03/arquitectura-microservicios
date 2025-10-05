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
		CitizenID:   user.CitizenID, // Social Security ID
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

// FullCardResult represents the joined result from users and cards tables
type FullCardResult struct {
	// User fields
	UserID          string `gorm:"column:user_id"`
	UserToken       string `gorm:"column:user_token"`
	UserName        string `gorm:"column:name"`
	UserLastname    string `gorm:"column:lastname"`
	UserBirthDate   string `gorm:"column:birth_date"`
	UserCountryCode string `gorm:"column:country_code"`
	UserSocialID    string `gorm:"column:citizen_id"`
	UserCreatedAt   string `gorm:"column:user_created_at"`

	// Card fields
	CardID        string `gorm:"column:card_id"`
	CardPAN       string `gorm:"column:pan"`
	CardCVV       string `gorm:"column:cvv"`
	CardExpiry    string `gorm:"column:expiry_date"`
	CardType      string `gorm:"column:card_type"`
	CardStatus    string `gorm:"column:card_status"`
	CardCreatedAt string `gorm:"column:card_created_at"`
}

// GetCardsByCitizenID retrieves all users and their cards for a specific citizen ID
func (p *PostgresService) GetCardsByCitizenID(citizenID string) ([]models.FullCard, error) {
	var results []FullCardResult

	// Single query with JOIN
	result := p.db.Table("users").
		Select(`
			users.id as user_id,
			users.user_token,
			users.name,
			users.lastname,
			users.birth_date,
			users.country_code,
			users.citizen_id,
			users.created_at as user_created_at,
			issued_cards.id as card_id,
			issued_cards.pan,
			issued_cards.cvv,
			issued_cards.expiry_date,
			issued_cards.card_type,
			issued_cards.status as card_status,
			issued_cards.created_at as card_created_at
		`).
		Joins("LEFT JOIN issued_cards ON users.id = issued_cards.user_id").
		Where("users.citizen_id = ?", citizenID).
		Find(&results)

	if result.Error != nil {
		return nil, result.Error
	}

	// Convert to FullCard models
	fullCards := make([]models.FullCard, len(results))
	for i, result := range results {
		fullCards[i] = models.FullCard{
			UserID:          result.UserID,
			UserToken:       result.UserToken,
			UserName:        result.UserName,
			UserLastname:    result.UserLastname,
			UserBirthDate:   result.UserBirthDate,
			UserCountryCode: result.UserCountryCode,
			UserSocialID:    result.UserSocialID,
			UserCreatedAt:   result.UserCreatedAt,
			CardID:          result.CardID,
			CardPAN:         result.CardPAN,
			CardCVV:         result.CardCVV,
			CardExpiry:      result.CardExpiry,
			CardType:        result.CardType,
			CardStatus:      result.CardStatus,
			CardCreatedAt:   result.CardCreatedAt,
		}
	}

	return fullCards, nil
}

// GetDB returns the GORM database instance for advanced queries if needed
func (p *PostgresService) GetDB() *gorm.DB {
	return p.db
}
