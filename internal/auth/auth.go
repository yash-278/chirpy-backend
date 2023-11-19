package auth

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

// HashPassword -
func HashPassword(password string) (string, error) {
	dat, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(dat), nil
}

// CheckPasswordHash -
func CheckPasswordHash(password, hash string) error {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
}

func CreateJWT(key string, id, expiresInSeconds int) (string, error) {

	if expiresInSeconds == 0 {
		expiresInSeconds = 86400
	}

	currentTime := time.Now().UTC()
	issueTime := jwt.NewNumericDate(currentTime)
	expireTime := jwt.NewNumericDate(currentTime.Add(time.Second * time.Duration(expiresInSeconds)))

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.RegisteredClaims{
		Issuer:    "chirpy",
		IssuedAt:  issueTime,
		ExpiresAt: expireTime,
		Subject:   fmt.Sprint(id),
	})

	signedKey, err := token.SignedString([]byte(key))

	if err != nil {
		return "", err
	}

	return signedKey, nil
}

func ValidateJWT(jwtStr, key string) (string, error) {

	token, err := jwt.ParseWithClaims(jwtStr, jwt.MapClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(key), nil
	})
	if err != nil {
		return "", err
	}

	data, err := token.Claims.GetSubject()
	if err != nil {
		return "", err
	}

	return data, err
}
