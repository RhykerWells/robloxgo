// Robloxgo - Roblox bindings for Go
// Available at https://github.com/RhykerWells/robloxgo
//
// Copyright 2025 Rhyker Wells <a.rhykerw@gmail.com>.  All rights reserved.
// License can be found in the LICENSE file of the repository.
//
// Package robloxgo provides Roblox binding for Go
package robloxgo

import (
	"bytes"
	"encoding/json"
	"net/http"
)

// get is an internal method that makes a GET request to the specified URL
//
// If the response status code is not 200 (OK/Successful), it
// returns a custom error describing the HTTP status code
func (c *Client) get(methodURL string, queryParams ...queryParam) (*http.Response, error) {
	req, err := http.NewRequest(http.MethodGet, methodURL, nil)
	if err != nil {
		return nil, err
	}
	for _, queryParam := range queryParams {
		req.Header.Set(queryParam.Key, queryParam.Value)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != ResponseOK.Code {
		return nil, getFullHttpError(resp.StatusCode)
	}

	return resp, nil
}

// post is an internal method that makes a POST request to the specified URL
//
// If the response status code is not 200 (OK/Successful), it
// returns a custom error describing the HTTP status code
func (c *Client) post(methodURL string, body any) (*http.Response, error) {
	var requestBody bytes.Buffer
	json.NewEncoder(&requestBody).Encode(body)
	req, err := http.NewRequest(http.MethodPost, methodURL, &requestBody)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != ResponseOK.Code {
		return nil, getFullHttpError(resp.StatusCode)
	}

	return resp, nil
}

type queryParam struct {
	Key string
	Value string
}