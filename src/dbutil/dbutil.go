package dbutil

import (
	"database/sql"
	"fmt"
	"os"

	_ "github.com/lib/pq"
)

type SqlConfig struct {
	Host     string
	Port     uint16
	User     string
	Password string
	Database string
}

func getConnectionString(conf *SqlConfig) string {
	return fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		conf.Host, conf.Port, conf.User, conf.Password, conf.Database,
	)
}

func Open(config *SqlConfig) (*sql.DB, error) {
	// construct a connection string and try to open the connection
	db, err := sql.Open("postgres", getConnectionString(config))
	if err != nil {
		return nil, err
	}

	// try to ping the database and see if it responds (tests network)
	err = db.Ping()
	if err != nil {
		return nil, err
	}

	// everything went well, give back the connection
	return db, nil
}

func RunSQLFile(db *sql.DB, fname string, sqlArgs ...interface{}) (*sql.Rows, error) {
	// read the file
	buf, err := os.ReadFile(fname)
	if err != nil {
		return nil, err
	}
	return db.Query(string(buf), sqlArgs...)
}
