package test

import (
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
