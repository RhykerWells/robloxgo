// Robloxgo - Roblox bindings for Go
// Available at https://github.com/RhykerWells/robloxgo
//
// Copyright 2025 Rhyker Wells <a.rhykerw@gmail.com>.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//
// Package robloxgo provides Roblox binding for Go
package robloxgo

import (
	"errors"
	"fmt"
)

type httpResponse struct {
	Code  int
	Error string
}

// See https://create.roblox.com/docs/cloud/reference/errors
var (
	ResponseOK                 = httpResponse{Code: 200, Error: "response ok"}
	ResponseInvalid            = httpResponse{Code: 400, Error: "invalid argument passed"}
	ResponsePermissionDenied   = httpResponse{Code: 403, Error: "missing permission scopes"}
	ResponseResourceNotFound   = httpResponse{Code: 404, Error: "resource not found"}
	ResponseAborted            = httpResponse{Code: 409, Error: "operation aborted"}
	ResponseLimited            = httpResponse{Code: 429, Error: "too many requests"}
	ResponseRequestTerminated  = httpResponse{Code: 499, Error: "system terminated request"}
	ResponseInternalError      = httpResponse{Code: 500, Error: "the service replied with internal server error"}
	ResponseServiceUnavailable = httpResponse{Code: 503, Error: "the service is currently unavailable"}
)

var httpResponses = map[int]httpResponse{
	200: ResponseOK,
	400: ResponseInvalid,
	403: ResponsePermissionDenied,
	404: ResponseResourceNotFound,
	409: ResponseAborted,
	429: ResponseLimited,
	499: ResponseRequestTerminated,
	500: ResponseInternalError,
	503: ResponseServiceUnavailable,
}

// getFullHttpError returns a formatted error for the given response HTTP status code.
func getFullHttpError(errorCode int) error {
	httpResponse := httpResponses[errorCode]

	errorMessage := fmt.Sprintf("http error %d: %s", httpResponse.Code, httpResponse.Error)
	return errors.New(errorMessage)
}
