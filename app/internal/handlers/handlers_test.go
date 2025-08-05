package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/netosts/goledger-challenge-besu/internal/models"
)

// MockContractUseCase is a mock implementation of usecases.ContractUseCaseInterface for testing
type MockContractUseCase struct {
	setValue      uint64
	getValue      uint64
	syncError     error
	setError      error
	getError      error
	checkResponse *models.CheckResponse
	checkError    error
}

func (m *MockContractUseCase) SetValue(value uint64) error {
	if m.setError != nil {
		return m.setError
	}
	m.setValue = value
	return nil
}

func (m *MockContractUseCase) GetValue() (uint64, error) {
	if m.getError != nil {
		return 0, m.getError
	}
	return m.getValue, nil
}

func (m *MockContractUseCase) SyncValue() error {
	return m.syncError
}

func (m *MockContractUseCase) CheckValue() (*models.CheckResponse, error) {
	if m.checkError != nil {
		return nil, m.checkError
	}
	return m.checkResponse, nil
}

func (m *MockContractUseCase) Close() {
	// No-op for mock
}

func setupTestRouter(mockUseCase *MockContractUseCase) *gin.Engine {
	gin.SetMode(gin.TestMode)
	handler := NewHandler(mockUseCase)

	router := gin.New()
	v1 := router.Group("/api/v1")
	{
		v1.GET("/health", handler.HealthCheck)
		v1.POST("/set", handler.SetValue)
		v1.GET("/get", handler.GetValue)
		v1.POST("/sync", handler.SyncValue)
		v1.GET("/check", handler.CheckValue)
	}

	return router
}

func TestHandler_HealthCheck(t *testing.T) {
	mockUseCase := &MockContractUseCase{}
	router := setupTestRouter(mockUseCase)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/v1/health", nil)
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("HealthCheck() status = %v, want %v", w.Code, http.StatusOK)
	}

	var response models.SuccessResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	if err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}

	expectedMessage := "Service is healthy"
	if response.Message != expectedMessage {
		t.Errorf("HealthCheck() message = %v, want %v", response.Message, expectedMessage)
	}
}

func TestHandler_SetValue_Success(t *testing.T) {
	mockUseCase := &MockContractUseCase{}
	router := setupTestRouter(mockUseCase)

	payload := models.SetValueRequest{Value: 42}
	jsonPayload, _ := json.Marshal(payload)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/v1/set", bytes.NewBuffer(jsonPayload))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("SetValue() status = %v, want %v", w.Code, http.StatusOK)
	}

	var response models.SuccessResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	if err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}

	expectedMessage := "Value set successfully"
	if response.Message != expectedMessage {
		t.Errorf("SetValue() message = %v, want %v", response.Message, expectedMessage)
	}

	if mockUseCase.setValue != 42 {
		t.Errorf("SetValue() called with %v, want 42", mockUseCase.setValue)
	}
}

func TestHandler_SetValue_InvalidJSON(t *testing.T) {
	mockUseCase := &MockContractUseCase{}
	router := setupTestRouter(mockUseCase)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/v1/set", bytes.NewBuffer([]byte("invalid json")))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("SetValue() status = %v, want %v", w.Code, http.StatusBadRequest)
	}

	var response models.ErrorResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	if err != nil {
		t.Fatalf("Failed to unmarshal error response: %v", err)
	}

	if response.Error == "" {
		t.Error("SetValue() should return error message for invalid JSON")
	}
}

func TestHandler_SetValue_InvalidValue(t *testing.T) {
	mockUseCase := &MockContractUseCase{}
	router := setupTestRouter(mockUseCase)

	payload := models.SetValueRequest{Value: 1e18 + 1} // Too large
	jsonPayload, _ := json.Marshal(payload)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/v1/set", bytes.NewBuffer(jsonPayload))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("SetValue() status = %v, want %v", w.Code, http.StatusBadRequest)
	}

	var response models.ErrorResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	if err != nil {
		t.Fatalf("Failed to unmarshal error response: %v", err)
	}

	expectedError := "value is too large"
	if response.Error != expectedError {
		t.Errorf("SetValue() error = %v, want %v", response.Error, expectedError)
	}
}

