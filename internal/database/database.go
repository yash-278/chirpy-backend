package database

import (
	"encoding/json"
	"errors"
	"os"
	"sync"
)

var ErrNotExist = errors.New("resource does not exist")

type DBStructure struct {
	Chirps        map[int]Chirp           `json:"chirps"`
	Users         map[int]User            `json:"users"`
	RefreshTokens map[string]RefreshToken `json:"refresh_tokens"`
}

type Chirp struct {
	Id       int    `json:"id"`
	Body     string `json:"body"`
	AuthorId int    `json:"author_id"`
}

type DB struct {
	path string
	mux  *sync.RWMutex
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

// ensureDB creates a new database file if it doesn't exist
func (db *DB) ensureDB() error {
	_, err := os.ReadFile(db.path)
	if errors.Is(err, os.ErrNotExist) {
		dbData := DBStructure{
			Chirps:        map[int]Chirp{},
			Users:         map[int]User{},
			RefreshTokens: map[string]RefreshToken{},
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
