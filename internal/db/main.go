package db

import (
	"database/sql"
	"fmt"
	_ "modernc.org/sqlite"
	"os"
)

var database *sql.DB

func init() {
	var err error
	database, err = sql.Open("sqlite", "db.sqlite")
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	err = Migrate(database)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func Database() *sql.DB {
	return database
}
