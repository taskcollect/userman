package dbutil

import (
	"database/sql"
	"encoding/json"
)

func InsertUser(db *sql.DB, user *User) error {
	q := `INSERT INTO users (
		username,
		password,
		registeredon,
		lastlogin,
		preferences,
		credentials
	) VALUES ( $1, $2, $3, $4, $5, $6  )`

	prefs, err := json.Marshal(user.Preferences)
	if err != nil {
		return err
	}

	creds, err := json.Marshal(user.Credentials)
	if err != nil {
		return err
	}

	_, err = RunQueryFailsafe(db, q,
		user.Username,
		user.Password,
		user.RegisteredOn.UTC(),
		user.LastLogin.UTC(),
		prefs,
		creds,
	)

	return err
}

func DoesUserExist(db *sql.DB, uname string) (bool, error) {
	var userExists bool
	err := db.QueryRow("SELECT EXISTS(SELECT 1 FROM users WHERE username = $1)", uname).Scan(&userExists)

	if err != nil {
		return false, err
	}

	return userExists, nil
}
