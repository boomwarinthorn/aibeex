# Fiber Hello World API

A simple Go web server using the Fiber framework that returns a JSON "Hello World" message.

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

## Built With

- [Go](https://golang.org/) - Programming language
- [Fiber](https://docs.gofiber.io/) - Web framework inspired by Express.js
