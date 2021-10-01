package handlers

import (
	"errors"
	"io/ioutil"
	"log"
	"main/dbutil"
	"net/http"

	"github.com/buger/jsonparser"
)

func (h *BaseHandler) DeleteUser(w http.ResponseWriter, r *http.Request) {
	if !EnsureMethod("DELETE", w, r) {
		return
	}

	// jsonparser.EachKey will iterate over all keys in the json object and pass values
	// to the switch statement.
	paths := [][]string{
		{"user"},
		{"secret"},
	}

	var (
		username string
		secret   string
	)

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
			username = string(value)
		case 1: // []string{"secret"}
			if vt != jsonparser.String {
				err = errors.New("secret must be string")
				return
			}
			secret = string(value)
		}
	}, paths...)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("invalid structure"))
		return
	}

	if username == "" || secret == "" {
		// the request body was missing a field
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("missing username and/or secret"))
		return
	}

	exists, err := dbutil.DoesUserExist(h.db, username)

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

	// delete the user
	err = dbutil.DeleteUser(h.db, username, secret)
	if err != nil {
		if err.Error() == "unauthorized" {
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte("no access to delete user"))
			return
		}

		// the database returned an error
		w.WriteHeader(http.StatusInternalServerError)
		log.Println(err.Error())
		return
	}
}
