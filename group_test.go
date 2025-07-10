// Robloxgo - Roblox bindings for Go
// Available at https://github.com/RhykerWells/robloxgo
//
// Copyright 2025 Rhyker Wells <a.rhykerw@gmail.com>.  All rights reserved.
// License can be found in the LICENSE file of the repository.
//
// Package robloxgo provides Roblox binding for Go
package robloxgo

import (
	"os"
	"testing"
)

func TestGetGroup_EmptyGroupID(t *testing.T) {
	apiKey := os.Getenv("RG_APIKEY")
	client, _ := Create(apiKey)

	group, err := client.GetGroupByID("")
	if err == nil {
		t.Fatal("expected error for empty groupID, got nil")
	}
	if group != nil {
		t.Fatalf("expected nil group, got %v", group)
	}
}

func TestGetGroup_InvalidGroupID(t *testing.T) {
	apiKey := os.Getenv("RG_APIKEY")
	client, _ := Create(apiKey)

	group, err := client.GetGroupByID("xxx")
	if err == nil {
		t.Fatalf("expected error for invalid groupID, got %v", err)
	}
	if group != nil {
		t.Fatalf("unexpected group, got %v", group)
	}
}

func TestGetGroup_PopulatedGroupID(t *testing.T) {
	apiKey := os.Getenv("RG_APIKEY")
	client, _ := Create(apiKey)

	group, err := client.GetGroupByID("7")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if group == nil {
		t.Fatal("expected group, got nil")
	}
}

func TestGetGroup_EmptyGroupname(t *testing.T) {
	apiKey := os.Getenv("RG_APIKEY")
	client, _ := Create(apiKey)

	group, err := client.GetGroupByGroupname("")
	if err == nil {
		t.Fatal("expected error for empty groupname, got nil")
	}
	if group != nil {
		t.Fatalf("expected nil group, got %v", group)
	}
}

func TestGetGroup_InvalidGroupname(t *testing.T) {
	apiKey := os.Getenv("RG_APIKEY")
	client, _ := Create(apiKey)

	group, err := client.GetGroupByGroupname("xxx")
	if err == nil {
		t.Fatal("expected error for invalid groupname, got nil")
	}
	if group != nil {
		t.Fatalf("expected nil group, got %v", group)
	}
}

func TestGetGroup_PopulatedGroupname(t *testing.T) {
	apiKey := os.Getenv("RG_APIKEY")
	client, _ := Create(apiKey)

	group, err := client.GetGroupByGroupname("Roblox")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if group == nil {
		t.Fatal("expected group, got nil")
	}
}

func TestGetGroupJoinRequests(t *testing.T) {
	apiKey := os.Getenv("RG_APIKEY")
	client, _ := Create(apiKey)

	group, err := client.GetGroupByID("36098297")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if group == nil {
		t.Fatal("expected group, got nil")
	}

	requests, err := group.GetJoinRequests()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(requests) == 0 {
		t.Fatal("expected requests, got empty")
	}
}

func TestGetGroupUsers(t *testing.T) {
	apiKey := os.Getenv("RG_APIKEY")
	client, _ := Create(apiKey)

	group, err := client.GetGroupByID("7")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if group == nil {
		t.Fatal("expected group, got nil")
	}

	members, err := group.GetMembers()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(members) <= 0 {
		t.Fatalf("Expected members, got")
	}
}

func TestGetUsersLegacyRole(t *testing.T) {
	apiKey := os.Getenv("RG_APIKEY")
	client, _ := Create(apiKey)

	group, err := client.GetGroupByID("7")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if group == nil {
		t.Fatal("expected group, got nil")
	}

	role, err := group.GetUsersLegacyRole("21557")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if role == nil {
		t.Fatal("expected role, got nil")
	}
}

func TestGetGroupIcon(t *testing.T) {
	apiKey := os.Getenv("RG_APIKEY")
	client, _ := Create(apiKey)

	group, err := client.GetGroupByID("7")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if group == nil {
		t.Fatal("expected group, got nil")
	}

	url, err := group.GetGroupIcon(false, false)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if url == "" {
		t.Fatal("expected link, got nil")
	}
}
