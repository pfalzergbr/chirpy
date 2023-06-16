package database

import (
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
