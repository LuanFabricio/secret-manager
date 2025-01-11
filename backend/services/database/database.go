package database

import (
	"database/sql"
	"log"
	"os"

	_ "github.com/lib/pq"
)

type Row interface {
	Scan(dest ...any) error;
}

type Database interface {
	QueryRow(query string, args ...any) *sql.Row;
	Query(query string, args ...any) (*sql.Rows, error);
	Begin() (*sql.Tx, error);
}

// NOTE: Maybe move back to *sql.DB
var db Database = nil;
func GetConnection() Database {
	if db == nil {
		db = bootConnection()
	}
	return db;
}

func bootConnection() *sql.DB {
	PSQL_CONNECTION_ENV := "SM_PSQL_CONNECTION"
	psql_connection, found := os.LookupEnv(PSQL_CONNECTION_ENV);

	log.Printf("PSQL Connection(connected? %v): %v\n", found, psql_connection)
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
