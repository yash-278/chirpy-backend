package database

import (
	"errors"
	"fmt"
)

type User struct {
	Id             int    `json:"id"`
	Email          string `json:"email"`
	HashedPassword string `json:"hashed_password"`
}

var ErrAlreadyExists = errors.New("already exists")

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

func (db *DB) UpdateUser(id int, email, hashedPassword string) (User, error) {
	DBStruct, _ := db.loadDB()

	user, err := db.GetUser(id)
	if err != nil {
		return User{}, err
	}

	// Update fields only if they are not empty
	if email != "" {
		user.Email = email
	}

	if hashedPassword != "" {
		user.HashedPassword = hashedPassword
	}

	// If you have more fields, add similar conditions for each field
	fmt.Println(email)
	DBStruct.Users[id] = user

	err = db.writeDB(DBStruct)
	if err != nil {
		return User{}, err
	}

	return user, nil
}
