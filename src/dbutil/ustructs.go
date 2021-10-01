package dbutil

import (
	"time"
)

var DefaultPreferences = []byte(`{
	"time24h": false,
	"accentColor": "blue",
}`)

type User struct {
	Username     string
	Secret       string
	Credentials  []byte
	Preferences  []byte
	RegisteredOn time.Time
}

type UserFieldSelector struct {
	WantCredentials bool
	WantPreferences bool
}

type UserRequestSchema struct {
	Username string
	Secret   string
	Fields   UserFieldSelector
}