func TestHandler_SetValue_UseCaseError(t *testing.T) {
	mockUseCase := &MockContractUseCase{
		setError: fmt.Errorf("blockchain connection failed"),
	}
	router := setupTestRouter(mockUseCase)

	payload := models.SetValueRequest{Value: 42}
	jsonPayload, _ := json.Marshal(payload)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/v1/set", bytes.NewBuffer(jsonPayload))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	if w.Code != http.StatusInternalServerError {
		t.Errorf("SetValue() status = %v, want %v", w.Code, http.StatusInternalServerError)
	}

	var response models.ErrorResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	if err != nil {
		t.Fatalf("Failed to unmarshal error response: %v", err)
	}

	if response.Error == "" {
		t.Error("SetValue() should return error message when use case fails")
	}
}

func TestHandler_GetValue_Success(t *testing.T) {
	mockUseCase := &MockContractUseCase{
		getValue: 123,
	}
	router := setupTestRouter(mockUseCase)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/v1/get", nil)
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("GetValue() status = %v, want %v", w.Code, http.StatusOK)
	}

	var response models.ValueResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	if err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}

	if response.Value != 123 {
		t.Errorf("GetValue() value = %v, want 123", response.Value)
	}
}

func TestHandler_GetValue_Error(t *testing.T) {
	mockUseCase := &MockContractUseCase{
		getError: fmt.Errorf("failed to read from blockchain"),
	}
	router := setupTestRouter(mockUseCase)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/v1/get", nil)
	router.ServeHTTP(w, req)

	if w.Code != http.StatusInternalServerError {
		t.Errorf("GetValue() status = %v, want %v", w.Code, http.StatusInternalServerError)
	}

	var response models.ErrorResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	if err != nil {
		t.Fatalf("Failed to unmarshal error response: %v", err)
	}

	if response.Error == "" {
		t.Error("GetValue() should return error message when use case fails")
	}
}

func TestHandler_SyncValue_Success(t *testing.T) {
	mockUseCase := &MockContractUseCase{}
	router := setupTestRouter(mockUseCase)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/v1/sync", nil)
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("SyncValue() status = %v, want %v", w.Code, http.StatusOK)
	}

	var response models.SuccessResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	if err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}

	expectedMessage := "Value synchronized successfully"
	if response.Message != expectedMessage {
		t.Errorf("SyncValue() message = %v, want %v", response.Message, expectedMessage)
	}
}

func TestHandler_SyncValue_Error(t *testing.T) {
	mockUseCase := &MockContractUseCase{
		syncError: fmt.Errorf("database connection failed"),
	}
	router := setupTestRouter(mockUseCase)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/v1/sync", nil)
	router.ServeHTTP(w, req)

	if w.Code != http.StatusInternalServerError {
		t.Errorf("SyncValue() status = %v, want %v", w.Code, http.StatusInternalServerError)
	}

	var response models.ErrorResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	if err != nil {
		t.Fatalf("Failed to unmarshal error response: %v", err)
	}

	if response.Error == "" {
		t.Error("SyncValue() should return error message when use case fails")
	}
}

func TestHandler_CheckValue_Success(t *testing.T) {
	mockResponse := &models.CheckResponse{
		IsEqual:         true,
		DatabaseValue:   42,
		BlockchainValue: 42,
	}
	mockUseCase := &MockContractUseCase{
		checkResponse: mockResponse,
	}
	router := setupTestRouter(mockUseCase)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/v1/check", nil)
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("CheckValue() status = %v, want %v", w.Code, http.StatusOK)
	}

	var response models.CheckResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	if err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}

	if response.IsEqual != true {
		t.Errorf("CheckValue() IsEqual = %v, want true", response.IsEqual)
	}
	if response.DatabaseValue != 42 {
		t.Errorf("CheckValue() DatabaseValue = %v, want 42", response.DatabaseValue)
	}
	if response.BlockchainValue != 42 {
		t.Errorf("CheckValue() BlockchainValue = %v, want 42", response.BlockchainValue)
	}
}

func TestHandler_CheckValue_Error(t *testing.T) {
	mockUseCase := &MockContractUseCase{
		checkError: fmt.Errorf("failed to compare values"),
	}
	router := setupTestRouter(mockUseCase)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/v1/check", nil)
	router.ServeHTTP(w, req)

	if w.Code != http.StatusInternalServerError {
		t.Errorf("CheckValue() status = %v, want %v", w.Code, http.StatusInternalServerError)
	}

	var response models.ErrorResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	if err != nil {
		t.Fatalf("Failed to unmarshal error response: %v", err)
	}

	if response.Error == "" {
		t.Error("CheckValue() should return error message when use case fails")
	}
}
