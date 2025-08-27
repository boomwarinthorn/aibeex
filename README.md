# Fiber Authentication API - Clean Architecture

A Go web server using the Fiber framework with Clean Architecture principles, providing user authentication and registration features with JWT token support and SQLite database persistence.

## 🏗️ Architecture Overview

This project follows **Clean Architecture** principles with clear separation of concerns:

- **Domain Layer**: Core business entities and repository interfaces
- **Use Case Layer**: Business logic and application services  
- **Infrastructure Layer**: Database implementations and external services
- **Presentation Layer**: HTTP handlers, DTOs, and middleware
- **Package Layer**: Reusable utilities (JWT, validation)

## 📁 Project Structure

```
fiber-hello-world/
├── cmd/
│   └── api/
│       └── main.go                 # Application entry point
├── config/
│   └── config.go                   # Configuration management
├── internal/
│   ├── domain/
│   │   ├── entity/
│   │   │   └── user.go             # User domain entity
│   │   └── repository/
│   │       └── user_repository.go  # Repository interface
│   ├── usecase/
│   │   └── user_usecase.go         # Business logic layer
│   ├── infrastructure/
│   │   └── database/
│   │       └── sqlite_user_repository.go  # SQLite implementation
│   └── presentation/
│       ├── dto/
│       │   └── user_dto.go         # Data Transfer Objects
│       ├── handler/
│       │   └── user_handler.go     # HTTP request handlers
│       └── middleware/
│           └── jwt_middleware.go   # JWT authentication middleware
├── pkg/
│   ├── jwt/
│   │   └── jwt.go                  # JWT service utilities
│   └── validator/
│       └── validator.go            # Validation service
├── docs/                           # Swagger documentation files
├── go.mod                          # Go modules
├── go.sum                          # Go modules checksums
├── Makefile                        # Build and development commands
├── README.md                       # This file
└── users.db                        # SQLite database (auto-generated)
```

## 🚀 Features

- ✅ **Clean Architecture** with proper dependency injection
- ✅ Hello World JSON API endpoint
- ✅ User registration with comprehensive validation
- ✅ User login with JWT token generation (24h expiry)
- ✅ JWT-protected user profile endpoint
- ✅ Interactive Swagger/OpenAPI documentation
- ✅ SQLite database for persistent data storage
- ✅ Password hashing with bcrypt
- ✅ Email validation and duplicate prevention
- ✅ Structured error handling and responses
- ✅ Middleware-based JWT authentication
- ✅ Environment-based configuration
- ✅ Database auto-initialization

## 🛠️ Installation & Setup

### Prerequisites
- Go 1.19 or higher
- Git

### Steps
1. Clone this repository:
   ```bash
   git clone <repository-url>
   cd fiber-hello-world
   ```

2. Install dependencies:
   ```bash
   make deps
   # or
   go mod download && go mod tidy
   ```

## 🚀 Running the Application

### Using Makefile (Recommended)
```bash
# Run the application
make run

# Run in development mode with hot reload (requires air)
make dev

# Build the application
make build

# Run tests
make test

# View all available commands
make help
```

### Using Go directly
```bash
# Run from project root
go run cmd/api/main.go
```

The server will start on `http://localhost:3000`

### Environment Variables
You can configure the application using environment variables:

```bash
export PORT=8080
export JWT_SECRET=your-super-secret-key
export DB_PATH=./data/users.db
```

## 📚 API Documentation

### Swagger UI
Interactive API documentation is available at: `http://localhost:3000/swagger/index.html`

The Swagger UI provides:
- Complete API endpoint documentation  
- Interactive testing interface
- Request/response schema definitions
- Authentication examples with Bearer tokens

### Regenerating Swagger Documentation
When you make changes to the API endpoints:

```bash
# Using Makefile
make swagger

# Using swag directly  
go run github.com/swaggo/swag/cmd/swag@latest init -g cmd/api/main.go -o docs
```

## 🗄️ Database Architecture

