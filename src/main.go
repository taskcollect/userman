package main

import (
	"database/sql"
	"log"
	"main/dbutil"
	"main/handlers"
	"net/http"
	"os"
	"strconv"
)

type ServerConfig struct {
	Sql      dbutil.SqlConfig
	BindAddr string
}

// server config, values here will get overriden by env
var config = ServerConfig{
	BindAddr: "0.0.0.0:2000",
}

func makeMux(db *sql.DB) *http.ServeMux {
	mux := http.NewServeMux()

	handler := handlers.NewBaseHandler(db)

	mux.HandleFunc("/v1/register", handler.NewUser)
	mux.HandleFunc("/v1/get", handler.GetUser)
	mux.HandleFunc("/v1/alter", handler.AlterUser)
	mux.HandleFunc("/v1/delete", handler.DeleteUser)

	return mux
}

func configure(c *ServerConfig) {
	bindAddr, exists := os.LookupEnv("BIND_ADDR")
	if exists {
		if bindAddr == "" {
			log.Fatalln("(cfg) empty bind address supplied, cannot bind")
		}
		c.BindAddr = bindAddr
	} else {
		log.Printf("(cfg) no bind address supplied, defaulting to '%s'", c.BindAddr)
	}

	dbHost, exists := os.LookupEnv("DB_HOST")
	if !exists {
		log.Fatalln("(cfg) no db host supplied (from DB_HOST)")
	}
	if dbHost == "" {
		log.Fatalln("(cfg) empty db host supplied")
	}
	c.Sql.Host = dbHost

	dbPortStr, exists := os.LookupEnv("DB_PORT")
	if !exists {
		log.Fatalln("(cfg) no db Port supplied (from DB_PORT)")
	}

	dbPortInt, err := strconv.Atoi(dbPortStr)
	if err != nil {
		log.Fatalln("(cfg) given port int conversion failed: " + err.Error())
	}

	if dbPortInt > 65535 || dbPortInt < 0 {
		log.Fatalln("(cfg) given port not in range 0-65535")
	}

	c.Sql.Port = uint16(dbPortInt)

	dbUser, exists := os.LookupEnv("DB_USER")
	if !exists {
		log.Fatalln("(cfg) no db user supplied (from DB_USER)")
	}
	if dbUser == "" {
		log.Fatalln("(cfg) empty db user supplied")
	}
	c.Sql.User = dbUser

	// this will be empty if the variable is not set, which in this case is ok
	c.Sql.Password = os.Getenv("DB_PASS")

	dbName, exists := os.LookupEnv("DB_NAME")
	if !exists {
		log.Fatalln("(cfg) no db name supplied (from DB_NAME)")
	}
	if dbName == "" {
		log.Fatalln("(cfg) empty db name supplied")
	}
	c.Sql.Database = dbName
}

func main() {
	log.Printf("Initializing config from environment variables...")

	configure(&config)

	log.Printf(
		"Trying to open connection to db '%s' as %s@%s:%d",
		config.Sql.Database, config.Sql.User, config.Sql.Host, config.Sql.Port,
	)

	db, err := dbutil.Open(&config.Sql)

	if err != nil {
		log.Fatalf("Could not connect: %s\n", err.Error())
	}

	// remember to close the connection once we exit
	defer db.Close()

	log.Println("Connection OK, setting up database...")

	err = dbutil.Initialize(db)
	if err != nil {
		log.Fatalf("Could not run db init: %s\n", err.Error())
	}

	log.Printf("Database inititalized, starting server binded to %s...", config.BindAddr)

	mux := makeMux(db)
	http.ListenAndServe(config.BindAddr, handlers.RequestLogger(mux))

	log.Println("Server exited. Cleaning up...")
}
