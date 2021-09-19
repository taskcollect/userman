package handlers

import (
	"encoding/json"
	"log"
	"main/dbutil"
	"net/http"
	"time"
)

type userRegData struct {
	Username    string                 `json:"name"`
	Password    string                 `json:"secret"`
	Credentials dbutil.UserCredentials `json:"creds"`
	Preferences dbutil.UserPreferences `json:"prefs"`
}

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

	var udata userRegData

	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()

	// construct a new decoder and decode the request body into the struct
	if err := dec.Decode(&udata); err != nil {
		// the decoder returned an error
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}

	// construct a real user object
	user := dbutil.User{
		Username:     udata.Username,
		Password:     udata.Password,
		Credentials:  udata.Credentials,
		Preferences:  udata.Preferences,
		RegisteredOn: time.Now(),
		LastLogin:    time.Now(),
	}

	exists, err := dbutil.DoesUserExist(h.db, user.Username)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if exists {
		w.WriteHeader(http.StatusConflict)
		w.Write([]byte("already registered"))
		return
	}

	// insert the user into the database
	if err := dbutil.InsertUser(h.db, &user); err != nil {
		// the database returned an error
		w.WriteHeader(http.StatusInternalServerError)
		log.Println("Error while trying to register new user '" + user.Username + "': " + err.Error())
		return
	}
}
