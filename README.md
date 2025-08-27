# Fiber Hello World API with Authentication

A Go web server using the Fiber framework that provides user authentication and registration features with JWT token support.

## Features

- ✅ Hello World JSON API endpoint
- ✅ User registration with validation
- ✅ User login with JWT token generation
- ✅ JWT-protected user profile endpoint
- ✅ Swagger/OpenAPI documentation
- ✅ SQLite database for persistent data storage
- ✅ Password hashing with bcrypt
- ✅ Email validation
- ✅ Input validation for all fields
- ✅ Duplicate email prevention
- ✅ JWT tokens expire after 24 hours
- ✅ Persistent data storage with SQLite database
- ✅ Database auto-initialization on first run

## Installation

1. Make sure you have Go 1.17 or higher installed
2. Clone this project
3. Install dependencies:
   ```bash
   go mod download
   ```

## Running the Application

```bash
go run main.go
```

The server will start on `http://localhost:3000`

## Swagger Documentation

Interactive API documentation is available at: `http://localhost:3000/swagger/index.html`

The Swagger UI provides:
- Complete API endpoint documentation
- Interactive testing interface
- Request/response schema definitions
- Authentication examples with Bearer tokens

### Regenerating Swagger Documentation

When you make changes to the API endpoints or add new ones, regenerate the Swagger docs:

```bash
go run github.com/swaggo/swag/cmd/swag@latest init
```

## Database

The application uses SQLite for persistent data storage. The database file (`users.db`) is automatically created in the project directory when you first run the application.

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
