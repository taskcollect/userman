package handlers

import (
	"errors"
	"io/ioutil"
	"log"
	"main/dbutil"
	"main/util"
	"net/http"
	"time"

	"github.com/buger/jsonparser"
)

func (h *BaseHandler) NewUser(w http.ResponseWriter, r *http.Request) {
	if !EnsureMethod("POST", w, r) {
		return
	}

	err := h.db.Ping()
	if err != nil {
		// the database won't respond
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// construct a real user object
	user := dbutil.User{
		Username:     "",
		Secret:       "",
		Credentials:  []byte{'{', '}'},
		Preferences:  []byte{'{', '}'},
		RegisteredOn: time.Now(),
		LastLogin:    time.Now(),
	}

	// jsonparser.EachKey will iterate over all keys in the json object and pass values
	// to the switch statement.
	paths := [][]string{
		{"username"},
		{"secret"},
		{"creds"},
		{"prefs"},
	}

	// convert response body to byte slice
	data, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err = nil

	jsonparser.EachKey(data, func(idx int, value []byte, vt jsonparser.ValueType, _ error) {
		if err != nil {
			return
		}

		switch idx {
		case 0: // []string{"username"}
			if vt != jsonparser.String {
				err = errors.New("username must be string")
				return
			}
			user.Username = string(value)
		case 1: // []string{"secret"}
			if vt != jsonparser.String {
				err = errors.New("secret must be string")
				return
			}
			user.Secret = string(value)
		case 2: // []string{"creds"},
			if vt != jsonparser.Object {
				err = errors.New("creds must be object")
				return
			}
			user.Credentials = value // just take the json directly
		case 3: // []string{"prefs"},
			if vt != jsonparser.Object {
				err = errors.New("creds must be object")
				return
			}
			user.Preferences = value // just take the json directly
		}
	}, paths...)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("invalid structure"))
		return
	}

	if user.Username == "" || user.Secret == "" {
		// the request body was missing a field
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("missing username and/or secret"))
		return
	}

	exists, err := dbutil.DoesUserExist(h.db, user.Username)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println(err.Error())
		return
	}

	if exists {
		w.WriteHeader(http.StatusConflict)
		w.Write([]byte("already registered"))
		return
	}

	prefsOverrides, err := util.RemoveDefaultKeys(user.Preferences, dbutil.DefaultPreferences, true, true)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("invalid preferences"))
		return
	}

	user.Preferences = prefsOverrides

	// insert the user into the database
	if err := dbutil.InsertUser(h.db, &user); err != nil {
		// the database returned an error
		w.WriteHeader(http.StatusInternalServerError)
		log.Println(err.Error())
		return
	}
}
