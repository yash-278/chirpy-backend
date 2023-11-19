package auth

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

var AccessTokenExpireTime int = 1
var RefreshTokenExpireTime int = 1440

type IssuerType struct {
	AccessIssuer  string
	RefreshIssuer string
}

var Issuer = IssuerType{
	AccessIssuer:  "chirpy-access",
	RefreshIssuer: "chirpy-refresh",
}

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

func CreateJWT(key, issuer string, id, expireTimeInSeconds int) (string, error) {

	currentTime := time.Now().UTC()
	issueTime := jwt.NewNumericDate(currentTime)
	expireTime := jwt.NewNumericDate(currentTime.Add(time.Hour * time.Duration(expireTimeInSeconds)))

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.RegisteredClaims{
		Issuer:    issuer,
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

func ValidateJWT(jwtStr, key string) (jwt.Claims, error) {

	token, err := jwt.ParseWithClaims(jwtStr, jwt.MapClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(key), nil
	})
	if err != nil {
		return jwt.MapClaims{}, err
	}

	return token.Claims, err
}
