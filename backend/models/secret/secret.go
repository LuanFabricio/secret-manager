package secret

import (
	"time"

	"secret-manager/backend/services/database"

	"github.com/gin-gonic/gin"
)


type SecretDTO struct {
	ID        *uint `json:"id"`
	UserID	  uint  `json:"user_id"`
	Name      string `json:"name"`
	Secret    string `json:"secret"`
	Encrypted bool `json:"encrypted"`
	CreatedAt *time.Time `json:"created_at"`
}

type SecretDB struct {
	ID        uint `json:"id"`
	UserID	  uint  `json:"user_id"`
	Name      string `json:"name"`
	Secret    string `json:"secret"`
	Encrypted bool `json:"encrypted"`
	CreatedAt time.Time `json:"created_at"`
}

func (sd *SecretDB) ToH() gin.H {
	return gin.H{
		"id": sd.ID,
		"user_id": sd.UserID,
		"name": sd.Name,
		"secret": sd.Secret,
		"encrypted": sd.Encrypted,
		"created_at": sd.CreatedAt,

	};
}

// TODO: Add a secret encryption (sync)
func Create(db database.Database, secret SecretDTO) (*SecretDB, error) {
	var new_secret SecretDB;

	err := db.QueryRow(
		`INSERT INTO secrets (user_id, name, secret, encrypted)
		VALUES ($1, $2, $3, $4)
		RETURNING id, user_id, name, secret, encrypted, created_at`,
		secret.UserID, secret.Name, secret.Secret, secret.Encrypted,
	).Scan(&new_secret.ID,
		&new_secret.UserID,
		&new_secret.Name,
		&new_secret.Secret,
		&new_secret.Encrypted,
		&new_secret.CreatedAt,
	)

	if err != nil {
		return nil, err
	}

	return &new_secret, nil
}
