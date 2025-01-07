# Chirpy

Chirpy is a social media backend API built with Go, featuring user authentication, post management, and premium user features.

## Features

- User authentication with JWT tokens and refresh tokens
- CRUD operations for "chirps" (posts)
- Profanity filtering
- Premium user upgrades (Chirpy Red)
- PostgreSQL database integration
- API key authentication for webhooks

## Tech Stack

- Go 1.23.2
- PostgreSQL
- JWT for authentication
- SQLC for type-safe SQL queries
- Goose for database migrations

## Project Structure

.
├── internal/
│ ├── auth/ # Authentication utilities
│ └── database/ # Database models and queries
├── sql/
│ ├── queries/ # SQLC query definitions
│ └── schema/ # Database migrations
├── main.go # Application entry point
├── chirps.go # Chirp-related handlers
├── users.go # User-related handlers
├── webhooks.go # Webhook handlers
└── README.md

## API Endpoints

### Authentication

- `POST /api/users` - Create a new user
- `POST /api/login` - Login user
- `POST /api/refresh` - Refresh access token
- `POST /api/revoke` - Revoke refresh token

### Chirps

- `GET /api/chirps` - Get all chirps
- `GET /api/chirps/{chirpID}` - Get specific chirp
- `POST /api/chirps` - Create new chirp
- `DELETE /api/chirps/{chirpID}` - Delete chirp

### User Management

- `PUT /api/users` - Update user information
- `POST /api/polka/webhooks` - Handle user upgrades to Chirpy Red

### System

- `GET /api/healthz` - Health check endpoint
- `GET /admin/metrics` - View system metrics
- `POST /admin/reset` - Reset application (development only)

## Setup

1. Clone the repository
2. Create a `.env` file with the following variables:
   ```
   DB_URL=your_postgresql_connection_string
   JWT_SECRET=your_jwt_secret
   POLKA_KEY=your_webhook_api_key
   PLATFORM=dev|prod
   ```
3. Run database migrations:
   ```bash
   goose -dir sql/schema postgres "your_connection_string" up
   ```
4. Start the server:
   ```bash
   go run .
   ```

## Database Schema

The application uses three main tables:

- `users` - Stores user information and authentication details
- `chirps` - Stores user posts with foreign key relationships
- `refresh_tokens` - Manages JWT refresh tokens

## Security Features

- Password hashing using bcrypt
- JWT-based authentication
- API key validation for webhooks
- Refresh token rotation
- SQL injection prevention through prepared statements

## Development Notes

- The application includes profanity filtering for chirp content
- Chirps are limited to 140 characters
- Premium features are managed through the Chirpy Red flag
- Development-only endpoints are protected by environment checks

## Contributing

1. Fork the repository
2. Create a feature branch
3. Commit your changes
4. Push to the branch
5. Create a Pull Request
