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
	CreatedAt time.Time
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

// newGroup returns a new Group instance associated with the provided Client.
//
// It is intended for internal use to ensure that each Group is linked to a Client,
// enabling the Group's methods to perform API calls.
func newGroup(client *Client) *Group {
	return &Group{
		Client: client,
	}
}

// GetGroupByID retrieves a Roblox group from the Open Cloud API using the provided group ID.
//
// It returns a Group instance associated with the current Client.
// An error is returned if the group ID is empty, if the HTTP request fails,
// if the response cannot be decoded, or if the group does not exist.
func (c *Client) GetGroupByID(groupID string) (*Group, error) {
	if groupID == "" {
		return nil, ErrNoGroupID
	}

	resp, err := c.get(EndpointCloudGroups+groupID, nil, nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	group := newGroup(c)
	err = json.NewDecoder(resp.Body).Decode(group)
	if err != nil {
		return nil, err
	}
	group.OwnerID = strings.TrimPrefix(group.OwnerID, "users/")

	return group, nil
}

// GetGroupByGroupname retrieves a Roblox group using the legacy Roblox API by its exact group name (case sensitive).
//
// It returns a Group instance associated with the current Client.
// An error is returned if the group name is empty, the HTTP request fails,
// the response cannot be decoded, or if the group cannot be found.
//
// Note: This method relies on the legacy endpoint at https://groups.roblox.com/v1/groups/search/lookup,
// which may be deprecated or removed by Roblox in the future.
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
	defer resp.Body.Close()

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
	defer resp.Body.Close()

	group := newGroup(c)
	err = json.NewDecoder(resp.Body).Decode(group)
	if err != nil {
		return nil, err
	}
	group.Groupname = legacyGroup.Name
	group.OwnerID = strings.TrimPrefix(group.OwnerID, "users/")

	return group, nil
}

// GetJoinRequests retrieves all pending join requests for the group using the Open Cloud API.
//
// It returns a slice of JoinRequest structs, each containing the user ID, username,
// and the timestamp the request was created.
//
// An error is returned if the HTTP request fails or if the response cannot be decoded.
// If an individual user lookup fails, that request is skipped and the remaining are still returned.
func (g *Group) GetJoinRequests() (requests []JoinRequest, err error) {
	methodURL := EndpointCloudGroups + g.ID.String() + "/join-requests"
	resp, err := g.Client.get(methodURL, nil, nil)
	if err != nil {
		return requests, err
	}
	defer resp.Body.Close()

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
		timestamp, _ := time.Parse(time.RFC3339, request.CreatedAt)
		requests = append(requests, JoinRequest{
			ID:        userID,
			Username:  user.Username,
			CreatedAt: timestamp.UTC(),
		})
	}

	return requests, err
}

// JoinRequestAccept approves a pending group join request for the specified user ID.
//
// Returns true if the request was successfully accepted.
// Returns an error if the user does not exist, the HTTP request fails,
// or the response cannot be decoded.
func (g *Group) JoinRequestAccept(userID string) (bool, error) {
	if userID == "" {
		return false, ErrNoUserID
	}

	_, err := g.Client.GetUserByID(userID)
	if err != nil {
		return false, err
	}

	methodURL := EndpointCloudGroups + g.ID.String() + "/join-requests/" + userID + ":accept"
	requestBody := map[string]interface{}{}
	resp, err := g.Client.post(methodURL, requestBody, nil, nil)
	if err != nil {
		return false, err
	}
	resp.Body.Close()

	return true, nil
}

// JoinRequestDecline rejects a pending group join request for the specified user ID.
//
// Returns true if the request was successfully declined.
// Returns an error if the user does not exist, the HTTP request fails,
// or the response cannot be decoded.
func (g *Group) JoinRequestDecline(userID string) (bool, error) {
	if userID == "" {
		return false, ErrNoUserID
	}

	_, err := g.Client.GetUserByID(userID)
	if err != nil {
		return false, err
	}

	methodURL := EndpointCloudGroups + g.ID.String() + "/join-requests/" + userID + ":decline"
	requestBody := map[string]interface{}{}
	resp, err := g.Client.post(methodURL, requestBody, nil, nil)
	if err != nil {
		return false, err
	}
	resp.Body.Close()

	return true, nil
}

