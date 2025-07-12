package robloxgo

import (
	"encoding/json"
)

// User represents a Roblox user and its associated metadata.
type User struct {
	// ID is the unique identifier of the user.
	ID json.Number `json:"id"`

	// Username is the user's Roblox account name.
	Username string `json:"name"`

	// Displayname is the user's chosen display name, if set.
	Displayname string `json:"displayName"`

	// Premium indicates whether the user has Roblox Premium.
	Premium bool `json:"premium"`

	// Locale is the user's preferred language or region setting.
	Locale string `json:"locale"`

	// CreatedAt is the ISO 8601 timestamp of when the account was created.
	CreatedAt string `json:"createTime"`

	// Client is the API client used to interact with the user.
	Client *Client
}

// newUser returns a new User instance associated with the provided Client.
//
// It is intended for internal use to ensure that each User is linked to a Client,
// enabling the User's methods to perform API calls.
func newUser(client *Client) *User {
	return &User{
		Client: client,
	}
}

// GetUserByID retrieves a Roblox user from the Open Cloud API using the provided user ID.
//
// It returns a User instance associated with the current Client.
// An error is returned if the user ID is empty, if the HTTP request fails,
// if the response cannot be decoded, or if the user does not exist.
func (c *Client) GetUserByID(userID string) (*User, error) {
	if userID == "" {
		return nil, ErrNoUserID
	}

	resp, err := c.get(EndPointCloudUsers+userID, nil, nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	user := newUser(c)
	err = json.NewDecoder(resp.Body).Decode(user)
	if err != nil {
		return nil, err
	}

	return user, nil
}

// GetUserByUsername retrieves a Roblox user by their username using the legacy Roblox API.
//
// It returns a User instance associated with the current Client.
// An error is returned if the username is empty, the HTTP request fails,
// the response cannot be decoded, or if the user does not exist.
//
// Note: This method depends on the legacy endpoint at https://users.roblox.com/v1/usernames/users,
// which may be deprecated or removed by Roblox in the future.
func (c *Client) GetUserByUsername(username string) (*User, error) {
	if username == "" {
		return nil, ErrNoUsername
	}

	requestBody := map[string]interface{}{"usernames": []string{username}, "excludeBannedUsers": true}
	resp, err := c.post(EndpointLegacyGetUsers, requestBody, nil, nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var Response struct {
		Data []User `json:"data"`
	}
	err = json.NewDecoder(resp.Body).Decode(&Response)
	if err != nil {
		return nil, err
	}

	if len(Response.Data) == 0 {
		return nil, ErrInvalidUsername
	}

	legacyUser := &Response.Data[0]
	resp, err = c.get(EndPointCloudUsers+legacyUser.ID.String(), nil, nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	user := newUser(c)
	err = json.NewDecoder(resp.Body).Decode(user)
	if err != nil {
		return nil, err
	}

	return user, nil
}

// GetUserThumbnailURI retrieves the user's thumbnail image URI using the Open Cloud API.
//
// The request can be customized using optional query parameters such as format, size,
// and circular cropping. Returns the thumbnail URI as a string.
//
// Returns an error if the HTTP request fails or if the response body cannot be decoded.
func (u *User) GetUserThumbnailURI(queryParams []queryParam) (string, error) {
	methodURL := EndPointCloudUsers + u.ID.String() + ":generateThumbnail"
	resp, err := u.Client.get(methodURL, nil, queryParams)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	var thumbnailResponse struct {
		Response struct {
			Type     string `json:"@type"`
			ImageURI string `json:"imageUri"`
		} `json:"response"`
	}
	err = json.NewDecoder(resp.Body).Decode(&thumbnailResponse)
	if err != nil {
		return "", err
	}

	return thumbnailResponse.Response.ImageURI, nil
}
