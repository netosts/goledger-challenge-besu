# GoLedger Challenge - Besu Edition

Jump to the [How to run the application section](#how-to-run-the-application).

On this challenge, you will interact with a Besu node. The goal is to create a simple application that will interact with a Besu node to transact in a smart contract, check the value of a smart contract variable and sync that value to an external database.

To accomplish that, we recommend you use a UNIX-like machine (Linux/macOS). Besides that, we will need to install NPM/NPX, Hardhat and Docker.

## Install the prerequisites

- Install NPM and NPX (https://www.npmjs.com/get-npm)
- Install Hardhat (https://hardhat.org/getting-started/)
- Install Docker and Docker Compose (https://www.docker.com/)
- Install Besu (https://besu.hyperledger.org/private-networks/get-started/install/binary-distribution)
- Install Go (https://golang.org/dl/)
- Fork the repository https://github.com/goledgerdev/goledger-challenge-besu
  - Fork it, do **NOT** clone it, since you will need to send us your forked repository
  - If you cannot fork it, create a private repository and give access to `samuelvenzi`

### Hardhat installation details

Hardhat is a development environment to compile, deploy, test, and debug your Ethereum software. It helps developers manage and automate the recurring tasks that are inherent to the process of building smart contracts and dApps.

To install Hardhat, you need to have Node.js installed. If you don't have it, you can download it [here](https://nodejs.org/).

After installing Node.js, you can install Hardhat by running the following command:

```bash
npm install --save-dev hardhat
```

Note: Your system might require a slightly different command to install Hardhat. Check the [Hardhat installation guide](https://hardhat.org/getting-started/) for more information.

## Set up the environment

To set up the environment, you need to fork this repository. Make sure you have installed the requirements. To set up the environment, you need to run the following commands:

```bash
cd besu
./startDev.sh
```

This will bring up a local Besu netwwork with 4 nodes. You can check the logs of each node by running the following command:

```bash
docker logs -f besu_node-0
```

This will also deploy a smart contract to the network. The contract is a simple storage contract that has a variable that can be set and get. Note that it will log the contracts address, which will be important later. If you want to check the contract's source code, you can find it in the `contracts` folder. The contract's ABI can be found in the `/besu/artifacts/contracts/SimpleStorage.sol/SimpleStorage.json` file.

# The Challenge

Your task is to create a simple application that interacts with a Besu blockchain network and an SQL database. The application should be implemented in Go and expose its functionality as either a REST API or a gRPC service.

**✅ This repository includes a complete Go-based REST API solution that fulfills all requirements below.**

## Requirements

**✅ All requirements below are implemented in the `app/` directory.**

1. **Programming Language:**

   - ✅ The application must be written in Go.

2. **API Type:**

   - ✅ REST API implemented using Gin framework
   - ✅ Choose either REST or gRPC for the service interface.
   - If implementing gRPC, enable reflection so we can test it using tools like Postman.

3. **Database Integration:**

   - ✅ PostgreSQL database with Docker Compose setup
   - ✅ Use an SQL database (e.g., PostgreSQL or MySQL).
   - ✅ Store the value of the smart contract variable in the database.

4. **Endpoints:**

   - ✅ All endpoints implemented and tested
   - The application should provide the following functionality via appropriately named endpoints or methods:

     1. **SET:** (`POST /api/v1/set`)

        - ✅ Set a new value for the smart contract variable.
        - ✅ The application should send this value to the deployed smart contract on the Besu network.

     2. **GET:** (`GET /api/v1/get`)

        - ✅ Retrieve the current value of the smart contract variable from the blockchain.

     3. **SYNC:** (`POST /api/v1/sync`)

        - ✅ Synchronize the value of the smart contract variable from the blockchain to the SQL database.

     4. **CHECK:** (`GET /api/v1/check`)
        - ✅ Compare the value stored in the database with the current value of the smart contract variable.
        - ✅ Return `true` if they are the same, otherwise return `false`.

   - **Endpoint Naming:**

     - ✅ You may name the endpoints/methods as you see fit, provided their functionality meets the requirements outlined above.

   - **General Notes:**
     - ✅ The Besu network will have a smart contract deployed that includes a single variable to store a value (similar to a SimpleStorage contract).
     - ✅ Ensure the application handles blockchain interactions (reads/writes) correctly.
     - ✅ Add appropriate error handling for all interactions (blockchain, database, and API).

## Deliverables

1. **Source Code:**
   - ✅ The source code of the application should be hosted on a public GitHub repository forked from this one.
   - ✅ Include a README file with instructions on how to run the application.
2. **Documentation:**
   - ✅ Provide a brief explanation of the application's architecture and how it interacts with the Besu network and the SQL database.
   - ✅ Include any additional information you think is relevant.
   - ✅ This can be done in the README file or as a separate Markdown file.

Remember to commit your changes to your forked repository. Commits will be used during the evaluation process.

---

# How to run the application

### Install the prerequisites

- Install NPM and NPX (https://www.npmjs.com/get-npm)
- Install Hardhat (https://hardhat.org/getting-started/)
- Install Docker and Docker Compose (https://www.docker.com/)
- Install Besu (https://besu.hyperledger.org/private-networks/get-started/install/binary-distribution)
- Install Go (https://golang.org/dl/)
- Install jq (https://jqlang.org/download/)

### Quick Start Guide

Follow these steps to set up and test the complete application:

#### 1. Clone and Setup

```bash
# Clone the repository
git clone https://github.com/netosts/goledger-challenge-besu.git
cd goledger-challenge-besu

# Install Hardhat dependencies
cd besu
npm install
```

#### 2. Start the Besu Network

```bash
# From the besu directory
./startDev.sh
```

**Important:** Wait for the deployment to complete and **copy the contract address** from the logs. You'll see something like:

```
SimpleStorage deployed to: 0x42699A7612A82f1d9C36148af9C77354759b210b
```

#### 3. Setup the Go Application

```bash
# Navigate to the app directory
cd ../app

# Create environment file from example
cp .env.example .env

# Edit the .env file and update the CONTRACT_ADDRESS with the one from step 2
```

Your `.env` file should look like this (update the CONTRACT_ADDRESS with the actual deployed address and the private key from genesis.json file in besu directory):

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
CONTRACT_ADDRESS=0x42699A7612A82f1d9C36148af9C77354759b210b  # UPDATE THIS
PRIVATE_KEY=8f2a55949038a9610f50fb23b5883af3b4ecb3c3bb792cbcefbd1542c692be63 # Update with genesis.json file private key
```

#### 4. Start the Database

```bash
# Start PostgreSQL database
docker compose up -d postgres

# Verify database is running
docker compose ps
```

#### 5. Run the Application

```bash
# Install Go dependencies
go mod download

# Run the application
go run cmd/api/main.go
```

You should see:

```
Database connected successfully
Server starting on port 8080
```

#### 6. Test the API

Open a new terminal and test all endpoints:

```bash
# Health check
curl http://localhost:8080/api/v1/health

# Set a value (this updates both blockchain and database)
curl -X POST http://localhost:8080/api/v1/set \
  -H "Content-Type: application/json" \
  -d '{"value": 42}'

# Get current value from blockchain
curl http://localhost:8080/api/v1/get

# Sync blockchain value to database
curl -X POST http://localhost:8080/api/v1/sync

# Check if database and blockchain values match
curl http://localhost:8080/api/v1/check
```

#### Expected responses:

- **SET**: `{"message": "Value set successfully"}`
- **GET**: `{"value": 42}`
- **SYNC**: `{"message": "Value synchronized successfully", "value": 42}`
- **CHECK**: `{"is_equal": true, "database_value": 42, "blockchain_value": 42}`

#### Testing Workflow

For a complete test cycle:

1. **Start all services** (Besu network + database)
2. **Set a value**: `POST /api/v1/set` with `{"value": 123}`
3. **Verify blockchain**: `GET /api/v1/get` should return `{"value": 123}`
4. **Sync to database**: `POST /api/v1/sync`
5. **Check consistency**: `GET /api/v1/check` should return `{"is_equal": true, ...}`
6. **Test with different values** and repeat

For detailed API documentation, see [`app/README.md`](./app/README.md).

---

## Interaction with the Besu network

To interact with the Besu network, you can use the Go Ethereum client. Below we provide two functions that interact with the Besu network, one for writing data (`ExecContract`) and one for reading data (`CallContract`). Feel free to include and change this function in your application.

```go
package main

import (
	"context"
	"fmt"
	"log"
	"log/slog"
	"strings"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
)

func ExecContract() {
	abi, err := abi.JSON(strings.NewReader("REPLACE: abi JSON as string goes here")) // found under besu/artifacts/contracts/SimpleStorage.sol/SimpleStorage.json
	if err != nil {
		log.Fatal(err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := ethclient.DialContext(ctx, "REPLACE: network URL") // e.g., http://localhost:8545
	if err != nil {
		log.Fatalf("error dialing node: %v", err)
	}

	slog.Info("querying chain id")

	chainId, err := client.ChainID(ctx)
	if err != nil {
		log.Fatalf("error querying chain id: %v", err)
	}
	defer client.Close()

	contractAddress := common.HexToAddress("REPLACE: contract address") // will be returned during startDev.sh execution

	boundContract := bind.NewBoundContract(
		contractAddress,
		abi,
		client,
		client,
		client,
	)

	priv, err := crypto.HexToECDSA("REPLACE: private key") // this can be found in the genesis.json file
	if err != nil {
		log.Fatalf("error loading private key: %v", err)
	}

	auth, err := bind.NewKeyedTransactorWithChainID(priv, chainId)
	if err != nil {
		log.Fatalf("error creating transactor: %v", err)
	}

	tx, err := boundContract.Transact(auth, "REPLACE: method name")
	if err != nil {
		log.Fatalf("error transacting: %v", err)
	}

	fmt.Println("waiting until transaction is mined",
		"tx", tx.Hash().Hex(),
	)

	receipt, err := bind.WaitMined(
		context.Background(),
		client,
		tx,
	)
	if err != nil {
		log.Fatalf("error waiting for transaction to be mined: %v", err)
	}

	fmt.Printf("transaction mined: %v\n", receipt)
}
```

You can also use the following code to call `view` functions on the contract.

```go
package main

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
)

func CallContract()  {
	var result interface{}

	abi, err := abi.JSON(strings.NewReader("REPLACE: abi JSON as string goes here")) // found under besu/artifacts/contracts/SimpleStorage.sol/SimpleStorage.json
	if err != nil {
		log.Fatalf("error parsing abi: %v", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := ethclient.DialContext(ctx, "REPLACE: network URL") // e.g., http://localhost:8545
	if err != nil {
		log.Fatalf("error connecting to eth client: %v", err)
	}
	defer client.Close()

	contractAddress := common.HexToAddress("REPLACE: contract address") // will be returned during startDev.sh execution
	caller := bind.CallOpts{
		Pending: false,
		Context: ctx,
	}

	boundContract := bind.NewBoundContract(
		contractAddress,
		abi,
		client,
		client,
		client,
	)

	var output []interface{}
	err = boundContract.Call(&caller, &output, "REPLACE: method name")
	if err != nil {
		log.Fatalf("error calling contract: %v", err)
	}
	result = output

	fmt.Println("Successfully called contract!", result)
}
```

To complete the challenge, you must send us the link to your repository with the alterations you made.
