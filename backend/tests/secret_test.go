package test

import (
	"os"
	"testing"
	"time"

	secret_model "secret-manager/backend/models/secret"
	user_model "secret-manager/backend/models/user"
	"secret-manager/backend/services/database"
)

func setupTestSecret(t *testing.T) {
	destroy_secrets := "drop table if exists secrets;"
	destroy_users := "drop table if exists users;"

	secret_query_bytes, err :=  os.ReadFile("../sql/02_secret/02_create_secret_table.sql");
	if err != nil {
		t.Fatal(err)
	}

	user_query_bytes, err :=  os.ReadFile("../sql/01_user/01_create_user_table.sql");
	if err != nil {
		t.Fatal(err)
	}

	queries := []string{
		destroy_secrets, destroy_users,
		string(user_query_bytes), string(secret_query_bytes),
	}

	RunTXQuery(t, queries)
}

func TestCreateSecret(t *testing.T) {
	setupTestSecret(t);

	db := database.GetConnection()
	user_db, err := user_model.Create(db, "secret_user_test", "somepass")

	if err != nil {
		t.Fatal(err)
	}

	user_id := *user_db.ID
	name := "test_secret"
	secret := "some super secret"
	encrypted := false
	secret_dto := secret_model.SecretDTO{
		UserID: user_id,
		Name: name,
		Secret: secret,
		Encrypted: encrypted,
	}
	secret_db, err := secret_model.Create(db, secret_dto)

	if err != nil {
		t.Fatal(err)
	}

	if secret_db.UserID != user_id {
		t.Fatal("Wrong user id")
	}

	if secret_db.Encrypted != encrypted {
		t.Fatal("Wrong encrypted flag")
	}

	if secret_db.Secret != secret {
		t.Fatal("Wrong secret")
	}

	if secret_db.Name != name {
		t.Fatal("Wrong secret name")
	}

	if time.Now().Compare(secret_db.CreatedAt) == -1 {
		t.Fatal("Secret created at is on future")
	}
}

func TestFindSecretByID(t *testing.T) {
	setupTestSecret(t)

	db := database.GetConnection()
	user_db, err := user_model.Create(db, "secret_user_test", "somepass")
	if err != nil {
		t.Fatal(err)
	}

	user_id := *user_db.ID
	secret_dto := secret_model.SecretDTO{
		UserID: user_id,
		Name: "some secret name",
		Secret: "some secret",
		Encrypted: false,
	}
	secret_db, err := secret_model.Create(db, secret_dto)
	if err != nil {
		t.Fatal(err)
	}

	secret_find_db, err := secret_model.FindByID(db, secret_db.ID)
	if err != nil {
		t.Fatal(err)
	}

	if err = secret_find_db.Compare(secret_db); err != nil {
		t.Fatal(err)
	}
}

func TestFindSecretByUserID(t *testing.T) {
	setupTestSecret(t)

	db := database.GetConnection()
	user_db, err := user_model.Create(db, "secret_user_test", "somepass")
	if err != nil {
		t.Fatal(err)
	}

	user_id := *user_db.ID
	secret_dto := secret_model.SecretDTO{
		UserID: user_id,
		Name: "some secret name",
		Secret: "some secret",
		Encrypted: false,
	}

	secrets_len := 10
	for i := 0; i < secrets_len; i++ {
		secret_model.Create(db, secret_dto)
	}

	secrets, err := secret_model.FindByUserID(db, user_id)

	if len(secrets) != secrets_len {
		t.Fatalf("The secrets length dont match")
	}
}

func TestDeleteByUserID(t *testing.T) {
	setupTestSecret(t)

	db := database.GetConnection()
	user_db, err := user_model.Create(db, "secret_user_test", "somepass")
	if err != nil {
		t.Fatal(err)
	}

	user_id := *user_db.ID
	secret_dto := secret_model.SecretDTO{
		UserID: user_id,
		Name: "some secret name",
		Secret: "some secret",
		Encrypted: false,
	}

	created_secret, err := secret_model.Create(db, secret_dto)
	if err != nil {
		t.Fatal(err)
	}

	deleted_secret, err := secret_model.DeleteById(db, created_secret.ID)
	if err != nil {
		t.Fatal(err)
	}

	if deleted_secret.ID != created_secret.ID {
		t.Fatalf("Deleted id dont match with created id (%v != %v)",
			deleted_secret, created_secret.ID)
	}

	find_secret, err := secret_model.FindByID(db, deleted_secret.ID)
	if find_secret != nil || err == nil {
		t.Fatalf("The secret should not be in the database ")
	}
}
