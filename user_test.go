// Robloxgo - Roblox bindings for Go
// Available at https://github.com/RhykerWells/robloxgo

// Copyright 2025 Rhyker Wells <a.rhykerw@gmail.com>.  All rights reserved.
// License can be found in the LICENSE file of the repository.

// Package robloxgo provides Roblox binding for Go

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
