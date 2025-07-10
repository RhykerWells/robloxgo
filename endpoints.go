// Robloxgo - Roblox bindings for Go
// Available at https://github.com/RhykerWells/robloxgo
//
// Copyright 2025 Rhyker Wells <a.rhykerw@gmail.com>.  All rights reserved.
// License can be found in the LICENSE file of the repository.
//
// Package robloxgo provides Roblox binding for Go
package robloxgo

// CloudAPIVersion is the Opencloud API version used for REST and Websocket API.
var (
	CloudAPIVersion = "2"
)

// Roblox API Endpoints
// These are used internally to build request paths for the
// various functions that this package provides
var (
	EndpointRoblox = "https://roblox.com"

	// Cloud APIs
	EndpointCloud      = "https://apis.roblox.com/cloud/v"
	EndpointCloudAPI   = EndpointCloud + CloudAPIVersion + "/"
	EndPointCloudUsers = EndpointCloudAPI + "users/"
	EndpointCloudGroups = EndpointCloudAPI + "groups/"

	// Legacy APIs
	EndpointLegacyUsers    = "https://users.roblox.com"
	EndpointLegacyGetUsers = EndpointLegacyUsers + "/v1/usernames/users"
	EndpointLegacyGroups = "https://groups.roblox.com"
	EndpointLegacyGetGroups = EndpointLegacyGroups + "/v1/groups/search/lookup"
)
