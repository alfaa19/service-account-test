# Banking Service API

A simple banking service REST API built with Go that handles basic banking operations like account creation, deposits, and withdrawals.

## Features

- Create new bank accounts
- Check account balance
- Deposit money
- Withdraw money
- Data persistence using PostgreSQL
- Structured logging
- Docker support
- API documentation

## Prerequisites

- Docker and Docker Compose
- Go 1.21 or higher (for local development)
- Make (optional, for using Makefile commands)

## Getting Started

1. Clone the repository
```bash
git clone https://github.com/yourusername/service-account-test.git
cd service-account-test
```

2. Create `.env` file
```bash
cp .env.example .env
```

3. Start the services
```bash
docker compose up --build
```

The API will be available at `http://localhost:8080` for default

## API Endpoints

### Create Account
```http
POST /daftar
Content-Type: application/json

{
    "nama": "John Doe",
    "nik": "1234567890",
    "no_hp": "081234567890"
}
```

### Check Balance
```http
GET /saldo/:noRekening
```

### Deposit Money
```http
POST /tabung
Content-Type: application/json

{
    "no_rekening": "1234567890",
    "saldo": 100000
}
```

### Withdraw Money
```http
POST /tarik
Content-Type: application/json

{
    "no_rekening": "1234567890",
    "saldo": 50000
}
```

## Project Structure

```
.
├── cmd/
│   └── server/
│       └── main.go
├── config/
│   └── config.go
├── internal/
│   ├── handler/
│   ├── models/
│   ├── repository/
│   ├── routes/
│   └── service/
├── migrations/
│   └── init.sql
├── pkg/
│   └── logrus/
├── docker-compose.yml
├── Dockerfile
└── README.md
```

## Configuration

Configuration is handled through environment variables:

| Variable | Description | Default |
|----------|-------------|---------|
| DB_HOST | Database host | db |
| DB_PORT | Database port | 5432 |
| DB_NAME | Database name | account_service |
| DB_USER | Database user | postgres |
| DB_PASSWORD | Database password | postgres |
| LOG_LEVEL | Logging level | DEBUG |

## Development

For local development:

1. Install dependencies
```bash
go mod download
```

2. Run PostgreSQL (using Docker)
```bash
docker compose up db
```

3. Run the application
```bash
go run cmd/server/main.go
```

## Testing

Run the tests:
```bash
go test ./...
```
