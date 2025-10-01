package models

import (
	"time"

	"gorm.io/gorm"
)

// User represents a registered user
type User struct {
	Name        string `json:"name"`
	Lastname    string `json:"lastname"`
	BirthDate   string `json:"birth_date"`
	CountryCode string `json:"country_code"`
}

// RegisterRequest represents the request to register a user
type RegisterRequest struct {
	Name        string `json:"name" binding:"required"`
	Lastname    string `json:"lastname" binding:"required"`
	BirthDate   string `json:"birth_date" binding:"required"`
	CountryCode string `json:"country_code" binding:"required"`
}

// RegisterResponse represents the response after user registration
type RegisterResponse struct {
	Token string `json:"token"`
}

// UserRecord represents a user record in the database
type UserRecord struct {
	ID          string         `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	UserToken   string         `json:"user_token" gorm:"uniqueIndex;not null"`
	Name        string         `json:"name" gorm:"not null"`
	Lastname    string         `json:"lastname" gorm:"not null"`
	BirthDate   string         `json:"birth_date" gorm:"type:date;not null"`
	CountryCode string         `json:"country_code" gorm:"not null"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `json:"deleted_at" gorm:"index"`
}

// IssuedCardRecord represents an issued card record in the database
type IssuedCardRecord struct {
	ID         string         `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	UserID     string         `json:"user_id" gorm:"type:uuid;not null"`
	User       UserRecord     `json:"user" gorm:"foreignKey:UserID;references:ID"`
	UserToken  string         `json:"user_token" gorm:"not null"`
	PAN        string         `json:"pan" gorm:"not null"`
	CVV        string         `json:"cvv" gorm:"not null"`
	ExpiryDate string         `json:"expiry_date" gorm:"type:date;not null"`
	CardType   string         `json:"card_type" gorm:"not null"`
	Status     string         `json:"status" gorm:"not null"`
	CreatedAt  time.Time      `json:"created_at"`
	UpdatedAt  time.Time      `json:"updated_at"`
	DeletedAt  gorm.DeletedAt `json:"deleted_at" gorm:"index"`
}

// FailedAttemptRecord represents a failed attempt record in the database
type FailedAttemptRecord struct {
	ID            string         `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	UserID        string         `json:"user_id" gorm:"type:uuid;not null"`
	User          UserRecord     `json:"user" gorm:"foreignKey:UserID;references:ID"`
	UserToken     string         `json:"user_token" gorm:"not null"`
	CardType      string         `json:"card_type" gorm:"not null"`
	DeclineReason string         `json:"decline_reason" gorm:"not null"`
	Status        string         `json:"status" gorm:"not null"`
	CreatedAt     time.Time      `json:"created_at"`
	UpdatedAt     time.Time      `json:"updated_at"`
	DeletedAt     gorm.DeletedAt `json:"deleted_at" gorm:"index"`
}

// TableName methods for GORM
func (UserRecord) TableName() string {
	return "users"
}

func (IssuedCardRecord) TableName() string {
	return "issued_cards"
}

func (FailedAttemptRecord) TableName() string {
	return "failed_attempts"
}
