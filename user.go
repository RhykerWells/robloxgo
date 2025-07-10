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

	// The client used to connect to Roblox
	Client *Client
}


// userUser creates a new User instance associated with the given Client.
//
// This function is intended for internal use to ensure that every user
// has a reference to the Client, enable methods on the User object to make api calls.
func newUser(client *Client) *User {
	return &User{
		Client: client,
	}
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

	user := newUser(c)
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

	user := newUser(c)
	err = json.NewDecoder(response.Body).Decode(user)
	if err != nil {
		return nil, err
	}

	return user, nil
}

// ThumbnailURI retrieves a Roblox user's thumbnail URI from the Open Cloud API.
//
// Returns an error if the HTTP request fails, or if the response body cannot
// be decoded.
func (u *User) GetUserThumbnailURI(queryParams ...queryParam) (string, error) {
	methodURL := EndPointCloudUsers + u.ID.String() + ":generateThumbnail"
	response, err := u.Client.get(methodURL, queryParams...)
	if err != nil {
		return "", err
	}

	var thumbnailResponse struct {
		Response struct {
			Type     string `json:"@type"`
			ImageURI string `json:"imageUri"`
		} `json:"response"`
	}
	err = json.NewDecoder(response.Body).Decode(&thumbnailResponse)
	if err != nil {
		return "", err
	}

	return thumbnailResponse.Response.ImageURI, nil
}
