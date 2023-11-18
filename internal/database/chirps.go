package database

import "errors"

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
		return chirp, ErrNotExist
	}

	return chirp, nil
}
