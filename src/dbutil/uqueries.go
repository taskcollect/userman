package dbutil

import (
	"database/sql"
	"errors"
	"main/security"
	"main/util"
	"strings"

	"github.com/buger/jsonparser"
)

func InsertNewUser(db *sql.DB, user *User) error {
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

// returns stored secret, stored creds, stored prefs overrides
func getUsefulData(db *sql.DB, uname, secret string) (string, []byte, []byte, error) {
	row := db.QueryRow("SELECT secret, preferences, credentials FROM users WHERE username = $1", uname)

	if row.Err() != nil {
		return "", nil, nil, row.Err()
	}

	// stuff that will be written into by the row scan
	var (
		storedSecret         string
		storedCreds          []byte
		storedPrefsOverrides []byte
	)
	// scan the data from the db into the above vars
	err := row.Scan(&storedSecret, &storedPrefsOverrides, &storedCreds)
	if err != nil {
		return "", nil, nil, err
	}

	return storedSecret, storedCreds, storedPrefsOverrides, nil
}

func GetUser(db *sql.DB, uname string, secret string, wantCreds bool, wantPrefs bool) ([]byte, error) {
	row := db.QueryRow("SELECT secret, preferences, credentials FROM users WHERE username = $1", uname)

	if row.Err() != nil {
		return nil, row.Err()
	}

	storedSecret, storedCreds, storedPrefsOverrides, err := getUsefulData(db, uname, secret)
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

func AlterUser(db *sql.DB, uname string, secret string, deltaPrefs, deltaCreds []byte) error {
	storedSecret, storedCreds, storedPrefsOverrides, err := getUsefulData(db, uname, secret)
	if err != nil {
		return err
	}

	// make sure the user is authorized to access the data
	if !security.VerifySecrets(secret, storedSecret) {
		return errors.New("unauthorized")
	}

	newPrefs, err := util.Merge(deltaPrefs, storedPrefsOverrides)
	if err != nil {
		return errors.New("invalid")
	}

	newPrefs, err = util.RemoveDefaultKeys(newPrefs, DefaultPreferences, true, true)
	if err != nil {
		if strings.HasPrefix(err.Error(), "invalid:") {
			return errors.New("invalid")
		}
		return err
	}

	newCreds, err := util.Merge(deltaCreds, storedCreds)
	if err != nil {
		return errors.New("invalid")
	}

	q := `UPDATE users SET preferences = $1, credentials = $2 WHERE username = $3`

	// update the db
	_, err = RunQueryFailsafe(db, q, newPrefs, newCreds, uname)
	return err
}
