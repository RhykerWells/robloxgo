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
	"net/url"
)

// get is an internal method that makes a GET request to the specified URL
//
// If the response status code is not 200 (OK/Successful), it
// returns a custom error describing the HTTP status code
func (c *Client) get(methodURL string, headers []httpHeader, queryParams []queryParam) (*http.Response, error) {
	parsedURL, _ := url.Parse(methodURL)
	q := parsedURL.Query()
	for _, param := range queryParams {
		q.Set(param.Key, param.Value)
	}
	parsedURL.RawQuery = q.Encode()

	req, err := http.NewRequest(http.MethodGet, parsedURL.String(), nil)
	if err != nil {
		return nil, err
	}
	for _, header := range headers {
		req.Header.Set(header.Key, header.Value)
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
func (c *Client) post(methodURL string, body any, headers []httpHeader, queryParams []queryParam) (*http.Response, error) {
	var requestBody bytes.Buffer
	json.NewEncoder(&requestBody).Encode(body)

	parsedURL, _ := url.Parse(methodURL)
	q := parsedURL.Query()
	for _, param := range queryParams {
		q.Set(param.Key, param.Value)
	}
	parsedURL.RawQuery = q.Encode()

	req, err := http.NewRequest(http.MethodPost, parsedURL.String(), &requestBody)
	if err != nil {
		return nil, err
	}
	for _, header := range headers {
		req.Header.Set(header.Key, header.Value)
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

func (c *Client) patch(methodURL string, headers []httpHeader, body any) (bool, error) {
	var requestBody bytes.Buffer
	json.NewEncoder(&requestBody).Encode(body)

	parsedURL, _ := url.Parse(methodURL)
	req, err := http.NewRequest(http.MethodPatch, parsedURL.String(), &requestBody)
	if err != nil {
		return false, err
	}
	for _, header := range headers {
		req.Header.Set(header.Key, header.Value)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.client.Do(req)
	if err != nil {
		return false, err
	}

	if resp.StatusCode != ResponseOK.Code {
		return false, getFullHttpError(resp.StatusCode)
	}

	return true, nil
}

func (c *Client) delete(methodURL string, headers []httpHeader) (bool, error) {
	parsedURL, _ := url.Parse(methodURL)
	req, err := http.NewRequest(http.MethodDelete, parsedURL.String(), nil)
	if err != nil {
		return false, err
	}
	for _, header := range headers {
		req.Header.Set(header.Key, header.Value)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.client.Do(req)
	if err != nil {
		return false, err
	}

	if resp.StatusCode != ResponseOK.Code {
		return false, getFullHttpError(resp.StatusCode)
	}

	return true, nil
}

type httpHeader struct {
	// The key (case sensitive) for the HTTP header
	Key string
	// The value for the HTTP header
	Value string
}

type queryParam struct {
	// The key for the query parameter
	Key string
	// The value for the query parameter
	Value string
}
