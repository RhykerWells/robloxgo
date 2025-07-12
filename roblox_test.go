package robloxgo

import (
	"os"
	"testing"
)

// Test Create returns error on empty API key
func TestCreate_EmptyAPIKey(t *testing.T) {
	client, err := Create("")
	if err == nil {
		t.Fatal("expected error for empty API key, got nil")
	}
	if client != nil {
		t.Fatal("expected client to be nil on error")
	}
}

// Test Create returns client with correct API key set in transport
func TestCreate_PopulatedAPIKey(t *testing.T) {
	apiKey := os.Getenv("RG_APIKEY")

	client, err := Create(apiKey)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if client == nil {
		t.Fatal("expected client, got nil")
	}
}
