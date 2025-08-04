# Sample Game Backend API

This project is a Golang sample backend implementation that manages session-specific asset information using **go-memdb** as an in-memory database.

## Installation and Setup

### Prerequisites
- Go 1.21 or higher
- Git

### 1. Clone the repository
```bash
git clone <repository-url>
cd sample-game-backend
```

### 2. Install dependencies
```bash
go mod tidy
```

### 3. Run the server
```bash
go run main.go
```

The server will start on port 8080. Session-specific asset information is stored in the **go-memdb** in-memory database.

## Project Structure

```
sample-game-backend/
├── main.go                 # Application entry point
├── internal/
│   ├── config/            # Configuration management
│   ├── database/          # Database operations (go-memdb)
│   ├── handlers/          # HTTP request handlers
│   ├── middleware/        # HTTP middleware (auth, CORS)
│   ├── models/            # Data structures
│   └── services/          # Business logic
├── test/                  # Test files
└── session_db/            # Session database files
```

## Development

### Running Tests
```bash
go test ./...
```

### Building
```bash
go build -o sample-game-backend
```

## Support

For questions or issues, please create an issue in the repository.
