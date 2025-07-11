package robloxgo

import (
	"os"
	"testing"
)

func TestGetUser_EmptyUserID(t *testing.T) {
	apiKey := os.Getenv("RG_APIKEY")
	client, _ := Create(apiKey)

	user, err := client.GetUserByID("")
	if err == nil {
		t.Fatal("expected error for empty userID, got nil")
	}
	if user != nil {
		t.Fatalf("expected nil user, got %v", user)
	}
}

func TestGetUser_InvalidUserID(t *testing.T) {
	apiKey := os.Getenv("RG_APIKEY")
	client, _ := Create(apiKey)

	user, err := client.GetUserByID("xxx")
	if err == nil {
		t.Fatalf("expected error for invalid userID, got %v", err)
	}
	if user != nil {
		t.Fatalf("unexpected user, got %v", user)
	}
}

func TestGetUser_PopulatedUserID(t *testing.T) {
	apiKey := os.Getenv("RG_APIKEY")
	client, _ := Create(apiKey)

	user, err := client.GetUserByID("369780411")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if user == nil {
		t.Fatal("expected user, got nil")
	}
}

func TestGetUser_EmptyUsername(t *testing.T) {
	apiKey := os.Getenv("RG_APIKEY")
	client, _ := Create(apiKey)

	user, err := client.GetUserByUsername("")
	if err == nil {
		t.Fatal("expected error for empty username, got nil")
	}
	if user != nil {
		t.Fatalf("expected nil user, got %v", user)
	}
}

func TestGetUser_InvalidUsername(t *testing.T) {
	apiKey := os.Getenv("RG_APIKEY")
	client, _ := Create(apiKey)

	user, err := client.GetUserByUsername("xxx")
	if err == nil {
		t.Fatal("expected error for invalid username, got nil")
	}
	if user != nil {
		t.Fatalf("expected nil user, got %v", user)
	}
}

func TestGetUser_PopulatedUsername(t *testing.T) {
	apiKey := os.Getenv("RG_APIKEY")
	client, _ := Create(apiKey)

	user, err := client.GetUserByUsername("captainbarborsa")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if user == nil {
		t.Fatal("expected user, got nil")
	}
}

func TestGetUserThumnail(t *testing.T) {
	apiKey := os.Getenv("RG_APIKEY")
	client, _ := Create(apiKey)

	user, err := client.GetUserByID("369780411")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if user == nil {
		t.Fatal("expected user, got nil")
	}

	uri, err := user.GetUserThumbnailURI(nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if uri == "" {
		t.Fatal("expected user, got nil")
	}
}
