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
	"fmt"
	"strconv"
	"strings"
	"time"
)

// Group represents a Roblox group and its associated metadata.
type Group struct {
	// ID is the unique identifier of the group.
	ID json.Number `json:"id"`

	// Groupname is the display name of the group.
	Groupname string `json:"displayName"`

	// Description is the group's public description text.
	Description string `json:"description"`

	// OwnerID is the user ID of the group owner.
	OwnerID string `json:"owner"`

	// MemberCount is the number of users currently in the group.
	MemberCount json.Number `json:"memberCount"`

	// PublicEntry indicates whether users can join the group freely (true) or require approval (false).
	PublicEntry bool `json:"publicEntryAllowed"`

	// Locked reports whether the group is locked from further modifications.
	Locked bool `json:"locked"`

	// CreatedAt is the ISO 8601 timestamp of when the group was created.
	CreatedAt string `json:"createTime"`

	// Client is the API client used to interact with the group.
	Client *Client
}

// JoinRequest represents a user's request to join a Roblox group.
type JoinRequest struct {
	// ID is the unique identifier of the user making the request.
	ID string

	// Username is the Roblox username of the user.
	Username string

	// CreatedAt is the timestamp of when the join request was submitted.
	CreatedAt string
}

// GroupMember represents a user who is currently a member of a Roblox group.
type GroupMember struct {
	// ID is the unique identifier of the group member.
	ID string

	// Username is the Roblox username of the member.
	Username string

	// GroupRole is the member's role within the group.
	GroupRole GroupRole
}

// GroupRole represents a role within a Roblox group.
type GroupRole struct {
	// ID is the unique identifier of the role.
	ID json.Number `json:"id"`

	// Name is the display name of the role.
	Name string `json:"displayName"`

	// Rank is the hierarchical rank of the role within the group.
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
	resp, err := c.get(EndpointCloudGroups+groupID, nil, nil)
	if err != nil {
		return nil, err
	}

	group := newGroup(c)
	err = json.NewDecoder(resp.Body).Decode(group)
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
		return nil, ErrNoGroupname
	}

	groupHeader := queryParam{
		Key:   "groupName",
		Value: groupname,
	}
	resp, err := c.get(EndpointLegacyGetGroups, nil, []queryParam{groupHeader})
	if err != nil {
		return nil, err
	}

	var legacyResponse struct {
		Data []struct {
			ID   json.Number `json:"id"`
			Name string      `json:"name"`
		} `json:"data"`
	}
	err = json.NewDecoder(resp.Body).Decode(&legacyResponse)
	if err != nil {
		return nil, err
	}

	if len(legacyResponse.Data) == 0 || legacyResponse.Data[0].Name != groupname {
		return nil, ErrInvalidGroupname
	}

	legacyGroup := &legacyResponse.Data[0]

	resp, err = c.get(EndpointCloudGroups+legacyGroup.ID.String(), nil, nil)
	if err != nil {
		return nil, err
	}

	group := newGroup(c)
	err = json.NewDecoder(resp.Body).Decode(group)
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

// GetMembers retrieves all users from the group using the OpenCloud v2 API.
//
// This is a very hacky implementation on retrieving all the users.
// Neither the Legacy nor v2 OpenCloud API provide a way of retrieving
// just the user IDs from current group members.
//
// The OpenCloud endpoint imposes a 100 member max + 300 reqs/minute on member retrievals,
// because of this, when this function is called, we poll it every 200 milliseconds.
//
// I personally reccomend having keeping a local state of group users and update it every day due
// to how long the process of polling these users might take depending on group size.
//
// I am open to other suggestions of refactoring this if the limits are modified,
// or other methods of retrieving just the IDs are brought into the API.
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

			role, _ := g.GetUserRole(userID)

			members = append(members, GroupMember{
				ID:        userID,
				Username:  user.Username,
				GroupRole: *role,
			})
		}

		if membershipResponse.NextPage == "" {
			break
		}
		pageToken = membershipResponse.NextPage
	}

	return members, nil
}

