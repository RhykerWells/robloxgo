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

// get is an internal method that sends an internal HTTP GET request to the specified URL.
//
// It returns the HTTP response if the status code is 200 (OK).
// If the response status code is not 200, it returns an error
// containing the status and response body.
//
// The functions caller must close the returned response body.
func (c *Client) get(methodURL string, headers []httpHeader, parameters []queryParam) (*http.Response, error) {
	req, err := newHttpRequest(http.MethodGet, methodURL, nil, headers, parameters)
	if err != nil {
		return nil, err
	}

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}

	if err := httpErrorCheck(resp); err != nil {
		return nil, err
	}

	return resp, nil
}

// post is an internal method that sends a HTTP POST request to the specified URL with optional headers and a request body.
//
// It returns the HTTP response if the status code is 200 (OK).
// If the status code is not 200, it returns an error containing the status and response body.
//
// The caller must close the response body.
func (c *Client) post(methodURL string, body interface{}, headers []httpHeader, parameters []queryParam) (*http.Response, error) {
	req, err := newHttpRequest(http.MethodPost, methodURL, body, headers, parameters)
	if err != nil {
		return nil, err
	}

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}

	if err := httpErrorCheck(resp); err != nil {
		return nil, err
	}

	return resp, nil
}

// patch is an internal method that sends a HTTP PATCH request to the specified URL with optional headers.
//
// It returns the HTTP response if the status code is 200 (OK).
// If the status code is not 200, it returns an error containing the status and response body.
func (c *Client) patch(methodURL string, headers []httpHeader, body interface{}) (bool, error) {
	req, err := newHttpRequest(http.MethodPatch, methodURL, body, headers, nil)
	if err != nil {
		return false, err
	}

	resp, err := c.client.Do(req)
	if err != nil {
		return false, err
	}
	defer resp.Body.Close()

	if err := httpErrorCheck(resp); err != nil {
		return false, err
	}

	return true, nil
}

// delete is an internal method that sends a HTTP DELETE request to the specified URL with optional headers.
//
// It returns the HTTP response if the status code is 200 (OK).
// If the status code is not 200, it returns an error containing the status and response body.
func (c *Client) delete(methodURL string, headers []httpHeader) (bool, error) {
	req, err := newHttpRequest(http.MethodDelete, methodURL, nil, headers, nil)
	if err != nil {
		return false, err
	}

	resp, err := c.client.Do(req)
	if err != nil {
		return false, err
	}
	defer resp.Body.Close()

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

// httpErrorCheck validates the HTTP response status code.
//
// If the response status is not 200 OK, it reads and preserves the response body,
// then returns a formatted error including the status and response body contents.
//
// The response body is restored using io.NopCloser so it can still be read after the check.
// If the body cannot be read, a fallback error message is returned instead.
func httpErrorCheck(resp *http.Response) error {
	if resp.StatusCode == http.StatusOK {
		return nil
	}

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("http error %s: unable to read body: %v", resp.Status, err)
	}

	resp.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))
	return fmt.Errorf("http error %s: %s", resp.Status, string(bodyBytes))
}

// newHttpRequest constructs an HTTP request with optional query parameters, headers, and a JSON body.
//
// It accepts the HTTP method, request URL, optional request body (which is JSON-encoded if not nil),
// custom headers, and query parameters.
//
// If a body is provided and the method is not GET, the Content-Type is set to application/json.
// A User-Agent header is also added, including the library and Go runtime version.
//
// Returns the constructed *http.Request or an error if the request cannot be created.
func newHttpRequest(method string, methodURL string, body interface{}, headers []httpHeader, parameters []queryParam) (*http.Request, error) {
	parsedURL, _ := url.Parse(methodURL)

	q := parsedURL.Query()
	for _, parameter := range parameters {
		q.Set(parameter.Key, parameter.Value)
	}
	parsedURL.RawQuery = q.Encode()

	var requestBody bytes.Buffer
	if body != nil {
		json.NewEncoder(&requestBody).Encode(body)
	}

	req, err := http.NewRequest(method, parsedURL.String(), &requestBody)
	if err != nil {
		return nil, err
	}

	if method != http.MethodGet && body != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	for _, header := range headers {
		req.Header.Set(header.Key, header.Value)
	}

	req.Header.Set("User-Agent", robloxGoUserAgent)

	return req, nil
}
