package database

import (
	"encoding/json"
	"errors"
	"os"
	"sync"

	"golang.org/x/crypto/bcrypt"
)

type DBStructure struct {
	Chirps map[int]Chirp `json:"chirps"`
	Users  map[int]User  `json:"users"`
}

type Chirp struct {
	Id   int    `json:"id"`
	Body string `json:"body"`
}

type User struct {
	Id       int    `json:"id"`
	Email    string `json:"email"`
	Password []byte `json:"password"`
}

type DB struct {
	path string
	mux  *sync.RWMutex
}

type ResponseUser struct {
	Id    int    `json:"id"`
	Email string `json:"email"`
}

// NewDB creates a new database connection
// and creates the database file if it doesn't exist
func NewDB(path string) (*DB, error) {
	DbConfig := DB{
		path: path,
		mux:  &sync.RWMutex{},
	}

	_, err := os.ReadFile(path)

	if err == os.ErrNotExist {
		err = os.WriteFile(path, nil, 0666)
		if err != nil {
			return nil, errors.New("database cannot be created")
		}
	}

	return &DbConfig, nil
}

// CreateChirp creates a new chirp and saves it to disk
func (db *DB) CreateChirp(body string) (Chirp, error) {
	DBStruct, _ := db.loadDB()
	chirpLen := len(DBStruct.Chirps)
	chirp := Chirp{}

	newId := chirpLen + 1
	chirp = Chirp{
		Id:   newId,
		Body: body,
	}
	DBStruct.Chirps[newId] = chirp

	db.writeDB(DBStruct)

	return chirp, nil
}

// GetChirps returns all chirps in the database
func (db *DB) GetChirps() ([]Chirp, error) {
	chirps := []Chirp{}

	DbData, err := db.loadDB()
	if err != nil {
		return chirps, errors.New("db cannot be loaded")
	}

	for _, chirp := range DbData.Chirps {
		chirps = append(chirps, chirp)
	}

	return chirps, nil
}

// GetChirps returns a chirp in the database based on Id
func (db *DB) GetChirpById(chirpId int) (Chirp, error) {
	chirp := Chirp{}

	DbData, err := db.loadDB()
	if err != nil {
		return chirp, errors.New("db cannot be loaded")
	}

	chirp, ok := DbData.Chirps[chirpId]
	if !ok {
		return chirp, errors.New("chirp not found")
	}

	return chirp, nil
}

// CreateChirp creates a new chirp and saves it to disk
func (db *DB) CreateUser(email, password string) (ResponseUser, error) {
	DBStruct, _ := db.loadDB()

	_, err := db.getUser(email)
	if err != nil {
		return ResponseUser{}, errors.New("user with the email already exists")
	}

	newId := len(DBStruct.Users) + 1

	hash, err := bcrypt.GenerateFromPassword([]byte(password), 10)
	if err != nil {
		return ResponseUser{}, errors.New("password Has Failed")
	}

	DBStruct.Users[newId] = User{
		Id:       newId,
		Email:    email,
		Password: hash,
	}

	err = db.writeDB(DBStruct)
	if err != nil {
		return ResponseUser{}, err
	}

	return ResponseUser{
		Id:    newId,
		Email: email,
	}, nil
}

func (db *DB) LoginUser(email, password string) (ResponseUser, error) {
	user, err := db.getUser(email)
	if err != nil {
		return ResponseUser{}, errors.New("user does not exists")
	}

	err = bcrypt.CompareHashAndPassword(user.Password, []byte(password))
	if err != nil {
		return ResponseUser{}, errors.New("password is incorrect")
	}

	return ResponseUser{
		user.Id,
		user.Email,
	}, nil
}

// ensureDB creates a new database file if it doesn't exist
func (db *DB) ensureDB() error {
	_, err := os.ReadFile(db.path)
	if errors.Is(err, os.ErrNotExist) {
		dbData := DBStructure{
			Chirps: map[int]Chirp{},
			Users:  map[int]User{},
		}

		rawData, err := json.Marshal(dbData)
		if err != nil {
			return err
		}

		err = os.WriteFile(db.path, rawData, 0666)
		if err != nil {
			return errors.New("database cannot be created")
		}
	}

	return nil
}

// loadDB reads the database file into memory
func (db *DB) loadDB() (DBStructure, error) {
	db.mux.RLock()
	database := DBStructure{}
	err := db.ensureDB()
	if err != nil {
		return database, errors.New("db cannot be loaded")
	}

	rawData, err := os.ReadFile(db.path)
	if err != nil {
		return database, errors.New("file cannot be read")
	}
	db.mux.RUnlock()

	err = json.Unmarshal(rawData, &database)
	return database, nil
}

// writeDB writes the database file to disk
func (db *DB) writeDB(dbStructure DBStructure) error {
	db.mux.Lock()
	defer db.mux.Unlock()

	err := db.ensureDB()
	if err != nil {
		return err
	}

	rawJson, err := json.Marshal(dbStructure)
	if err != nil {
		return err
	}

	err = os.WriteFile(db.path, rawJson, 0666)
	if err != nil {
		return err
	}

	return nil
}

func (db *DB) getUser(email string) (User, error) {
	user := User{}
	database, err := db.loadDB()

	if err != nil {
		return user, errors.New("DB not loaded")
	}

	if len(database.Users) == 0 {
		return user, nil
	}

	for _, u := range database.Users {
		if u.Email == email {
			user = u
		}
		return user, nil
	}

	return user, errors.New("user not found")
}