// GetRoles retrieves all (legacy) roles from the group using the OpenCloud v2 API.
// Realistically this should retrieve the new
//
// This is a very hacky implementation on retrieving all the roles.
// Neither the Legacy nor v2 OpenCloud API provide a way of retrieving
// just the role IDs from current group roles.
//
// The OpenCloud endpoint imposes a 100 member max + 300 reqs/minute on role retrievals,
// because of this, when this function is called, we poll it every 200 milliseconds.
//
// I personally reccomend having keeping a local state of group roles and update it every day due
// to how long the process of polling these roles might take depending on number of roles.
//
// I am open to other suggestions of refactoring this if the limits are modified,
// or other methods of retrieving just the IDs are brought into the API.
//
// Please note that for larger groups it will take significantly longer to return
// the full role slice.
//
// TODO: Implement a Client/Session state and repoll this at set intervals instead?
func (g *Group) GetRoles() (roles []GroupRole, err error) {
	methodURL := EndpointCloudGroups + g.ID.String() + "/roles"
	var pageToken string

	rateLimit := time.NewTicker(200 * time.Millisecond)
	defer rateLimit.Stop()
	for {
		<-rateLimit.C

		query := []queryParam{{Key: "maxPageSize", Value: "20"}}
		if pageToken != "" {
			query = append(query, queryParam{Key: "pageToken", Value: pageToken})
		}

		resp, err := g.Client.get(methodURL, nil, query)
		if err != nil {
			return nil, err
		}

		var rolesResponse struct {
			NextPage   string      `json:"nextPageToken"`
			GroupRoles []GroupRole `json:"groupRoles"`
		}

		err = json.NewDecoder(resp.Body).Decode(&rolesResponse)
		resp.Body.Close()
		if err != nil {
			return nil, err
		}

		for _, role := range rolesResponse.GroupRoles {
			groupRole, err := g.GetRole(role.ID.String())
			if err != nil {
				continue
			}

			roles = append(roles, *groupRole)
		}

		if rolesResponse.NextPage == "" {
			break
		}
		pageToken = rolesResponse.NextPage
	}

	return roles, nil
}

// GetRole retrieves a group role from the Open Cloud API.
//
// Returns an error if the HTTP request fails, or if the response body cannot
// be decoded.
func (g *Group) GetRole(roleID string) (role *GroupRole, err error) {
	if roleID == "" {
		return nil, ErrNoRoleID
	}

	methodURL := EndpointCloudGroups + g.ID.String() + "/roles/" + roleID
	resp, err := g.Client.get(methodURL, nil, nil)
	if err != nil {
		return nil, err
	}

	err = json.NewDecoder(resp.Body).Decode(&role)
	if err != nil {
		return nil, err
	}

	return role, nil
}

// GetUserRole retrieves the users group role from both the Legacy Roblox API & Open Cloud API.
//
// Returns an error if the HTTP request fails, or if the response body cannot
// be decoded.
func (g *Group) GetUserRole(userID string) (*GroupRole, error) {
	user, err := g.Client.GetUserByID(userID)
	if err != nil {
		return nil, err
	}

	methodURL := EndpointLegacyGroups + "/v2/users/" + user.ID.String() + "/groups/roles"
	resp, err := g.Client.get(methodURL, nil, nil)
	if err != nil {
		return nil, err
	}

	var groupData struct {
		Data []struct {
			Group struct {
				ID json.Number `json:"id"`
			} `json:"group"`
			Role GroupRole `json:"role"`
		} `json:"data"`
	}
	err = json.NewDecoder(resp.Body).Decode(&groupData)
	if err != nil {
		return nil, err
	}

	for _, groupData := range groupData.Data {
		if groupData.Group.ID.String() != g.ID.String() {
			continue
		}
		groupRole, err := g.GetRole(groupData.Role.ID.String())
		if err != nil {
			return nil, err
		}
		return groupRole, nil
	}

	return nil, ErrUserHasNoRole
}

// UpdateUserRole updates a given users role in the group using the Open Cloud API
//
// Returns an error if the HTTP request fails, or if the response body cannot
// be decoded.
func (g *Group) UpdateUserRole(userID string, roleID string) (*GroupRole, error) {
	if userID == "" {
		return nil, ErrNoUserID
	}
	if roleID == "" {
		return nil, ErrNoRoleID
	}

	user, err := g.Client.GetUserByID(userID)
	if err != nil {
		return nil, err
	}

	role, err := g.GetRole(roleID)
	if err != nil {
		return nil, err
	}

	path := fmt.Sprintf("%s/memberships/%s", g.ID.String(), user.ID.String())
	requestBody := map[string]any{
		"path": "groups" + path,
		"user": "users/" + user.ID.String(),
		"role": "groups/" + g.ID.String() + "/roles/" + role.ID.String(),
	}
	_, err = g.Client.patch(EndpointCloudGroups+path, nil, requestBody)
	if err != nil {
		return nil, err
	}

	return role, nil
}

// RemoveUser removes a given user from the group using the legacy Roblox API
//
// Returns an error if the HTTP request fails, or if the response body cannot
// be decoded.
//
// This method may be deprecated if Roblox removes the
// legacy https://groups.roblox.com/v1/groups/{groupID}/users/{memberID} endpoint
func (g *Group) RemoveUser(userID string) (bool, error) {
	if userID == "" {
		return false, ErrNoUserID
	}

	ok, err := g.Client.delete(EndpointLegacyGroups+g.ID.String()+"/users/"+userID, nil)

	return ok, err
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
	resp, err := g.Client.get(EndpointLegacyGetGroupIcon, nil, querySet)
	if err != nil {
		return "", err
	}

	var thumbnailResponse struct {
		Data []struct {
			ImageURL string `json:"imageUrl"`
		} `json:"data"`
	}
	err = json.NewDecoder(resp.Body).Decode(&thumbnailResponse)
	if err != nil {
		return "", err
	}

	return thumbnailResponse.Data[0].ImageURL, nil
}
