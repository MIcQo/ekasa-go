package ekasa

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestClient_GetCustomer(t *testing.T) {
	mockResponse := CustomerDto{
		ID:        "cust-123",
		FirstName: "John",
		LastName:  "Doe",
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("Expected GET request, got %s", r.Method)
		}
		if r.URL.Path != "/v1/customers/cust-123" {
			t.Errorf("Expected path /v1/customers/cust-123, got %s", r.URL.Path)
		}
		if r.Header.Get("Authorization") == "" {
			t.Error("Missing Authorization header")
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		if err := json.NewEncoder(w).Encode(mockResponse); err != nil {
			t.Errorf("Encode failed: %v", err)
		}
	}))
	defer server.Close()

	client := NewClient("pub", "priv", WithBaseURL(server.URL))
	customer, err := client.GetCustomer(context.Background(), "cust-123")

	if err != nil {
		t.Fatalf("GetCustomer failed: %v", err)
	}

	if customer.ID != mockResponse.ID {
		t.Errorf("Expected customer ID %s, got %s", mockResponse.ID, customer.ID)
	}
}

func TestClient_RegisterReceipt(t *testing.T) {
	mockResponse := ReceiptRegistrationDto{
		ID:    "reg-123",
		State: RegistrationStateProcessed,
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("Expected POST request, got %s", r.Method)
		}

		var req CreateReceiptRegistrationDto
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			t.Errorf("Failed to decode request body: %v", err)
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		if err := json.NewEncoder(w).Encode(mockResponse); err != nil {
			t.Errorf("Encode failed: %v", err)
		}
	}))
	defer server.Close()

	client := NewClient("pub", "priv", WithBaseURL(server.URL))
	req := &CreateReceiptRegistrationDto{
		Printer: ReceiptPrinterDto{Name: "pos"},
		Request: CreateReceiptRegistrationRequestDto{
			ExternalID: "ext-123",
		},
	}

	res, err := client.RegisterReceipt(context.Background(), req)
	if err != nil {
		t.Fatalf("RegisterReceipt failed: %v", err)
	}

	if res.ID != "reg-123" || res.State != RegistrationStateProcessed {
		t.Errorf("Unexpected response: %+v", res)
	}
}

func TestClient_HandleError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		pd := ProblemDetails{
			Title:  "Not Found",
			Detail: "Customer not found",
			Status: http.StatusNotFound,
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		if err := json.NewEncoder(w).Encode(pd); err != nil {
			t.Errorf("Encode failed: %v", err)
		}
	}))
	defer server.Close()

	client := NewClient("pub", "priv", WithBaseURL(server.URL))
	_, err := client.GetCustomer(context.Background(), "none")

	if err == nil {
		t.Error("Expected error, got nil")
	}
	if !contains(err.Error(), "api error (404)") {
		t.Errorf("Unexpected error message: %v", err)
	}
}

func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || (len(substr) > 0 && (s[:len(substr)] == substr || contains(s[1:], substr))))
}
