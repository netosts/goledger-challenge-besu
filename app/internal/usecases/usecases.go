package usecases

import (
	"context"
	"fmt"
	"log"
	"math/big"
	"os"
	"strings"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/netosts/goledger-challenge-besu/internal/models"
	"github.com/netosts/goledger-challenge-besu/internal/repositories"
)

type ContractUseCase struct {
	repo repositories.Repository
}

func NewContractUseCase(repo repositories.Repository) *ContractUseCase {
	return &ContractUseCase{
		repo: repo,
	}
}

// SimpleStorage ABI - this should match the deployed contract
const SimpleStorageABI = `[
	{
		"inputs": [],
		"name": "get",
		"outputs": [
			{
				"internalType": "uint256",
				"name": "",
				"type": "uint256"
			}
		],
		"stateMutability": "view",
		"type": "function"
	},
	{
		"inputs": [
			{
				"internalType": "uint256",
				"name": "x",
				"type": "uint256"
			}
		],
		"name": "set",
		"outputs": [],
		"stateMutability": "nonpayable",
		"type": "function"
	}
]`

func (uc *ContractUseCase) SetValue(value uint64) error {
	if err := uc.setBlockchainValue(value); err != nil {
		return fmt.Errorf("failed to set blockchain value: %w", err)
	}

	// if err := uc.repo.SetValue(value); err != nil {
	// 	return fmt.Errorf("failed to update database: %w", err)
	// }

	return nil
}

func (uc *ContractUseCase) GetValue() (uint64, error) {
	return uc.getBlockchainValue()
}

func (uc *ContractUseCase) SyncValue() error {
	blockchainValue, err := uc.getBlockchainValue()
	if err != nil {
		return fmt.Errorf("failed to get blockchain value: %w", err)
	}

	if err := uc.repo.SetValue(blockchainValue); err != nil {
		return fmt.Errorf("failed to sync to database: %w", err)
	}

	return nil
}

func (uc *ContractUseCase) CheckValue() (*models.CheckResponse, error) {
	blockchainValue, err := uc.getBlockchainValue()
	if err != nil {
		return nil, fmt.Errorf("failed to get blockchain value: %w", err)
	}

	storedValue, err := uc.repo.GetLatestValue()
	if err != nil {
		return nil, fmt.Errorf("failed to get database value: %w", err)
	}

	response := &models.CheckResponse{
		IsEqual:         blockchainValue == storedValue.Value,
		DatabaseValue:   storedValue.Value,
		BlockchainValue: blockchainValue,
	}

	return response, nil
}

func (uc *ContractUseCase) setBlockchainValue(value uint64) error {
	abi, err := abi.JSON(strings.NewReader(SimpleStorageABI))
	if err != nil {
		return fmt.Errorf("failed to parse ABI: %w", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	nodeURL := getEnvOrDefault("NODE_URL", "http://localhost:8545")
	client, err := ethclient.DialContext(ctx, nodeURL)
	if err != nil {
		return fmt.Errorf("failed to connect to Ethereum client: %w", err)
	}
	defer client.Close()

	chainID, err := client.ChainID(ctx)
	if err != nil {
		return fmt.Errorf("failed to get chain ID: %w", err)
	}

	contractAddress := common.HexToAddress(getEnvOrDefault("CONTRACT_ADDRESS", ""))
	if contractAddress == (common.Address{}) {
		return fmt.Errorf("CONTRACT_ADDRESS environment variable is required")
	}

	boundContract := bind.NewBoundContract(
		contractAddress,
		abi,
		client,
		client,
		client,
	)

	privateKeyHex := getEnvOrDefault("PRIVATE_KEY", "")
	if privateKeyHex == "" {
		return fmt.Errorf("PRIVATE_KEY environment variable is required")
	}

	priv, err := crypto.HexToECDSA(privateKeyHex)
	if err != nil {
		return fmt.Errorf("failed to parse private key: %w", err)
	}

	auth, err := bind.NewKeyedTransactorWithChainID(priv, chainID)
	if err != nil {
		return fmt.Errorf("failed to create transactor: %w", err)
	}

	auth.GasLimit = uint64(300000)

	log.Printf("Setting blockchain value to %d", value)
	tx, err := boundContract.Transact(auth, "set", big.NewInt(int64(value)))
	if err != nil {
		return fmt.Errorf("failed to execute transaction: %w", err)
	}

	fmt.Println("waiting until transaction is mined",
		"tx", tx.Hash().Hex(),
	)

	receipt, err := bind.WaitMined(ctx, client, tx)
	if err != nil {
		return fmt.Errorf("failed to wait for transaction to be mined: %w", err)
	}

	log.Printf("Transaction mined in block %d", receipt.BlockNumber.Uint64())
	return nil
}

func (uc *ContractUseCase) getBlockchainValue() (uint64, error) {
	abi, err := abi.JSON(strings.NewReader(SimpleStorageABI))
	if err != nil {
		return 0, fmt.Errorf("failed to parse ABI: %w", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	nodeURL := getEnvOrDefault("NODE_URL", "http://localhost:8545")
	client, err := ethclient.DialContext(ctx, nodeURL)
	if err != nil {
		return 0, fmt.Errorf("failed to connect to Ethereum client: %w", err)
	}
	defer client.Close()

	contractAddress := common.HexToAddress(getEnvOrDefault("CONTRACT_ADDRESS", ""))
	if contractAddress == (common.Address{}) {
		return 0, fmt.Errorf("CONTRACT_ADDRESS environment variable is required")
	}

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
	err = boundContract.Call(&caller, &output, "get")
	if err != nil {
		return 0, fmt.Errorf("failed to call contract: %w", err)
	}

	if len(output) == 0 {
		return 0, fmt.Errorf("no output from contract call")
	}

	result, ok := output[0].(*big.Int)
	if !ok {
		return 0, fmt.Errorf("unexpected output type from contract call")
	}

	return result.Uint64(), nil
}

func getEnvOrDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
