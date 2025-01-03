package user

import (
	"crypto/sha256"
	"database/sql"
	"fmt"
)

type User struct {
	ID *uint `json:"id"`
	Username string `json:"username"`
	Password string `json:"password"`
}

func Create(db* sql.DB, username string, password string) (*User, error) {
	hash := sha256.Sum256([]byte(password));
	hash_string := fmt.Sprintf("%x", hash);

	var db_user User;
	var err = db.QueryRow(
		`INSERT INTO users (username, hash)
		VALUES ($1, $2)
		RETURNING id, username`,
		username, hash_string,
	).Scan(&db_user.ID, &db_user.Username);

	if err != nil {
		return nil, err;
	}

	return &db_user, nil;
}

func FindByID(db *sql.DB, id uint) (*User, error) {
	var find_user User;
	var err = db.QueryRow(
		`SELECT	id, username FROM users
		WHERE id = $1`,
		id,
	).Scan(&find_user.ID, &find_user.Username);

	if err != nil {
		return nil, err;
	}

	return &find_user, nil;
}
