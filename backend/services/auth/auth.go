package auth

import (
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/dgrijalva/jwt-go"
)

const TOKEN_LIFE_SPAN = 4
const TOKEN_ENV_VAR = "SM_TOKEN_SECRET"

func GenerateToken(id uint) (string, error) {
	claims := jwt.MapClaims{}
	claims["authorized"] = true
	claims["id"] = id
	claims["exp"] = time.Now().Add(time.Duration(TOKEN_LIFE_SPAN) * time.Hour)

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(os.Getenv(TOKEN_ENV_VAR)));
}

func ValidateToken(userToken string) bool {
	_, err := parseToken(userToken)

	return err == nil
}

func ExtractTokenId(userToken string) (string, error) {
	parsedToken, err := parseToken(userToken)
	if err != nil {
		return "", err
	}

	if !parsedToken.Valid {
		return "", errors.New("Invalid token")
	}

	claims, ok := parsedToken.Claims.(jwt.MapClaims)
	if !ok {
		return "", errors.New("Failed to claim the parsed token")
	}

	return fmt.Sprintf("%s", claims["id"]), nil
}

func parseToken(userToken string) (*jwt.Token, error) {
	return jwt.Parse(userToken, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(os.Getenv(TOKEN_ENV_VAR)), nil
	})
}
