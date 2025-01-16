package test

import (
	"os"
	"testing"
	"time"

	user_model "secret-manager/backend/models/user"
	"secret-manager/backend/services/database"
)

func setupTestUser(t *testing.T) {
	query_bytes, err :=  os.ReadFile("../sql/01_user/01_create_user_table.sql");
	if err != nil {
		t.Fatal(err)
	}
	destroy_query := "drop table if exists secrets; drop table if exists users;"

	queries := []string { destroy_query, string(query_bytes)}
	RunTXQuery(t, queries)
}

func TestCreateUser(t *testing.T) {
	setupTestUser(t)

	db := database.GetConnection()
	test_username := "test_username"
	test_pass := "somepass"
	user, err := user_model.Create(db, test_username, test_pass)

	if err != nil {
		t.Fatalf("Test user create fail. Error: %v", err)
	}

	if !user.Active {
		t.Fatalf("User didnt initialized as active");
	}

	if user.Username != test_username {
		t.Fatalf("User username dont match");
	}

	if time.Now().Compare(user.CreatedAt) == -1 {
		t.Fatalf("User created at is on future");
	}
}

func TestFindByUserID(t *testing.T) {
	setupTestUser(t)

	db := database.GetConnection()
	new_user, err := user_model.Create(db, "test_username", "somepass")
	if err != nil {
		t.Fatalf("Test user create fail. Error: %v", err)
	}

	find_user, err := user_model.FindByID(db, *new_user.ID)
	if err != nil {
		t.Fatalf("Test user find by id fail. Error: %v", err)
	}

	if *new_user.ID != *find_user.ID {
		t.Fatalf("User ID not match.")
	}

	if new_user.Username != find_user.Username {
		t.Fatalf("User Username not match")
	}

	if new_user.Hash != find_user.Hash {
		t.Fatalf("User Hash not match")
	}

	if new_user.Active != find_user.Active {
		t.Fatalf("User Active not match")
	}

	if !new_user.CreatedAt.Equal(find_user.CreatedAt){
		t.Fatalf("User CreatedAt not match")
	}
}

func TestFindByUsername(t *testing.T) {
	setupTestUser(t)

	db := database.GetConnection()
	new_user, err := user_model.Create(db, "test_username", "somepass")
	if err != nil {
		t.Fatalf("Test user create fail. Error: %v", err)
	}

	find_user, err := user_model.FindByUsername(db, new_user.Username)
	if err != nil {
		t.Fatalf("Test user find by id fail. Error: %v", err)
	}

	if *new_user.ID != *find_user.ID {
		t.Fatalf("User ID not match.")
	}

	if new_user.Username != find_user.Username {
		t.Fatalf("User Username not match")
	}

	if new_user.Hash != find_user.Hash {
		t.Fatalf("User Hash not match")
	}

	if new_user.Active != find_user.Active {
		t.Fatalf("User Active not match")
	}

	if !new_user.CreatedAt.Equal(find_user.CreatedAt){
		t.Fatalf("User CreatedAt not match")
	}
}
