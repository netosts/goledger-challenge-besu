package integration

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/netosts/goledger-challenge-besu/internal/database"
	"github.com/netosts/goledger-challenge-besu/internal/handlers"
	"github.com/netosts/goledger-challenge-besu/internal/models"
	"github.com/netosts/goledger-challenge-besu/internal/repositories"
	"github.com/netosts/goledger-challenge-besu/internal/routes"
)

// These are integration tests that require a real database
// Run with: go test -v ./tests/integration/...
// Note: These tests require a running PostgreSQL database

func setupIntegrationTest(t *testing.T) (*gin.Engine, func()) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	// Set up test database connection
	dbConfig := &database.Config{
		Host:     "localhost",
		Port:     "5432",
		User:     "postgres",
		Password: "password",
		DBName:   "besu_challenge_test", // Use a separate test database
	}

	db, err := dbConfig.Connect()
	if err != nil {
		t.Skipf("Skipping integration test: failed to connect to database: %v", err)
	}

	// Initialize schema
	err = database.InitializeSchema(db)
	if err != nil {
		db.Close()
		t.Fatalf("Failed to initialize test database schema: %v", err)
	}

	// Clean up any existing data
	_, err = db.Exec("DELETE FROM contract_values")
	if err != nil {
		db.Close()
		t.Fatalf("Failed to clean test database: %v", err)
	}

	// Create repository and use case (without blockchain connection for integration tests)
	repo := repositories.NewPostgresRepository(db)

	// For integration tests, we'll create a mock use case that only tests database operations
	mockUseCase := &MockContractUseCaseIntegration{
		repo: repo,
	}

	handler := handlers.NewHandler(mockUseCase)
	router := routes.SetupRoutes(handler)

	// Return cleanup function
	cleanup := func() {
		db.Close()
	}

	return router, cleanup
}

// MockContractUseCaseIntegration mocks blockchain operations but uses real database
type MockContractUseCaseIntegration struct {
	repo            repositories.Repository
	blockchainValue uint64
}

func (m *MockContractUseCaseIntegration) SetValue(value uint64) error {
	m.blockchainValue = value
	return nil
}

func (m *MockContractUseCaseIntegration) GetValue() (uint64, error) {
	return m.blockchainValue, nil
}

func (m *MockContractUseCaseIntegration) SyncValue() error {
	return m.repo.SetValue(m.blockchainValue)
}

func (m *MockContractUseCaseIntegration) CheckValue() (*models.CheckResponse, error) {
	storedValue, err := m.repo.GetLatestValue()
	if err != nil {
		return nil, err
	}

	response := &models.CheckResponse{
		IsEqual:         m.blockchainValue == storedValue.Value,
		DatabaseValue:   storedValue.Value,
		BlockchainValue: m.blockchainValue,
	}

	return response, nil
}

func (m *MockContractUseCaseIntegration) Close() {
	// No-op for mock
}

func TestIntegration_CompleteWorkflow(t *testing.T) {
	router, cleanup := setupIntegrationTest(t)
	defer cleanup()

	// Test complete workflow: Set -> Get -> Sync -> Check

	// 1. Set a value
	setPayload := models.SetValueRequest{Value: 123}
	setBody, _ := json.Marshal(setPayload)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/v1/set", bytes.NewBuffer(setBody))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("Set request failed: status = %v, body = %s", w.Code, w.Body.String())
	}

	// 2. Get the value
	w = httptest.NewRecorder()
	req, _ = http.NewRequest("GET", "/api/v1/get", nil)
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("Get request failed: status = %v, body = %s", w.Code, w.Body.String())
	}

	var getValue models.ValueResponse
	err := json.Unmarshal(w.Body.Bytes(), &getValue)
	if err != nil {
		t.Fatalf("Failed to unmarshal get response: %v", err)
	}

	if getValue.Value != 123 {
		t.Errorf("Get value = %v, want 123", getValue.Value)
	}

	// 3. Sync to database
	w = httptest.NewRecorder()
	req, _ = http.NewRequest("POST", "/api/v1/sync", nil)
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("Sync request failed: status = %v, body = %s", w.Code, w.Body.String())
	}

	// 4. Check consistency
	w = httptest.NewRecorder()
	req, _ = http.NewRequest("GET", "/api/v1/check", nil)
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("Check request failed: status = %v, body = %s", w.Code, w.Body.String())
	}

	var checkResponse models.CheckResponse
	err = json.Unmarshal(w.Body.Bytes(), &checkResponse)
	if err != nil {
		t.Fatalf("Failed to unmarshal check response: %v", err)
	}

	if !checkResponse.IsEqual {
		t.Errorf("Check IsEqual = %v, want true", checkResponse.IsEqual)
	}

	if checkResponse.DatabaseValue != 123 {
		t.Errorf("Check DatabaseValue = %v, want 123", checkResponse.DatabaseValue)
	}

	if checkResponse.BlockchainValue != 123 {
		t.Errorf("Check BlockchainValue = %v, want 123", checkResponse.BlockchainValue)
	}
}

func TestIntegration_DatabasePersistence(t *testing.T) {
	router, cleanup := setupIntegrationTest(t)
	defer cleanup()

	// Test that values persist in database across multiple sync operations

	// Set initial value and sync
	setPayload := models.SetValueRequest{Value: 100}
	setBody, _ := json.Marshal(setPayload)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/v1/set", bytes.NewBuffer(setBody))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("Set request failed: %v", w.Code)
	}

	// Sync to database
	w = httptest.NewRecorder()
	req, _ = http.NewRequest("POST", "/api/v1/sync", nil)
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("First sync failed: %v", w.Code)
	}

	// Change blockchain value and sync again
	setPayload = models.SetValueRequest{Value: 200}
	setBody, _ = json.Marshal(setPayload)

	w = httptest.NewRecorder()
	req, _ = http.NewRequest("POST", "/api/v1/set", bytes.NewBuffer(setBody))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("Second set request failed: %v", w.Code)
	}

	// Sync again
	w = httptest.NewRecorder()
	req, _ = http.NewRequest("POST", "/api/v1/sync", nil)
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("Second sync failed: %v", w.Code)
	}

	// Check that database has the latest value
	w = httptest.NewRecorder()
	req, _ = http.NewRequest("GET", "/api/v1/check", nil)
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("Check request failed: %v", w.Code)
	}

	var checkResponse models.CheckResponse
	err := json.Unmarshal(w.Body.Bytes(), &checkResponse)
	if err != nil {
		t.Fatalf("Failed to unmarshal check response: %v", err)
	}

	if checkResponse.DatabaseValue != 200 {
		t.Errorf("DatabaseValue = %v, want 200", checkResponse.DatabaseValue)
	}

	if !checkResponse.IsEqual {
		t.Errorf("Values should be equal after sync, got IsEqual = %v", checkResponse.IsEqual)
	}
}

func TestIntegration_HealthCheck(t *testing.T) {
	router, cleanup := setupIntegrationTest(t)
	defer cleanup()

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/v1/health", nil)
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Health check failed: status = %v", w.Code)
	}

	var response models.SuccessResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	if err != nil {
		t.Fatalf("Failed to unmarshal health check response: %v", err)
	}

	if response.Message != "Service is healthy" {
		t.Errorf("Health check message = %v, want 'Service is healthy'", response.Message)
	}
}
