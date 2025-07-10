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

// JoinRequest stores all data for a pending user request to join a Roblox group.
type JoinRequest struct {
	// The ID of the user.
	ID        string

	// The user's username.
	Username  string

	// The user's join request date.
	CreatedAt string
}

// GroupMember stores all data for a user who is currently a member of a Roblox group.
type GroupMember struct {
	// The ID of the member.
	ID        string

	// The member's username.
	Username  string

	// The member's legacy group role.
	LegacyGroupRole LegacyGroupRole
}

// LegacyGroupRole stores all data for the legacy role type within a Roblox group,
type LegacyGroupRole struct {
	// The ID of the role.
	ID   json.Number `json:"id"`

	// The role's name.
	Name string      `json:"name"`

	// The role's heirarchial rank.
	Rank json.Number `json:"rank"`
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


// GetJoinRequests retrieves all pending join requests for the group from the Open Cloud AP.
//
// It returns these requests in a slice of JoinRequest structs.
//
// Each join request includes the user ID, username, and the timestamp the request was created.
//
// Returns an error if the HTTP request fails, or if the response body cannot
// be decoded.
func (g *Group) GetJoinRequests() (requests []JoinRequest, err error) {
	methodURL := EndpointCloudGroups + g.ID.String() + "/join-requests"
	resp, err := g.Client.get(methodURL, nil, nil)
	if err != nil {
		return requests, err
	}

	var requestData struct {
		GroupJoinRequests []struct {
			User      string `json:"user"`
			CreatedAt string `json:"createTime"`
		} `json:"groupJoinRequests"`
	}
	err = json.NewDecoder(resp.Body).Decode(&requestData)
	if err != nil {
		return requests, err
	}

	for _, request := range requestData.GroupJoinRequests {
		userID := strings.TrimPrefix(request.User, "users/")
		user, err := g.Client.GetUserByID(userID)
		if err != nil {
			continue
		}

		requests = append(requests, JoinRequest{
			ID:        userID,
			Username:  user.Username,
			CreatedAt: request.CreatedAt,
		})
	}

	return requests, err
}

// JoinRequestAccept accepts a pending group join request for the specified user ID.
//
// Returns an error if the HTTP request fails, or if the response body cannot
// be decoded.
//
// Returns true if the join request was successfully accepted.
func (g *Group) JoinRequestAccept(userID string) (bool, error) {
	_, err := g.Client.GetUserByID(userID)
	if err != nil {
		return false, err
	}

	methodURL := EndpointCloudGroups + g.ID.String() + "/join-requests" + userID + ":accept"
	requestBody := map[string]interface{}{}
	_, err = g.Client.post(methodURL, requestBody, nil, nil)
	if err != nil {
		return false, err
	}

	return true, nil
}

// JoinRequestAccept declines a pending group join request for the specified user ID.
//
// Returns an error if the HTTP request fails, or if the response body cannot
// be decoded.
//
// Returns true if the join request was successfully declined.
func (g *Group) JoinRequestDecline(userID string) (bool, error) {
	_, err := g.Client.GetUserByID(userID)
	if err != nil {
		return false, err
	}

	methodURL := EndpointCloudGroups + g.ID.String() + "/join-requests" + userID + ":decline"
	requestBody := map[string]interface{}{}
	_, err = g.Client.post(methodURL, requestBody, nil, nil)
	if err != nil {
		return false, err
	}

	return true, nil
}

// GetMembers retrieves all user IDs from the group using the OpenCloud v2 API.
//
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
//
// TODO: Implement a Client/Session state and repoll this at set intervals instead?
func (g *Group) GetMembers() (members []GroupMember, err error) {
	methodURL := EndpointCloudGroups + g.ID.String() + "/memberships"
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
			NextPage        string `json:"nextPageToken"`
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
			userID := strings.TrimPrefix(member.User, "users/")

			user, err := g.Client.GetUserByID(userID)
			if err != nil {
				continue
			}

			role, _ := g.GetUsersLegacyRole(userID)

			members = append(members, GroupMember{
				ID:              userID,
				Username:        user.Username,
				LegacyGroupRole: *role,
			})
		}

		if membershipResponse.NextPage == "" {
			break
		}
		pageToken = membershipResponse.NextPage
	}
	return members, nil
}

// GetUsersLegacyRole retrieves the users (legacy) group role from the Legacy Roblox AP.
//
// Returns an error if the HTTP request fails, or if the response body cannot
// be decoded.
func (g *Group) GetUsersLegacyRole(userID string) (*LegacyGroupRole, error) {
	user, err := g.Client.GetUserByID(userID)
	if err != nil {
		return nil, err
	}

	methodURL := EndpointLegacyGroups + "/v1/users/" + user.ID.String() + "/groups/roles"
	resp, err := g.Client.get(methodURL, nil, nil)
	if err != nil {
		return nil, err
	}

	var groupData struct {
		Data []struct {
			Group struct {
				ID json.Number `json:"id"`
			} `json:"group"`
			Role LegacyGroupRole `json:"role"`
		} `json:"data"`
	}
	err = json.NewDecoder(resp.Body).Decode(&groupData)
	if err != nil {
		return nil, err
	}

	for _, group := range groupData.Data {
		if group.Group.ID.String() != g.ID.String() {
			continue
		}
		return &group.Role, nil
	}

	return nil, errors.New("user has no role")
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
