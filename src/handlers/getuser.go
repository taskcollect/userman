package handlers

import (
	"errors"
	"io/ioutil"
	"log"
	"main/dbutil"
	"net/http"

	"github.com/buger/jsonparser"
)

func (h *BaseHandler) GetUser(w http.ResponseWriter, r *http.Request) {
	if !EnsureMethod("GET", w, r) {
		return
	}

	err := h.db.Ping()
	if err != nil {
		// the database won't respond
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	ureq := dbutil.UserRequestSchema{
		Username: "",
		Secret:   "",
		Fields: dbutil.UserFieldSelector{
			WantCredentials: false,
			WantPreferences: false,
		},
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
			ureq.Username = string(value)
		case 1: // []string{"secret"}
			if vt != jsonparser.String {
				err = errors.New("secret must be string")
				return
			}

			ureq.Secret = string(value)
		case 2: // []string{"creds"}
			bval, e := jsonparser.ParseBoolean(value)
			if e != nil {
				err = e
				return
			}
			ureq.Fields.WantCredentials = bval
		case 3: // []string{"prefs"}
			bval, e := jsonparser.ParseBoolean(value)
			if e != nil {
				err = e
				return
			}
			ureq.Fields.WantPreferences = bval
		}
	}, paths...)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("invalid structure"))
		return
	}

	if ureq.Username == "" || ureq.Secret == "" {
		// the request body was missing a field
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("missing username and/or secret"))
		return
	}

	exists, err := dbutil.DoesUserExist(h.db, ureq.Username)

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

	uresp, err := dbutil.GetUser(
		h.db,
		ureq.Username,
		ureq.Secret,
		ureq.Fields.WantCredentials,
		ureq.Fields.WantPreferences,
	)

	if err != nil {
		if err.Error() == "unauthorized" {
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte("no access to data"))
		} else {
			w.WriteHeader(http.StatusInternalServerError)
			log.Println(err.Error())
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_, err = w.Write(uresp)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println(err.Error())
		return
	}
}
