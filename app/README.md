# Besu Blockchain Challenge - Go Application Documentation

This is a simple Go application that acts as a bridge between a blockchain network (Besu) and a database (PostgreSQL).

Imagine you have:

- A **smart contract** on the blockchain that stores a number
- A **database** on your computer that also stores a number
- This **Go application** that helps keep these numbers in sync

## How It Works

### Core Components

1. **REST API Server** ([`handlers/handlers.go`](app/internal/handlers/handlers.go))

   - Provides web endpoints you can call with simple HTTP requests
   - Handles incoming requests and returns responses in JSON format

2. **Business Logic** ([`usecases/usecases.go`](app/internal/usecases/usecases.go))

   - Contains the main logic for talking to the blockchain
   - Manages the coordination between blockchain and database operations

3. **Database Layer** ([`repositories/repositories.go`](app/internal/repositories/repositories.go))

   - Handles all database operations (get, update)
   - Keeps a local copy of the blockchain value

4. **Blockchain Communication**
   - Uses the Ethereum Go library to send transactions and read data
   - Connects to the Besu network running on your local machine

## How the App interacts with Blockchain

### Reading from Blockchain (GET operation)

```go
// This gets the current stored number from the smart contract
boundContract.Call(&caller, &output, "get")
```

### Writing to Blockchain (SET operation)

```go
// This updates the number stored in the smart contract
boundContract.Transact(auth, "set", big.NewInt(value))
```

## How the App interacts with PostgreSQL Database

The application uses a PostgreSQL database that is running in a Docker container to store a local copy of the blockchain value. This provides fast access to data without having to query the blockchain every time.

### Database Schema

The app creates a simple table called `contract_values`:

```sql
CREATE TABLE IF NOT EXISTS contract_values (
    id SERIAL PRIMARY KEY,
    value BIGINT NOT NULL,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
```

### Database Operations

#### 1. **Saving Values** (Used in SYNC operation)

```go
// Saves the blockchain value to the database
repo.SaveValue(value)
```

**What happens:**

- Inserts a new record with the current blockchain value
- Automatically timestamps when the sync occurred
- Keeps a history of all synced values

#### 2. **Getting Latest Value** (Used in CHECK operation)

```go
// Gets the most recent value from the database
value, err := repo.GetLatestValue()
```

**What happens:**

- Queries the database for the most recently stored value
- Returns the value so it can be compared with blockchain
- Used to check if database and blockchain are in sync

#### 3. **Database Connection**

The app connects to PostgreSQL using these environment variables:

```bash
DB_HOST=              # Where PostgreSQL is running
DB_PORT=              # PostgreSQL port
DB_USER=              # Database username
DB_PASSWORD=          # Database password
DB_NAME=              # Database name
```

## Why This Architecture?

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

1. **Clean Separation**: Each layer has one job

   - Handlers: Deal with HTTP requests/responses
   - Use Cases: Contain business logic
   - Repositories: Handle database operations

2. **Error Handling**: Each layer can handle errors appropriately

   - Network errors when talking to blockchain
   - Database connection issues
   - Invalid user input

3. **Testability**: Each component can be tested independently

### Components:

- **Handlers**: HTTP request/response handling and validation
- **Use Cases**: Core business logic for blockchain and database operations
- **Repositories**: Database abstraction layer with PostgreSQL implementation
- **Models**: Data structures for requests, responses, and database entities
- **Database**: Connection management and schema initialization
- **Routes**: API endpoint definitions and routing

## API Endpoints

### Base URL: `http://localhost:8080/api/v1`

| Method | Endpoint  | Description                           |
| ------ | --------- | ------------------------------------- |
| `POST` | `/set`    | Set new value in smart contract       |
| `GET`  | `/get`    | Get current value from blockchain     |
| `POST` | `/sync`   | Sync blockchain value to database     |
| `GET`  | `/check`  | Compare database vs blockchain values |
| `GET`  | `/health` | Health check                          |
