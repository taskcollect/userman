package dbutil

import (
	"time"
)

type UserCredentials struct {
	GoogleToken string `json:"google"`
}

type UserPreferences struct {
	TimeIn24H bool `json:"24h"`
}

type User struct {
	Username     string
	Password     string
	Credentials  UserCredentials
	Preferences  UserPreferences
	RegisteredOn time.Time
	LastLogin    time.Time
}

type UserRegSchema struct {
	Username    string          `json:"user"`
	Password    string          `json:"secret"`
	Credentials UserCredentials `json:"creds"`
	Preferences UserPreferences `json:"prefs"`
}

type UserFieldSelector struct {
	WantCredentials bool `json:"creds"`
	WantPreferences bool `json:"prefs"`
}

type UserRequestSchema struct {
	Username string            `json:"user"`
	Password string            `json:"secret"`
	Fields   UserFieldSelector `json:"fields,omitempty"`
}

type UserResponseSchema struct {
	Credentials map[string]interface{} `json:"creds,omitempty"`
	Preferences map[string]interface{} `json:"prefs,omitempty"`
}
