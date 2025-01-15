package test

import (
	"fmt"
	"testing"

	"secret-manager/backend/services/auth"
)

func TestGenerateToken(t *testing.T) {
	var id uint = 42
	token, err := auth.GenerateToken(id)

	if err != nil {
		t.Fatal(err)
	}

	if !auth.ValidateToken(token) {
		t.Fatalf("Invalid token")
	}
}

func TestValidateToken(t *testing.T) {
	var id uint = 42
	token, err := auth.GenerateToken(id)

	if err != nil {
		t.Fatal(err)
	}

	if !auth.ValidateToken(token) {
		t.Fatalf("Invalid token, should be valid")
	}

	if auth.ValidateToken("wrong token") {
		t.Fatalf("Valid token, should be invalid")
	}
}

func TestExtracTokenId(t *testing.T)  {
	var id uint = 42
	token, err := auth.GenerateToken(id)

	if err != nil {
		t.Fatal(err)
	}

	token_id, err := auth.ExtractTokenId(token)
	id_string := fmt.Sprintf("%v", id)

	if err != nil {
		t.Fatal(err)
	}

	if token_id != id_string {
		t.Fatal("Tokens dont match")
	}
}
