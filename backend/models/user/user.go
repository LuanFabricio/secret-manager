package user

import (
	"crypto/sha256"
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/gin-gonic/gin"

	"secret-manager/backend/services/database"
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
	CreatedAt time.Time `json:"created_at"`
	Active bool `json:"active"`
}

func (ud *UserDB) ToH() gin.H {
	return gin.H{
		"id": ud.ID,
		"username": ud.Username,
		"created_at": ud.CreatedAt,
		"active": ud.Active,
	};
}

func Create(db database.Database, username string, password string) (*UserDB, error) {
	salt, found := os.LookupEnv("SALT");
	if !found {

		return nil, errors.New("Salt secret not found");
	}

	salted_password := salt + password;
	hash := sha256.Sum256([]byte(salted_password));
	hash_string := fmt.Sprintf("%x", hash);

	var db_user UserDB;
	var err = db.QueryRow(
		`INSERT INTO users (username, hash)
		VALUES ($1, $2)
		RETURNING id, username, hash, created_at, active`,
		username, hash_string,
	).Scan(&db_user.ID,
		&db_user.Username,
		&db_user.Hash,
		&db_user.CreatedAt,
		&db_user.Active,
	);

	if err != nil {
		return nil, err;
	}

	return &db_user, nil;
}

func FindByID(db database.Database, id uint) (*UserDB, error) {
	var find_user UserDB;
	row := db.QueryRow(
		`SELECT	id,
			username,
			hash,
			created_at,
			active
		FROM users
		WHERE id = $1`,
		id,
	);

	err := row.Scan(
		&find_user.ID,
		&find_user.Username,
		&find_user.Hash,
		&find_user.CreatedAt,
		&find_user.Active,
	);

	if err != nil {
		return nil, err;
	}

	return &find_user, nil;
}
