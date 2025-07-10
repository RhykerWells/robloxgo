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
	"strconv"
	"strings"
	"time"
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

	// A slice of Roblox userIDs currently present in the Group.
	Members []string

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

// GetGroupByGroupname retrieves a Roblox group from the Legacy Roblox API by their group Groupname (case sensitive).
//
// Returns an error if the HTTP request fails, if the response body cannot
// be decoded, or if the group does not exist.
//
// This method may be deprecated if Roblox removes the
// legacy https://groups.roblox.com/v1/groups/search/lookup endpoint
func (c *Client) GetGroupByGroupname(groupname string) (*Group, error) {
	if groupname == "" {
		return nil, errors.New("no groupname")
	}

	groupHeader := queryParam{
		Key:   "groupName",
		Value: groupname,
	}
	response, err := c.get(EndpointLegacyGetGroups, nil, []queryParam{groupHeader})
	if err != nil {
		return nil, err
	}

	var legacyResponse struct {
		Data []struct {
			ID   json.Number `json:"id"`
			Name string      `json:"name"`
		} `json:"data"`
	}
	err = json.NewDecoder(response.Body).Decode(&legacyResponse)
	if err != nil {
		return nil, err
	}

	if len(legacyResponse.Data) == 0 || legacyResponse.Data[0].Name != groupname {
		return nil, errors.New("invalid groupname provided")
	}

	legacyGroup := &legacyResponse.Data[0]

	response, err = c.get(EndpointCloudGroups+legacyGroup.ID.String(), nil, nil)
	if err != nil {
		return nil, err
	}

	group := newGroup(c)
	err = json.NewDecoder(response.Body).Decode(group)
	if err != nil {
		return nil, err
	}
	group.Groupname = legacyGroup.Name

	return group, nil
}

// GetGroupIcon retrieves a Roblox group's thumbnail URl from the legacy Roblox API.
//
// Returns an error if the HTTP request fails, or if the response body cannot
// be decoded.
//
// This method may be deprecated if Roblox removes the
// legacy https://thumbnails.roblox.com/v1/groups/icons endpoint
func (g *Group) GetGroupIcon(large bool, isCircular bool) (string, error) {
	size := "150x150"
	if large {
		size = "420x420"
	}
	querySet := []queryParam{
		{
			Key:   "groupIds",
			Value: g.ID.String(),
		},
		{
			Key:   "format",
			Value: "Png",
		},
		{
			Key:   "size",
			Value: size,
		},
		{
			Key:   "isCircular",
			Value: strconv.FormatBool(isCircular),
		},
	}
	response, err := g.Client.get(EndpointLegacyGetGroupIcon, nil, querySet)
	if err != nil {
		return "", err
	}

	var thumbnailResponse struct {
		Data []struct {
			ImageURL string `json:"imageUrl"`
		} `json:"data"`
	}
	err = json.NewDecoder(response.Body).Decode(&thumbnailResponse)
	if err != nil {
		return "", err
	}

	return thumbnailResponse.Data[0].ImageURL, nil
}

// This is a very hacky implementation on retrieving all the user IDs.
// Neither the Legacy nor v2 OpenCloud API provide a way of retrieving
// just the user IDs from current group members.
//
// The OpenCloud endpoint imposes a 100 member max + 300 reqs/minute on member retrievals,
// because of this, we when this function is called, we poll it every 200 milliseconds.
//
// I personally reccomend having keeping a local state of group users and update it every day due
// to how long the process of polling these users might take depending on group size.
//
// I am open to other suggestions of refactoring this if the limits are modified,
// of other methods of retrieving just the IDs are brought into the API.
//
// Please note that for larger groups it will take significantly longer to return
// the full member slice.
func (g *Group) GetMembers() (members []string, err error) {
	methodURL := EndpointCloudGroups+g.ID.String() + "/memberships"
	var pageToken string

	rateLimit := time.NewTicker(200 * time.Millisecond)
	defer rateLimit.Stop()
	for {
		<-rateLimit.C

		query := []queryParam{{Key: "maxPageSize", Value: "100"}}
        if pageToken != "" {
            query = append(query, queryParam{Key: "pageToken", Value: pageToken})
        }

		resp, err := g.Client.get(methodURL, nil, query)
        if err != nil {
            return nil, err
        }

		var membershipResponse struct {
			NextPage string `json:"nextPageToken"`
			GroupMemberShip []struct {
				User string `json:"user"`
			} `json:"groupMemberships"`
		}

		err = json.NewDecoder(resp.Body).Decode(&membershipResponse)
        resp.Body.Close()
        if err != nil {
            return nil, err
        }

		for _, member := range membershipResponse.GroupMemberShip {
			members = append(members, strings.TrimPrefix(member.User, "users/"))
		}

		if membershipResponse.NextPage == "" {
			break
		}
		pageToken = membershipResponse.NextPage
	}
	return members, nil
}
