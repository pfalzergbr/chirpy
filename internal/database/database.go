package database

import (
	"encoding/json"
	"fmt"
	"os"
	"sync"
)

type DB struct {
	path string
	mux  *sync.RWMutex
}

type Chirp struct {
	Id   int    `json:"id"`
	Body string `json:"body"`
}

type DBStructure struct {
	Chirps map[int]Chirp `json:"chirps"`
}

func NewDB(path string) (*DB, error) {
	db := &DB{
		path: path,
		mux:  &sync.RWMutex{},
	}

	err := db.ensureDB()
	if err != nil {
		return nil, err
	}
	return db, nil

}

func (db *DB) ensureDB() error {
	_, err := os.ReadFile(db.path)

	if err != nil {
		err := os.WriteFile(db.path, []byte(`{"chirps":{}}`), 0644)
		fmt.Printf("File not yet exist: %v\n", err)

		if err != nil {
			return err
		}
	}

	return nil
}

func (db *DB) GetChirps() ([]Chirp, error){
	dbStruct, err := db.loadDB()
	if err != nil {
		return nil, err
	}

	chirps := make([]Chirp, 0, len(dbStruct.Chirps))

	for _, chirp := range dbStruct.Chirps {
		chirps = append(chirps, chirp)
	}

	return chirps, nil
}

func (db *DB) CreateChirp(body string) (Chirp, error) {
	dbStruct, err := db.loadDB()
	if err != nil {
		return Chirp{}, err
	}

	id := len(dbStruct.Chirps) + 1
	chirp := Chirp{
		Id:   id,
		Body: body,
	}

	dbStruct.Chirps[id] = chirp

	err = db.writeDB(dbStruct)
	if err != nil {
		return Chirp{}, err
	}

	return chirp, nil
}

func (db *DB) loadDB() (DBStructure, error) {
	db.mux.RLock()
	defer db.mux.RUnlock()

	data, err := os.ReadFile(db.path)
	if err != nil {
		return DBStructure{}, err
	}

	dbStruct := DBStructure{}
	err = json.Unmarshal(data, &dbStruct)

	if err != nil {
		return DBStructure{}, err
	}

	return dbStruct, nil
}

func (db *DB) writeDB(dbStructure DBStructure) error {
	db.mux.Lock()
	defer db.mux.Unlock()

	json, err := json.Marshal(dbStructure)

	if err != nil {
		return err
	}

	err = os.WriteFile(db.path, []byte(json), 0644)
	if err != nil {
		return err
	}
	return nil
}
