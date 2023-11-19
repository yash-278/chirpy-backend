package database

import (
	"errors"
	"time"
)

type RefreshToken struct {
	Token string    `json:"token"`
	Time  time.Time `json:"time"`
}

// CreateChirp creates a new chirp and saves it to disk
func (db *DB) RevokeRefreshToken(token string) error {
	DBStruct, _ := db.loadDB()

	_, err := db.GetRefreshTokenById(token)
	if err == nil {
		return errors.New("RefreshToken exists")
	}

	DBStruct.RefreshTokens[token] = RefreshToken{
		Token: token,
		Time:  time.Now().UTC(),
	}

	db.writeDB(DBStruct)

	return nil
}

// GetRefreshTokenById returns a RefreshToken in the database based on Id
func (db *DB) GetRefreshTokenById(token string) (RefreshToken, error) {
	refresh_token := RefreshToken{}

	DbData, err := db.loadDB()
	if err != nil {
		return refresh_token, errors.New("db cannot be loaded")
	}

	refresh_token, ok := DbData.RefreshTokens[token]
	if !ok {
		return refresh_token, ErrNotExist
	}

	return refresh_token, nil
}
