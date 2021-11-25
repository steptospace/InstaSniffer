// Package api provides primitives to interact with the openapi HTTP API.
//
// Code generated by github.com/deepmap/oapi-codegen version v1.9.0 DO NOT EDIT.
package api

import (
	"time"
)

// ErrInfo defines model for ErrInfo.
type ErrInfo struct {
	Description string `json:"description"`
	Err         string `json:"err"`
}

// ImportantInfo defines model for ImportantInfo.
type ImportantInfo struct {
	Avatar    string    `json:"avatar"`
	Bio       string    `json:"bio"`
	CreatedAt time.Time `json:"created_at"`
	Images    []Media   `json:"images"`
	Name      string    `json:"name"`
	Username  string    `json:"username"`
	Videos    []Media   `json:"videos"`
}

// Media defines model for Media.
type Media struct {
	Description string `json:"description"`
	Url         string `json:"url"`
}

// User defines model for User.
type User struct {
	Id   string `json:"id"`
	Name string `json:"name"`
}

// PostUsersJSONBody defines parameters for PostUsers.
type PostUsersJSONBody User

// PostUsersJSONRequestBody defines body for PostUsers for application/json ContentType.
type PostUsersJSONRequestBody PostUsersJSONBody