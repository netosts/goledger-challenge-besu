# Besu Blockchain Challenge - Go Application

A Go-based REST API application that interacts with a Besu blockchain network and PostgreSQL database to manage smart contract values.

## Architecture Overview

The application follows a clean architecture pattern:

```
app/
├── cmd/               # Application entry point
├── internal/
│   ├── handlers/          # HTTP request handlers (REST API layer)
│   ├── usecases/          # Business logic layer
│   ├── repositories/      # Data access layer
│   ├── models/           # Data structures and DTOs
│   ├── database/         # Database connection and schema
│   └── routes/           # API route definitions
└── docker-compose.yml    # Docker configuration for dependencies
```

### Components:

- **Handlers**: HTTP request/response handling and validation
- **Use Cases**: Core business logic for blockchain and database operations
- **Repositories**: Database abstraction layer with PostgreSQL implementation
- **Models**: Data structures for requests, responses, and database entities
- **Database**: Connection management and schema initialization
- **Routes**: API endpoint definitions and routing

## Features

Four main endpoints that fulfill the challenge requirements:

1. **SET** (`POST /api/v1/set`) - Set value in smart contract and database
2. **GET** (`GET /api/v1/get`) - Retrieve current value from blockchain
3. **SYNC** (`POST /api/v1/sync`) - Synchronize blockchain value to database
4. **CHECK** (`GET /api/v1/check`) - Compare database and blockchain values

## Quick Start

### 1. Start Dependencies

```bash
# Start Besu network (in project root)
cd ../besu && ./startDev.sh

# Start PostgreSQL database
cd app && docker-compose up -d postgres
```

### 2. Configure Environment

Copy the contract address from Besu deployment logs to your `.env` file:

```bash
# Database Configuration
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=password
DB_NAME=besu_challenge

# Server Configuration
PORT=8080

# Blockchain Configuration
NODE_URL=http://localhost:8545
PRIVATE_KEY="8f2a55949038a9610f50fb23b5883af3b4ecb3c3bb792cbcefbd1542c692be63" # Update with genesis.json file private key
CONTRACT_ADDRESS="0x42699A7612A82f1d9C36148af9C77354759b210b"  # Update from Besu logs
```

### 3. Run Application

```bash
# Install dependencies and run
go mod download
go run cmd/api/main.go
```

API available at: `http://localhost:8080`

## API Endpoints

### Base URL: `http://localhost:8080/api/v1`

| Method | Endpoint  | Description                           |
| ------ | --------- | ------------------------------------- |
| `POST` | `/set`    | Set new value in smart contract       |
| `GET`  | `/get`    | Get current value from blockchain     |
| `POST` | `/sync`   | Sync blockchain value to database     |
| `GET`  | `/check`  | Compare database vs blockchain values |
| `GET`  | `/health` | Health check                          |

### Example Usage

```bash
# Set value
curl -X POST localhost:8080/api/v1/set \
  -H "Content-Type: application/json" \
  -d '{"value": 42}'

# Get value
curl localhost:8080/api/v1/get

# Sync values
curl -X POST localhost:8080/api/v1/sync

# Check consistency
curl localhost:8080/api/v1/check
```

### Response Examples

**SET Response:**

```json
{ "message": "Value set successfully" }
```

**GET Response:**

```json
{ "value": 42 }
```

**CHECK Response:**

```json
{
  "is_equal": true,
  "database_value": 42,
  "blockchain_value": 42
}
```

## Database Schema

```sql
CREATE TABLE stored_values (
    id SERIAL PRIMARY KEY,
    value BIGINT NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);
```

## Technology Stack

- **Go 1.21+** with Gin web framework
- **PostgreSQL** for data persistence
- **go-ethereum** for blockchain interaction
- **Docker Compose** for database setup

## Important Notes

1. **Prerequisites**:

   - Running Besu network with deployed SimpleStorage contract
   - PostgreSQL database (via Docker Compose)
   - Correct contract address and private key in `.env`

2. **Testing Workflow**:
   - Start Besu network and database
   - Update `.env` with contract address from logs
   - Update `.env` with private key from genesis.json file
   - Run application and test endpoints in sequence: SET → GET → CHECK → SYNC → CHECK
