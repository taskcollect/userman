package dbutil

import (
	"database/sql"
	"fmt"

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

// you should use this when you want to modify the database
// starts a transaction, runs the query and commits the transaction
// can always rollback if something goes wrong
func RunQueryFailsafe(db *sql.DB, query string, args ...interface{}) (*sql.Rows, error) {
	// initialize transaction
	tx, err := db.Begin()
	if err != nil {
		return nil, err
	}

	stmt, err := tx.Prepare(query)
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	// execute the query
	rows, err := stmt.Query(args...)
	if err != nil {
		stmt.Close()
		tx.Rollback()
		return nil, err
	}

	// commit any changes made
	err = tx.Commit()
	if err != nil {
		// something went wrong with committing
		stmt.Close()
		tx.Rollback()
		return nil, err
	}
	return rows, nil
}

func Initialize(db *sql.DB) error {
	q := `
	CREATE TABLE IF NOT EXISTS users (
		username VARCHAR(32) NOT NULL PRIMARY KEY,
		secret VARCHAR(96) NOT NULL,
		credentials JSON NOT NULL,
		preferences JSON NOT NULL,
		registeredOn TIMESTAMP NOT NULL,
		lastLogin TIMESTAMP
	);
	`

	_, err := RunQueryFailsafe(db, q)
	return err
}
