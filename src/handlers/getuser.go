package handlers

import (
	"encoding/json"
	"log"
	"main/dbutil"
	"net/http"
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

	var ureq dbutil.UserRequestSchema

	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()

	// construct a new decoder and decode the request body into the struct
	if err := dec.Decode(&ureq); err != nil {
		// the decoder returned an error
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("invalid json"))
		return
	}

	if ureq.Username == "" || ureq.Password == "" {
		// the request body was missing a field
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("missing username or password"))
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

	if !ureq.Fields.WantCredentials && !ureq.Fields.WantPreferences {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("no field requested"))
		return
	}

	uresp, err := dbutil.GetUser(
		h.db,
		ureq.Username,
		ureq.Password,
		ureq.Fields.WantCredentials,
		ureq.Fields.WantPreferences,
	)

	if err != nil {
		if err.Error() == "unauthorized" {
			w.WriteHeader(http.StatusUnauthorized)
		} else {
			w.WriteHeader(http.StatusInternalServerError)
			log.Println(err.Error())
		}
		return
	}

	enc := json.NewEncoder(w)
	err = enc.Encode(uresp)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println(err.Error())
		return
	}
}
