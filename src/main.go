package main

import (
	"database/sql"
	"log"
	"main/dbutil"
	"main/handlers"
	"net/http"
)

var defaultConfig = dbutil.SqlConfig{
	Host:     "database",
	Port:     5432,
	User:     "dev",
	Password: "dev",
	Database: "taskcollect",
}

func makeMux(db *sql.DB) *http.ServeMux {
	// server := &http.Server{
	// 	Addr: ":2000",
	// }

	mux := http.NewServeMux()

	handler := handlers.NewBaseHandler(db)

	mux.HandleFunc("/v1/register", handler.NewUser)
	mux.HandleFunc("/v1/get", handler.GetUser)
	mux.HandleFunc("/v1/alter", handler.AlterUser)

	return mux
}

func main() {
	log.Printf(
		"Trying to open connection to db '%s' as %s@%s:%d",
		defaultConfig.Database, defaultConfig.User, defaultConfig.Host, defaultConfig.Port,
	)

	db, err := dbutil.Open(&defaultConfig)

	if err != nil {
		log.Fatalf("Could not connect: %v\n", err)
	}

	log.Println("Connection OK, setting up database...")

	err = dbutil.Initialize(db)
	if err != nil {
		log.Fatalf("Could not run db init: %v\n", err)
	}

	log.Println("Database inititalized, starting server...")

	mux := makeMux(db)
	http.ListenAndServe(":2000", handlers.RequestLogger(mux))

	// remember to close the connection once we exit
	defer db.Close()
}
