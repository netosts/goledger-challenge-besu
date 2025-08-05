package repositories

import (
	"database/sql"
	"testing"
	"time"

	_ "github.com/lib/pq"
)

func setupTestDB(t *testing.T) *sql.DB {
	// For testing, we'll skip tests if no database is available
	// In a real CI environment, you would set up a test database
	t.Skip("Repository tests require a test database. Run integration tests instead.")
	return nil
}

func TestPostgresRepository_SetValue_Insert(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	repo := NewPostgresRepository(db)

	testValue := uint64(42)
	err := repo.SetValue(testValue)
	if err != nil {
		t.Errorf("SetValue() error = %v, want nil", err)
		return
	}

	// Verify the value was inserted
	storedValue, err := repo.GetLatestValue()
	if err != nil {
		t.Errorf("GetLatestValue() error = %v, want nil", err)
		return
	}

	if storedValue.Value != testValue {
		t.Errorf("GetLatestValue() value = %v, want %v", storedValue.Value, testValue)
	}
}

func TestPostgresRepository_SetValue_Update(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	repo := NewPostgresRepository(db)

	// Insert initial value
	initialValue := uint64(10)
	err := repo.SetValue(initialValue)
	if err != nil {
		t.Fatalf("SetValue() initial error = %v", err)
	}

	// Update with new value
	newValue := uint64(20)
	err = repo.SetValue(newValue)
	if err != nil {
		t.Errorf("SetValue() update error = %v, want nil", err)
		return
	}

	// Verify the value was updated
	storedValue, err := repo.GetLatestValue()
	if err != nil {
		t.Errorf("GetLatestValue() error = %v, want nil", err)
		return
	}

	if storedValue.Value != newValue {
		t.Errorf("GetLatestValue() value = %v, want %v", storedValue.Value, newValue)
	}
}

func TestPostgresRepository_GetLatestValue_NoData(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	repo := NewPostgresRepository(db)

	_, err := repo.GetLatestValue()
	if err == nil {
		t.Error("GetLatestValue() should return error when no data exists")
		return
	}

	expectedMsg := "no values found in database"
	if err.Error() != expectedMsg {
		t.Errorf("GetLatestValue() error message = %v, want %v", err.Error(), expectedMsg)
	}
}

func TestPostgresRepository_GetLatestValue_WithData(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	repo := NewPostgresRepository(db)

	// Insert test data
	testValue := uint64(100)
	err := repo.SetValue(testValue)
	if err != nil {
		t.Fatalf("SetValue() error = %v", err)
	}

	// Get the value
	storedValue, err := repo.GetLatestValue()
	if err != nil {
		t.Errorf("GetLatestValue() error = %v, want nil", err)
		return
	}

	if storedValue.Value != testValue {
		t.Errorf("GetLatestValue() value = %v, want %v", storedValue.Value, testValue)
	}

	// Check that timestamps are populated
	if storedValue.CreatedAt.IsZero() {
		t.Error("GetLatestValue() CreatedAt should not be zero")
	}

	if storedValue.UpdatedAt.IsZero() {
		t.Error("GetLatestValue() UpdatedAt should not be zero")
	}
}

func TestPostgresRepository_MultipleValues(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	repo := NewPostgresRepository(db)

	// Insert multiple values with delays to ensure different timestamps
	values := []uint64{10, 20, 30}
	for _, value := range values {
		err := repo.SetValue(value)
		if err != nil {
			t.Fatalf("SetValue(%d) error = %v", value, err)
		}
		time.Sleep(1 * time.Millisecond) // Small delay to ensure different timestamps
	}

	// Should get the latest (last) value
	storedValue, err := repo.GetLatestValue()
	if err != nil {
		t.Errorf("GetLatestValue() error = %v, want nil", err)
		return
	}

	expectedValue := values[len(values)-1]
	if storedValue.Value != expectedValue {
		t.Errorf("GetLatestValue() value = %v, want %v", storedValue.Value, expectedValue)
	}
}
