// Robloxgo - Roblox bindings for Go
// Available at https://github.com/RhykerWells/robloxgo

// Copyright 2025 Rhyker Wells <a.rhykerw@gmail.com>.  All rights reserved.
// License can be found in the LICENSE file of the repository.

// Package robloxgo provides Roblox binding for Go

package robloxgo

import (
	"encoding/json"
)

// A User stores all data for an individual Roblox user.
type User struct {
	// The ID of the user.
	ID string `json:"id"`

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