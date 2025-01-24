package secret

import (
	"errors"
	"fmt"
	"time"

	"secret-manager/backend/services/database"

	"github.com/gin-gonic/gin"
)


type SecretDTO struct {
	ID        *uint `json:"id"`
	UserID	  uint  `json:"user_id" binding:"required"`
	Name      string `json:"name" binding:"required"`
	Secret    string `json:"secret" binding:"required"`
	Encrypted bool `json:"encrypted"`
	EncryptionKey string `json:"encryption_key"`
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

func (sd *SecretDB) Compare(other *SecretDB) error {
	if sd.ID != other.ID {
		return errors.New(fmt.Sprintf("Wrong ID %d != %d", sd.ID, other.ID))
	}

	if sd.Name != other.Name {
		return errors.New(fmt.Sprintf("Wrong name %v != %v", sd.Name, other.Name))
	}

	if sd.UserID != other.UserID {
		return errors.New(
			fmt.Sprintf("Wrong user id %d != %d", sd.UserID, other.UserID))
	}

	if sd.Secret != other.Secret {
		return errors.New(
			fmt.Sprintf("Wrong secret %v != %v", sd.Secret, other.Secret))
	}

	if sd.Encrypted != other.Encrypted {
		return errors.New(
			fmt.Sprintf("Wrong secret %v != %v", sd.Encrypted, other.Encrypted))
	}

	if sd.CreatedAt.Compare(other.CreatedAt) != 0 {
		return errors.New(
			fmt.Sprintf("Wrong secret %v != %v", sd.CreatedAt, other.CreatedAt))
	}

	return nil
}

func (sd *SecretDB) loadFromRow(row database.Row) error {
	return row.Scan(&sd.ID,
		&sd.UserID,
		&sd.Name,
		&sd.Secret,
		&sd.Encrypted,
		&sd.CreatedAt,
	)
}

// TODO: Add a secret encryption (sync) support
func Create(db database.Database, secret SecretDTO) (*SecretDB, error) {
	if secret.Encrypted /*&& secret.EncryptionKey == ""*/{
		return nil, errors.New("Secret encryption dont have support, yet")
		// return nil, errors.New("The encryption key should be passed")
	}

	row := db.QueryRow(
		`INSERT INTO secrets (user_id, name, secret, encrypted)
		VALUES ($1, $2, $3, $4)
		RETURNING id, user_id, name, secret, encrypted, created_at`,
		secret.UserID, secret.Name, secret.Secret, secret.Encrypted,
	)

	var new_secret SecretDB;
	err := new_secret.loadFromRow(row)
	if err != nil {
		return nil, err
	}

	return &new_secret, nil
}

func FindByID(db database.Database, secret_id uint) (*SecretDB, error) {
	var find_secret SecretDB

	row := db.QueryRow(
		`SELECT
			s.id,
			s.user_id,
			s.name,
			s.secret,
			s.encrypted,
			s.created_at
		FROM secrets s
		WHERE s.id = $1`,
		secret_id,
	)

	err := find_secret.loadFromRow(row)
	if err != nil {
		return nil, err
	}

	return &find_secret, nil
}

func FindByUserID(db database.Database, user_id uint) ([]SecretDB, error) {
	row, err := db.Query(
		`SELECT
			s.id,
			s.user_id,
			s.name,
			s.secret,
			s.encrypted,
			s.created_at
		FROM secrets s
		WHERE s.user_id = $1`,
		user_id,
	)
	if err != nil {
		return nil, err
	}

	secrets := make([]SecretDB, 0)
	for row.Next() {
		var secret_db SecretDB
		err = secret_db.loadFromRow(row)
		if err != nil {
			return nil, err
		}

		secrets = append(secrets, secret_db)
	}

	return secrets, nil
}

func DeleteById(db database.Database, secret_id uint) (*SecretDB, error) {
	row := db.QueryRow(
		`DELETE FROM secrets
		WHERE id = $1
		RETURNING
			id,
			user_id,
			name,
			secret,
			encrypted,
			created_at
		`,
		secret_id,
	)

	var secret_db SecretDB
	if err := secret_db.loadFromRow(row); err != nil {
		return nil, err
	}

	return &secret_db, nil
}

func UpdateByID(db database.Database, secret_id uint, name, secret string) (*SecretDB, error) {
	row := db.QueryRow(
		`UPDATE secrets
		SET name = $2,
		    secret = $3
		WHERE id = $1
		RETURNING
			id,
			user_id,
			name,
			secret,
			encrypted,
			created_at
		`,
		secret_id, name, secret,
	)

	var updated_secret SecretDB
	err := updated_secret.loadFromRow(row)
	if err != nil {
		return nil, err
	}

	return &updated_secret, nil
}
