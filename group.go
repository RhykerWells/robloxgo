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
)

// A Group stores all data for an individual Roblox group.
type Group struct {
	// The ID of the group.
	ID json.Number `json:"id"`

	// The group's name.
	Groupname string `json:"displayName"`

	// The group's description.
	Description string `json:"description"`

	// The group's owner ID.
	OwnerID string `json:"owner"`

	// The group's member count.
	MemberCount json.Number `json:"memberCount"`

	// The group's entry state.
	// Returns true if public, false if set to request only
	PublicEntry bool `json:"publicEntryAllowed"`

	// The groups locked state
	Locked bool `json:"locked"`

	// The group's account creation date.
	CreatedAt string `json:"createTime"`

	// The client used to connect to Roblox.
	Client *Client
}

// newGroup creates a new Group instance associated with the given Client.
//
// This function is intended for internal use to ensure that every group
// has a reference to the Client, and enable methods on the Group object to make api calls.
func newGroup(client *Client) *Group {
	return &Group{
		Client: client,
	}
}

// GetGroupByID retrieves a Roblox group from the Open Cloud API by their group ID.
//
// Returns an error if the HTTP request fails, if the response body cannot
// be decoded, or if the group does not exist.
func (c *Client) GetGroupByID(groupID string) (*Group, error) {
	response, err := c.get(EndpointCloudGroups+groupID, nil, nil)
	if err != nil {
		return nil, err
	}

	group := newGroup(c)
	err = json.NewDecoder(response.Body).Decode(group)
	if err != nil {
		return nil, err
	}

	return group, nil
}
