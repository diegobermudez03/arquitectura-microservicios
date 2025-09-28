# Cards Microservice (Service A)

This is the main orchestrator microservice that handles user registration, credit card issuance requests, and webhook processing.

## Features

- **User Registration**: Register users and store session data in Redis and PostgreSQL
- **Card Issuance**: Initiate credit card issuance requests
- **Webhook Processing**: Handle issuance results from webhook service
- **Notifications**: Send notifications to users via notifications service
- **Auto-Migration**: Automatic database schema creation and updates using GORM

## Database Schema

The microservice uses **GORM** for automatic database migrations. The schema is automatically created and updated when the application starts.

### Tables (Auto-Created by GORM)

#### 1. Users Table
Stores user information with unique tokens:
- `id` (UUID, Primary Key, Auto-generated)
- `user_token` (TEXT, Unique)
- `name`, `lastname`, `birth_date`, `country_code`
- `created_at`, `updated_at`, `deleted_at` (Soft Delete Support)

#### 2. Issued Cards Table
Stores successfully issued cards:
- `id` (UUID, Primary Key, Auto-generated)
- `user_id` (UUID, Foreign Key to users)
- `user_token` (TEXT)
- `pan`, `cvv`, `expiry_date`, `card_type`
- `status`, `created_at`, `updated_at`, `deleted_at`

#### 3. Failed Attempts Table
Stores declined card issuance attempts:
- `id` (UUID, Primary Key, Auto-generated)
- `user_id` (UUID, Foreign Key to users)
- `user_token` (TEXT)
- `card_type`, `decline_reason`, `status`
- `created_at`, `updated_at`, `deleted_at`

## API Endpoints

### POST /register
Registers a new user and returns a user token.

**Request:**
```json
{
  "name": "John",
  "lastname": "Doe", 
  "birth_date": "2002-03-25",
  "country_code": "US"
}
```

**Response:**
```json
{
  "token": "generated_user_token"
}
```

### POST /issue
Initiates a credit card issuance request.

**Request:**
```json
{
  "card_type": "credit",
  "user_token": "user_token_value"
}
```

**Response:** 202 Accepted

### POST /webhook
Webhook endpoint for receiving issuance results.

**Request:** WebhookEvent with IssuerResponse data

**Response:** 200 OK

### GET /health
Health check endpoint.

**Response:**
```json
{
  "status": "healthy"
}
```

## Environment Variables

Create a `.env` file based on `.env.example`:

- `REDIS_ADDR`: Redis server address
- `REDIS_PASSWORD`: Redis password (optional)
- `POSTGRES_URL`: PostgreSQL connection string
- `WEBHOOK_URL`: Webhook service URL
- `NOTIFICATIONS_URL`: Notifications service URL
- `SUSCRIPTOR_TOKEN`: Suscriptor token for webhook service
- `PORT`: Server port (default: 8080)

## Setup

1. **Install dependencies:**
```bash
go mod tidy
```

2. **Create PostgreSQL database:**
```bash
createdb cards_db
# or
psql -c "CREATE DATABASE cards_db;"
```

3. **Configure environment variables in `.env` file**

4. **Run the service:**
```bash
go run main.go
```

The application will automatically:
- Connect to PostgreSQL
- Create all required tables and indexes
- Set up foreign key relationships
- Apply any schema migrations

## Dependencies

- **Gin** (HTTP framework)
- **GORM** (ORM with auto-migration)
- **Redis** (session storage)
- **PostgreSQL** (persistent storage)
- **UUID** (request identification)
- **Godotenv** (environment variables)

## GORM Features

- **Auto-Migration**: Tables are created automatically
- **Soft Deletes**: Records are marked as deleted instead of removed
- **Foreign Keys**: Automatic relationship management
- **Indexes**: Automatic index creation for performance
- **Timestamps**: Automatic created_at/updated_at tracking

## Project Structure

```
cards/
 main.go              # Application entry point with auto-migration
 handlers/            # HTTP handlers
    register.go     # User registration
    issue.go        # Card issuance
    webhook.go      # Webhook processing
 models/             # GORM models with auto-migration
    user.go         # User and database models
    request.go      # Request models
    card.go         # Card and webhook models
 internal/           # Internal services
    redis.go        # Redis service
    postgres.go     # GORM PostgreSQL service
 schema.sql          # Database schema reference
 .env.example        # Environment variables example
 go.mod             # Go module file
```

## Migration Notes

- **No manual migration needed**: GORM handles everything automatically
- **Schema updates**: Add new fields to models and restart the application
- **Data preservation**: Existing data is preserved during schema updates
- **Rollback**: Use GORM migrations for complex rollback scenarios
