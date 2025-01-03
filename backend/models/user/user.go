package user

import (
	"crypto/sha256"
	"database/sql"
	"fmt"
)

type UserDTO struct {
	ID *uint `json:"id"`
	Username string `json:"username"`
	Password string `json:"password"`
}

type UserDB struct {
	ID *uint `json:"id"`
	Username string `json:"username"`
	Hash string `json:"hash"`
}

func Create(db* sql.DB, username string, password string) (*UserDB, error) {
	hash := sha256.Sum256([]byte(password));
	hash_string := fmt.Sprintf("%x", hash);

	var db_user UserDB;
	var err = db.QueryRow(
		`INSERT INTO users (username, hash)
		VALUES ($1, $2)
		RETURNING id, username, hash`,
		username, hash_string,
	).Scan(&db_user.ID, &db_user.Username, &db_user.Hash);

	if err != nil {
		return nil, err;
	}

	return &db_user, nil;
}

func FindByID(db *sql.DB, id uint) (*UserDB, error) {
	var find_user UserDB;
	row := db.QueryRow(
		`SELECT	id, username, hash FROM users
		WHERE id = $1`,
		id,
	);

	err := row.Scan(&find_user.ID, &find_user.Username, &find_user.Hash);

	if err != nil {
		return nil, err;
	}

	return &find_user, nil;
}
