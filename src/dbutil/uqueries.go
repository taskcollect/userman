package dbutil

import (
	"database/sql"
	"encoding/json"
	"errors"
	"main/security"
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

func GetUser(
	db *sql.DB,
	uname string,
	password string,
	wantCreds bool,
	wantPrefs bool,
) (*UserResponseSchema, error) {
	row := db.QueryRow("SELECT password, preferences, credentials FROM users WHERE username = $1", uname)

	if row.Err() != nil {
		return nil, row.Err()
	}

	var storedPassword string
	var credsJSON []byte
	var prefsJSON []byte

	err := row.Scan(&storedPassword, &prefsJSON, &credsJSON)
	if err != nil {
		return nil, err
	}

	if !security.VerifySecrets(password, storedPassword) {
		return nil, errors.New("unauthorized")
	}

	// decode the credentials and preferences
	var creds map[string]interface{}
	if wantCreds {
		// decode the credentials json
		if err := json.Unmarshal(credsJSON, &creds); err != nil {
			return nil, err
		}
	}

	var prefs map[string]interface{}
	if wantPrefs {
		// decode the preferences json
		if err := json.Unmarshal(prefsJSON, &prefs); err != nil {
			return nil, err
		}
	}

	return &UserResponseSchema{
		Credentials: creds,
		Preferences: prefs,
	}, nil
}
