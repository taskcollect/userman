package handlers

import (
	"encoding/json"
	"log"
	"main/dbutil"
	"net/http"
	"time"
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

	var ureg dbutil.UserRegSchema

	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()

	// construct a new decoder and decode the request body into the struct
	if err := dec.Decode(&ureg); err != nil {
		// the decoder returned an error
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("invalid json"))
		return
	}

	if ureg.Username == "" || ureg.Password == "" {
		// the request body was missing a field
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("missing username or password"))
		return
	}

	// construct a real user object
	user := dbutil.User{
		Username:     ureg.Username,
		Password:     ureg.Password,
		Credentials:  ureg.Credentials,
		Preferences:  ureg.Preferences,
		RegisteredOn: time.Now(),
		LastLogin:    time.Now(),
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

	// insert the user into the database
	if err := dbutil.InsertUser(h.db, &user); err != nil {
		// the database returned an error
		w.WriteHeader(http.StatusInternalServerError)
		log.Println(err.Error())
		return
	}
}
