# Fiber Hello World API with Authentication

A Go web server using the Fiber framework that provides user authentication and registration features.

## Features

- ✅ Hello World JSON API endpoint
- ✅ User registration with validation
- ✅ Password hashing with bcrypt
- ✅ Email validation
- ✅ Input validation for all fields
- ✅ Duplicate email prevention

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

## Built With

- [Go](https://golang.org/) - Programming language
- [Fiber](https://docs.gofiber.io/) - Web framework inspired by Express.js
- [bcrypt](https://pkg.go.dev/golang.org/x/crypto/bcrypt) - Password hashing
- [Validator](https://github.com/go-playground/validator) - Input validation

## Security Features

- Passwords are hashed using bcrypt before storage
- Input validation prevents malformed data
- Email uniqueness validation
- Secure password requirements (minimum 6 characters)
