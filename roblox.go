// Robloxgo - Roblox bindings for Go
// Available at https://github.com/RhykerWells/robloxgo

// Copyright 2025 Rhyker Wells <a.rhykerw@gmail.com>.  All rights reserved.
// License can be found in the LICENSE file of the repository.

// Package robloxgo provides Roblox binding for Go

package robloxgo

import (
	"net/http"
	"errors"
)

// VERSION of DiscordGo. Follows Semantic Versioning. (https://semver.org)
const VERSION = "0.1.0"

// Create creates a new Roblox client with the provided API key.
func Create(apikey string) (*Client, error) {
	if apikey == "" {
		return nil, errors.New("no api key provided")
	}

	httpClient := &http.Client{
		Transport: &APIVerificationStruct{
			APIKey:    apikey,
			Transport: http.DefaultTransport,
		},
	}

	client := &Client{
		client: httpClient,
	}

	return client, nil
}

type Client struct {
	client  *http.Client
}

type APIVerificationStruct struct {
	APIKey    string
	Transport http.RoundTripper
}

func (a *APIVerificationStruct) RoundTrip(req *http.Request) (*http.Response, error) {
	req.Header.Set("X-API-KEY", a.APIKey)
	return a.Transport.RoundTrip(req)
}