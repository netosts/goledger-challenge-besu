package models

import (
	"testing"
)

func TestSetValueRequest_IsValid(t *testing.T) {
	tests := []struct {
		name    string
		value   uint64
		wantErr bool
	}{
		{
			name:    "valid small value",
			value:   42,
			wantErr: false,
		},
		{
			name:    "valid zero value",
			value:   0,
			wantErr: false,
		},
		{
			name:    "valid large value",
			value:   1000000,
			wantErr: false,
		},
		{
			name:    "valid max allowed value",
			value:   1e18,
			wantErr: false,
		},
		{
			name:    "invalid too large value",
			value:   1e18 + 1,
			wantErr: true,
		},
		{
			name:    "invalid very large value",
			value:   18446744073709551615, // max uint64
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &SetValueRequest{
				Value: tt.value,
			}
			err := r.IsValid()
			if (err != nil) != tt.wantErr {
				t.Errorf("SetValueRequest.IsValid() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestSetValueRequest_IsValid_ErrorMessage(t *testing.T) {
	r := &SetValueRequest{
		Value: 1e18 + 1,
	}
	err := r.IsValid()
	if err == nil {
		t.Error("SetValueRequest.IsValid() should return error for large value")
		return
	}

	expectedMsg := "value is too large"
	if err.Error() != expectedMsg {
		t.Errorf("SetValueRequest.IsValid() error message = %v, want %v", err.Error(), expectedMsg)
	}
}
