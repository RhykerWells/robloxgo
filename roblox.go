// Robloxgo - Roblox bindings for Go
// Available at https://github.com/RhykerWells/robloxgo
//
// Copyright 2025 Rhyker Wells <a.rhykerw@gmail.com>.  All rights reserved.
// License can be found in the LICENSE file of the repository.
//
// Package robloxgo provides Roblox binding for Go
package robloxgo

import (
	"net/http"
)

// Version of RobloxGo. Follows Semantic Versioning. (https://semver.org)
const Version = "1.0.0-alpha.1"

// Create initialises and returns a new Roblox client with the provided API key.
// The client automatically attaches the API key to all outgoing requests via the "X-API-KEY" header
//
// Returns an error if the API key is empty
func Create(apikey string) (*Client, error) {
	if apikey == "" {
		return nil, ErrNoAPIKey
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

// Client represents the created http client and will serve as a base for
// all help functions to be accessed from
type Client struct {
	client *http.Client
}

type APIVerificationStruct struct {
	APIKey    string
	Transport http.RoundTripper
}

func (a *APIVerificationStruct) RoundTrip(req *http.Request) (*http.Response, error) {
	req.Header.Set("X-API-KEY", a.APIKey)
	return a.Transport.RoundTrip(req)
}
