package dbutil

import "time"

type UserCredentials struct {
	GoogleToken string `json:"google,omitempty"`
}

type UserPreferences struct {
	TimeIn24H bool `json:"24h,omitempty"`
}

type User struct {
	Username     string
	Password     string
	Credentials  UserCredentials
	Preferences  UserPreferences
	RegisteredOn time.Time
	LastLogin    time.Time
}
