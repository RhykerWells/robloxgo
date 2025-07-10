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
	"fmt"
	"io"
	"net/http"
	"net/url"
)

// get is an internal method that  sends an internal HTTP GET request to the specified URL.
//
// It returns the HTTP response if the status code is 200 (OK).
// If the response status code is not 200, it returns an error
// containing the status and response body.
//
// The functions caller must close the returned response body.
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

	if err := httpErrorCheck(resp); err != nil {
		return nil, err
	}

	return resp, nil
}

// post is an internal method that sends an HTTP POST request to the specified URL with optional headers and a request body.
//
// It returns the HTTP response if the status code is 200 (OK).
// If the status code is not 200, it returns an error containing the status and response body.
//
// The caller must close the response body.
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

	if err := httpErrorCheck(resp); err != nil {
		return nil, err
	}

	return resp, nil
}

// patch is an internal method that sends an HTTP PATCH request to the specified URL with optional headers.
//
// It returns the HTTP response if the status code is 200 (OK).
// If the status code is not 200, it returns an error containing the status and response body.
//
// The caller must close the response body.
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

	if err := httpErrorCheck(resp); err != nil {
		return false, err
	}

	return true, nil
}

// delete is an internal method that sends an HTTP DELETE request to the specified URL with optional headers.
//
// It returns the HTTP response if the status code is 200 (OK).
// If the status code is not 200, it returns an error containing the status and response body.
//
// The caller must close the response body.
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

	if resp.StatusCode != http.StatusOK {
		return false, fmt.Errorf("http error %d: %s", resp.StatusCode, resp.Status)
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

func httpErrorCheck(resp *http.Response) error {
	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("http error %s: unable to read body: %v", resp.Status, err)
	}

	resp.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("http error %s: %v", resp.Status, resp.Body)
	}
	return nil
}
