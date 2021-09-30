package handlers

import (
	"errors"
	"io/ioutil"
	"log"
	"main/dbutil"
	"net/http"

	"github.com/buger/jsonparser"
)

func (h *BaseHandler) AlterUser(w http.ResponseWriter, r *http.Request) {
	if !EnsureMethod("POST", w, r) {
		return
	}

	paths := [][]string{
		{"user"},
		{"secret"},
		{"creds"},
		{"prefs"},
	}

	// make a special struct to keep all of the variables
	user := struct {
		Username   string
		Secret     string
		DeltaCreds []byte
		DeltaPrefs []byte
	}{
		// just initialize these with empty json, so we don't have to check if they're empty bytearrays
		DeltaCreds: []byte{'{', '}'},
		DeltaPrefs: []byte{'{', '}'},
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
		case 0: // []string{"user"}
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
			user.DeltaCreds = value // just take the json directly
		case 3: // []string{"prefs"},
			if vt != jsonparser.Object {
				err = errors.New("prefs must be object")
				return
			}
			user.DeltaPrefs = value // just take the json directly
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

	if !exists {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("user not found"))
		return
	}

	err = dbutil.AlterUser(h.db, user.Username, user.Secret, user.DeltaPrefs, user.DeltaCreds)

	if err != nil {
		switch err.Error() {
		case "unauthorized":
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte("no access to alter data"))
		case "invalid":
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("invalid alteration payload"))
		default:
			w.WriteHeader(http.StatusInternalServerError)
			log.Println(err.Error())
		}
	}
}
