package main

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/golang-jwt/jwt/v5"
	"github.com/yash-278/chirpy-backend/auth"
)

func (cfg *apiConfig) UtilAuthRequest(w http.ResponseWriter, r *http.Request, issuer string) (int, jwt.Claims, error) {
	token := GetBearerToken(r)
	if token == "" {
		return 0, jwt.MapClaims{}, errors.New("bad token")
	}

	claims, err := auth.ValidateJWT(token, cfg.secretKey)
	if err != nil {
		return 0, jwt.MapClaims{}, err
	}
	if sub, _ := claims.GetIssuer(); sub != issuer {
		return 0, jwt.MapClaims{}, errors.New("unauthorized")
	}

	strId, _ := claims.GetSubject()

	id, err := strconv.Atoi(strId)
	if err != nil {
		return 0, jwt.MapClaims{}, err
	}

	return id, claims, nil
}
