package test

import (
	"testing"

	"github.com/joho/godotenv"

	"secret-manager/backend/services/database"
)

func RunTXQuery(t *testing.T, queries []string) {
	godotenv.Load()
	db := database.GetConnection()

	tx, err := db.Begin()
	if err != nil {
		t.Fatalf("[Create table - create] Error: %v\n", err)
	}

	for _, query := range queries{
		_, err = tx.Exec(query)
		if err != nil {
			tx.Rollback()
			t.Fatalf("[Create table - exec] Error: %v\n", err)
		}
	}

	err = tx.Commit()
	if err != nil {
		t.Fatalf("[Create table - commit] Error: %v\n", err)
	}
}
