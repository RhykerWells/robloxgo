// Robloxgo - Roblox bindings for Go
// Available at https://github.com/RhykerWells/robloxgo
//
// Copyright 2025 Rhyker Wells <a.rhykerw@gmail.com>.  All rights reserved.
// License can be found in the LICENSE file of the repository.
//
// Package robloxgo provides Roblox binding for Go
package robloxgo

import (
	"encoding/json"
	"errors"
)

// A User stores all data for an individual Roblox user.
type User struct {
	// The ID of the user.
	ID json.Number `json:"id"`

	// The user's username.
	Username string `json:"name"`

	// The user's display name, if it is set.
	Displayname string `json:"displayName"`

	// The user's premium status.
	Premium bool `json:"premium"`

	// The user's chosen language option.
	Locale string `json:"locale"`

	// The user's account creation date.
	CreatedAt string `json:"createTime"`
}

// GetUserByID retrieves a Roblox user from the Open Cloud API by their user ID.
//
// Returns an error if the HTTP request fails, if the response body cannot
// be decoded, or if the user does not exist.
func (c *Client) GetUserByID(userID string) (*User, error) {
	response, err := c.get(EndPointCloudUsers + userID)
	if err != nil {
		return nil, err
	}

	user := new(User)
	err = json.NewDecoder(response.Body).Decode(user)
	if err != nil {
		return nil, err
	}

	return user, nil
}

// GetUserByUsername retrieves a Roblox user from the Legacy Roblox API by their user Username.
//
// Returns an error if the HTTP request fails, if the response body cannot
// be decoded, or if the user does not exist.
//
// This method may be deprecated if Roblox removes the
// legacy https://users.roblox.com/v1/usernames/users endpoint
func (c *Client) GetUserByUsername(username string) (*User, error) {
	if username == "" {
		return nil, errors.New("no username")
	}

	requestBody := map[string]interface{}{"usernames": []string{username}, "excludeBannedUsers": true}
	response, err := c.post(EndpointLegacyGetUsers, requestBody)
	if err != nil {
		return nil, err
	}

	var Response struct {
		Data []User `json:"data"`
	}
	err = json.NewDecoder(response.Body).Decode(&Response)
	if err != nil {
		return nil, err
	}

	if len(Response.Data) == 0 {
		return nil, errors.New("invalid username provided")
	}

	legacyUser := &Response.Data[0]

	response, err = c.get(EndPointCloudUsers + legacyUser.ID.String())
	if err != nil {
		return nil, err
	}

	user := new(User)
	err = json.NewDecoder(response.Body).Decode(user)
	if err != nil {
		return nil, err
	}

	return user, nil
}