The application uses SQLite for persistent data storage with repository pattern implementation.

### Database Schema

**Users Table:**
```sql
CREATE TABLE users (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    email TEXT UNIQUE NOT NULL,
    password TEXT NOT NULL,
    full_name TEXT NOT NULL,
    phone_number TEXT NOT NULL,
    birthday TEXT NOT NULL,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP
);
```

### Database Features:
- Automatic table creation on first run
- Email uniqueness constraint  
- Auto-incrementing user IDs
- Timestamps for user creation
- Persistent storage across server restarts

## 🏛️ Clean Architecture Layers

### 1. Domain Layer (`internal/domain/`)
The core of the application containing business entities and repository interfaces.

**Entities** (`entity/`):
- `User`: Core business entity representing a user with validation methods
- Pure Go structs with no external dependencies
- Contains business logic methods like `WithoutPassword()`

**Repository Interfaces** (`repository/`):
- `UserRepository`: Defines data access contract
- Database-agnostic interface
- Allows easy switching between different storage implementations

### 2. Use Case Layer (`internal/usecase/`)
Contains application-specific business logic and orchestrates the flow of data.

**Features:**
- `UserUseCase`: Handles user registration, authentication, and retrieval
- Implements business rules (password hashing, email validation)
- Coordinates between domain entities and repositories
- Returns domain entities or business errors

### 3. Infrastructure Layer (`internal/infrastructure/`)
Implements external concerns like databases and external services.

**Database** (`database/`):
- `SQLiteUserRepository`: Concrete implementation of UserRepository
- Handles SQLite-specific operations
- Database connection management
- SQL query implementations

### 4. Presentation Layer (`internal/presentation/`)
Handles HTTP concerns and user interface.

**Handlers** (`handler/`):
- `UserHandler`: HTTP request/response handling
- Converts between DTOs and domain entities
- HTTP status code management
- Error response formatting

**DTOs** (`dto/`):
- Request/Response data structures
- Input validation tags
- JSON serialization structures
- Separate from domain entities

**Middleware** (`middleware/`):
- `JWTMiddleware`: Token validation and user context
- Cross-cutting concerns
- Request/response processing

### 5. Package Layer (`pkg/`)
Reusable utilities and services used across the application.

**JWT Service** (`jwt/`):
- Token generation and validation
- Claims management
- Security configurations

**Validator Service** (`validator/`):
- Input validation wrapper
- Struct validation using tags
- Error formatting

### 6. Configuration (`config/`)
Application configuration management with environment variable support.

## 🔄 Dependency Flow

```
Presentation Layer → Use Case Layer → Domain Layer ← Infrastructure Layer
```

- **Presentation** depends on **Use Case**
- **Use Case** depends on **Domain** (entities & interfaces)
- **Infrastructure** implements **Domain** interfaces
- **Domain** has no dependencies (dependency inversion)

This structure ensures:
- ✅ **Testability**: Each layer can be tested in isolation
- ✅ **Maintainability**: Clear separation of concerns
- ✅ **Flexibility**: Easy to swap implementations
- ✅ **Scalability**: New features follow established patterns

## API Endpoints

### GET `/`
Returns a JSON response with "Hello World" message.

**Response:**
```json
{
  "message": "Hello World"
}
```

**Example:**
```bash
curl http://localhost:3000/
```

### POST `/register`
Register a new user account.

**Request Body:**
```json
{
  "email": "user@example.com",
  "password": "password123",
  "fullName": "John Doe",
  "phoneNumber": "0812345678",
  "birthday": "1990-01-15"
}
```

**Validation Rules:**
- `email`: Must be a valid email format
- `password`: Minimum 6 characters
- `fullName`: Minimum 2 characters
- `phoneNumber`: Minimum 10 characters
- `birthday`: Must be in YYYY-MM-DD format