// GetMembers retrieves all users in the group using the Open Cloud v2 API.
//
// Due to current limitations of both the legacy and Open Cloud APIs, there is no
// direct way to fetch only the user IDs of group members. This method works around that
// by paginating over the full member list (100 users per request) and polling the
// endpoint every 200 milliseconds to respect Roblox’s rate limit of 300 requests/minute.
//
// For large groups, this process can be slow. It is recommended to cache member data
// locally and update it periodically instead of calling this method frequently.
//
// Returns a slice of GroupMember structs. An error is returned if any request fails
// or a response cannot be decoded. Individual user lookups that fail are skipped.
//
// TODO: Consider caching state and repolling periodically in a background session.
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
		defer resp.Body.Close()

		var membershipResponse struct {
			NextPage        string `json:"nextPageToken"`
			GroupMembership []struct {
				User string `json:"user"`
			} `json:"groupMemberships"`
		}

		err = json.NewDecoder(resp.Body).Decode(&membershipResponse)
		if err != nil {
			return nil, err
		}

		for _, member := range membershipResponse.GroupMembership {
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

// GetRoles returns all roles defined within the group.
//
// Each role is retrieved and resolved into a complete GroupRole object.
// If a role lookup fails, it is skipped.
// Returns an error if the HTTP request fails or the response cannot be decoded.
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
		defer resp.Body.Close()

		var rolesResponse struct {
			NextPage   string      `json:"nextPageToken"`
			GroupRoles []GroupRole `json:"groupRoles"`
		}

		err = json.NewDecoder(resp.Body).Decode(&rolesResponse)
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

// GetRole retrieves a specific group role by its role ID.
//
// Returns a GroupRole associated with the given role ID.
// Returns an error if the role ID is empty, the HTTP request fails,
// or the response body cannot be decoded.
func (g *Group) GetRole(roleID string) (role *GroupRole, err error) {
	if roleID == "" {
		return nil, ErrNoRoleID
	}

	methodURL := EndpointCloudGroups + g.ID.String() + "/roles/" + roleID
	resp, err := g.Client.get(methodURL, nil, nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	err = json.NewDecoder(resp.Body).Decode(&role)
	if err != nil {
		return nil, err
	}

	return role, nil
}

// GetUserRole retrieves the role of a specific user within the group.
//
// This method uses both the legacy Roblox API and the Open Cloud API to determine
// the user's role. It returns a GroupRole pointer if the user has a role in the group.
// Returns an error if the user does not exist, if the user has no role in the group,
// if the HTTP request fails, or if the response body cannot be decoded.
func (g *Group) GetUserRole(userID string) (*GroupRole, error) {
	if userID == "" {
		return nil, ErrNoUserID
	}

	user, err := g.Client.GetUserByID(userID)
	if err != nil {
		return nil, err
	}

	methodURL := EndpointLegacyGroups + "/v2/users/" + user.ID.String() + "/groups/roles"
	resp, err := g.Client.get(methodURL, nil, nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

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

// UpdateUserRole sets a user's role in the group using the Open Cloud API.
//
// Returns the updated GroupRole if the operation is successful.
// Returns an error if the user ID or role ID is empty, the user or role cannot be found,
// the HTTP request fails, or the response cannot be decoded.
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
	requestBody := map[string]string{
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

// RemoveUser removes a user from the group using the legacy Roblox API.
//
// Returns true if the user was successfully removed.
// Returns an error if the user ID is empty, the HTTP request fails, or the response cannot be decoded.
//
// Note: This method uses the legacy endpoint at
// https://groups.roblox.com/v1/groups/{groupID}/users/{memberID}, which may be deprecated in the future.
func (g *Group) RemoveUser(userID string) (bool, error) {
	if userID == "" {
		return false, ErrNoUserID
	}

	ok, err := g.Client.delete(EndpointLegacyGroups+g.ID.String()+"/users/"+userID, nil)

	return ok, err
}

// GetGroupIcon retrieves the group's thumbnail image URL using the legacy Roblox API.
//
// The size of the icon can be set to either 150x150 or 420x420 based on the `large` flag.
// The `isCircular` flag determines whether the returned icon is circular.
//
// Returns the image URL as a string. An error is returned if the HTTP request fails
// or if the response cannot be decoded.
//
// Note: This method uses the legacy endpoint at
// https://thumbnails.roblox.com/v1/groups/icons, which may be deprecated in the future.
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
	defer resp.Body.Close()

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
