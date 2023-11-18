package database

import (
	"errors"
)

type User struct {
	Id             int    `json:"id"`
	Email          string `json:"email"`
	HashedPassword string `json:"hashed_password"`
}

var ErrAlreadyExists = errors.New("already exists")

// CreateChirp creates a new chirp and saves it to disk
func (db *DB) CreateUser(email, hashed_password string) (User, error) {
	DBStruct, _ := db.loadDB()

	if _, err := db.GetUserByEmail(email); !errors.Is(err, ErrNotExist) {
		return User{}, ErrAlreadyExists
	}

	newId := len(DBStruct.Users) + 1

	newUser := User{
		Id:             newId,
		Email:          email,
		HashedPassword: hashed_password,
	}

	DBStruct.Users[newId] = newUser

	err := db.writeDB(DBStruct)
	if err != nil {
		return User{}, err
	}

	return newUser, nil
}

func (db *DB) GetUser(id int) (User, error) {
	dbStructure, err := db.loadDB()
	if err != nil {
		return User{}, err
	}

	user, ok := dbStructure.Users[id]
	if !ok {
		return User{}, ErrNotExist
	}

	return user, nil
}

func (db *DB) GetUserByEmail(email string) (User, error) {
	dbStructure, err := db.loadDB()
	if err != nil {
		return User{}, err
	}

	for _, user := range dbStructure.Users {
		if user.Email == email {
			return user, nil
		}
	}

	return User{}, ErrNotExist
}
