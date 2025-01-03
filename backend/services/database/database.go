package database

import (
	"database/sql"
	"log"
	"os"

	_ "github.com/lib/pq"
)

var db *sql.DB = nil;
func GetConnection() *sql.DB {
	if db == nil {
		db = bootConnection()
	}
	return db;
}

func bootConnection() *sql.DB {
	PSQL_CONNECTION_ENV := "PSQL_CONNECTION"
	psql_connection, found := os.LookupEnv(PSQL_CONNECTION_ENV);

	if !found {
		log.Fatalf("Could not find \"%v\"\n", PSQL_CONNECTION_ENV);
	}

	db, err := sql.Open("postgres", psql_connection);
	if err != nil {
		log.Fatalf("Could not open PostgresSQL connection: %v\n", err);
		return nil;
	}

	return db;
}