**Success Response (201):**
```json
{
  "message": "User registered successfully",
  "user": {
    "id": 1,
    "email": "user@example.com",
    "fullName": "John Doe",
    "phoneNumber": "0812345678",
    "birthday": "1990-01-15",
    "createdAt": "2025-08-27T14:00:00Z"
  }
}
```

**Error Responses:**

*400 - Validation Failed:*
```json
{
  "error": "Validation failed",
  "message": "Field validation error details..."
}
```

*409 - Email Already Exists:*
```json
{
  "error": "Email already exists",
  "message": "User with this email already registered"
}
```

**Example:**
```bash
curl -X POST http://localhost:3000/register \
-H "Content-Type: application/json" \
-d '{
  "email": "test@example.com",
  "password": "password123",
  "fullName": "John Doe",
  "phoneNumber": "0812345678",
  "birthday": "1990-01-15"
}'
```

### GET `/me`
Get current user information using JWT token.

**Headers:**
```
Authorization: Bearer <jwt_token>
```

**Success Response (200):**
```json
{
  "message": "User information retrieved successfully",
  "user": {
    "id": 1,
    "email": "user@example.com",
    "fullName": "John Doe",
    "phoneNumber": "0812345678",
    "birthday": "1990-01-15",
    "createdAt": "2025-08-27T14:00:00Z"
  }
}
```

**Error Responses:**

*401 - Authorization Required:*
```json
"Authorization header required"
```

*401 - Bearer Token Required:*
```json
"Bearer token required"
```

*401 - Invalid Token:*
```json
"Invalid token"
```

*404 - User Not Found:*
```json
{
  "error": "User not found",
  "message": "User associated with this token no longer exists"
}
```

**Example:**
```bash
# First login to get token
TOKEN=$(curl -s -X POST http://localhost:3000/login \
-H "Content-Type: application/json" \
-d '{"email":"test@example.com","password":"password123"}' | jq -r .token)

# Then use token to get user info
curl -X GET http://localhost:3000/me \
-H "Authorization: Bearer $TOKEN"

### POST `/login`
Authenticate user and receive JWT token.

**Request Body:**
```json
{
  "email": "user@example.com",
  "password": "password123"
}
```

**Validation Rules:**
- `email`: Must be a valid email format
- `password`: Required field

**Success Response (200):**
```json
{
  "message": "Login successful",
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "user": {
    "id": 1,
    "email": "user@example.com",
    "fullName": "John Doe",
    "phoneNumber": "0812345678",
    "birthday": "1990-01-15",
    "createdAt": "2025-08-27T14:00:00Z"
  },
  "expiresAt": "2025-08-28T14:00:00Z"
}
```

**Error Responses:**

*400 - Validation Failed:*
```json
{
  "error": "Validation failed",
  "message": "Field validation error details..."
}
```

*401 - Invalid Credentials:*
```json
{
  "error": "Invalid credentials",
  "message": "Email or password is incorrect"
}
```

**Example:**
```bash
curl -X POST http://localhost:3000/login \
-H "Content-Type: application/json" \
-d '{
  "email": "test@example.com",
  "password": "password123"
}'

## Built With

- [Go](https://golang.org/) - Programming language
- [Fiber](https://docs.gofiber.io/) - Web framework inspired by Express.js
- [bcrypt](https://pkg.go.dev/golang.org/x/crypto/bcrypt) - Password hashing
- [Validator](https://github.com/go-playground/validator) - Input validation
- [JWT](https://github.com/golang-jwt/jwt) - JSON Web Token implementation
- [Swagger](https://github.com/gofiber/swagger) - API documentation and testing
- [SQLite](https://modernc.org/sqlite) - Pure Go SQLite database driver

## Security Features

- Passwords are hashed using bcrypt before storage
- JWT tokens for secure authentication (24-hour expiry)
- Authorization header validation (Bearer token format)
- Token signature verification and claims validation
- Input validation prevents malformed data
- Email uniqueness validation
- Secure password requirements (minimum 6 characters)
- Credentials are never exposed in API responses
