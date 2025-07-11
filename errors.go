package robloxgo

import "errors"

var (
	ErrNoAPIKey = errors.New("no api key provided")

	ErrNoUserID        = errors.New("no user id provided")
	ErrNoUsername      = errors.New("no username provided")
	ErrInvalidUsername = errors.New("invalid username provide")
	ErrUserHasNoRole   = errors.New("this user has no role")

	ErrNoGroupID        = errors.New("no group id provided")
	ErrNoGroupname      = errors.New("no group name provided")
	ErrInvalidGroupname = errors.New("invalid group name provided")

	ErrNoRoleID = errors.New("no role id provided")
)
