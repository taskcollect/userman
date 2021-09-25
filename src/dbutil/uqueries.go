package dbutil

import (
	"database/sql"
	"errors"
	"main/security"
	"main/util"

	"github.com/buger/jsonparser"
)

func InsertUser(db *sql.DB, user *User) error {
	q := `INSERT INTO users (
		username,
		secret,
		registeredon,
		lastlogin,
		preferences,
		credentials
	) VALUES ( $1, $2, $3, $4, $5, $6  )`

	_, err := RunQueryFailsafe(db, q,
		user.Username,
		user.Secret,
		user.RegisteredOn.UTC(),
		user.LastLogin.UTC(),
		user.Preferences,
		user.Credentials,
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
	secret string,
	wantCreds bool,
	wantPrefs bool,
) ([]byte, error) {
	row := db.QueryRow("SELECT secret, preferences, credentials FROM users WHERE username = $1", uname)

	if row.Err() != nil {
		return nil, row.Err()
	}

	// stuff that will be written into by the row scan
	var storedSecret string
	var storedCreds []byte
	var storedPrefsOverrides []byte

	// scan the data from the db into the above vars
	err := row.Scan(&storedSecret, &storedPrefsOverrides, &storedCreds)
	if err != nil {
		return nil, err
	}

	// make sure the user is authorized to access the data
	if !security.VerifySecrets(secret, storedSecret) {
		return nil, errors.New("unauthorized")
	}

	// prepare the output json
	out := []byte{'{', '}'}

	// decode the credentials and preferences

	/*
		this might be a little bit insecure as the json does not get validated
		at any point in the code, from input to storage to output. however, given
		that this service would be behind others, it's not a big deal.
	*/

	if wantPrefs {
		// the db only stores overrides for the preferences, so we need to merge with the default
		fullPrefs, err := util.AddDefaultKeys(storedPrefsOverrides, DefaultPreferences, true)

		if err != nil {
			return nil, err
		}

		// add the preferences to the output in the prefs key
		out, err = jsonparser.Set(out, fullPrefs, "prefs")

		if err != nil {
			return nil, err
		}
	}

	if wantCreds {
		out, err = jsonparser.Set(out, storedCreds, "creds")
		if err != nil {
			return nil, err
		}
	}

	return out, nil
}
