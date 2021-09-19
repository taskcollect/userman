package main

import (
	"log"
	"main/dbutil"
)

var defaultConfig = dbutil.SqlConfig{
	Host:     "database",
	Port:     5432,
	User:     "dev",
	Password: "dev",
	Database: "taskcollect",
}

func main() {
	log.Printf(
		"trying to open connection to db '%s' as %s@%s:%d",
		defaultConfig.Database, defaultConfig.User, defaultConfig.Host, defaultConfig.Port,
	)

	db, err := dbutil.Open(&defaultConfig)

	if err != nil {
		log.Fatalf("Could not connect: %v\n", err)
		return
	}

	log.Println("Connection OK, setting up database...")

	_, err = dbutil.RunSQLFile(db, "dbutil/init.sql")
	if err != nil {
		log.Fatalf("Could not run db init: %v\n", err)
		return
	}

	log.Println("Database inititalized!")

	// remember to close the connection once we exit
	defer db.Close()
}
