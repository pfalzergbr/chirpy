package database

import (
	"encoding/json"
	"fmt"
	"os"
	"sync"

	"golang.org/x/crypto/bcrypt"
)

type DB struct {
	path string
	mux  *sync.RWMutex
}

type Chirp struct {
	Id   int    `json:"id"`
	Body string `json:"body"`
}

type User struct {
	Id       int    `json:"id"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type UserResponse struct {
	Id       int    `json:"id"`
	Email    string `json:"email"`
}

type DBStructure struct {
	Chirps map[int]Chirp `json:"chirps"`
	Users  map[int]User  `json:"users"`
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
		err := os.WriteFile(db.path, []byte(`{"chirps":{}, "users":{}}`), 0644)
		fmt.Printf("File not yet exist: %v\n", err)

		if err != nil {
			return err
		}
	}

	return nil
}

func (db *DB) GetChirps() (DBStructure, error) {
	dbStruct, err := db.loadDB()
	if err != nil {
		return DBStructure{}, err
	}

	return dbStruct, nil
}

func (db *DB) CreateUser(email string, password string) (UserResponse, error) {
	dbStruct, err := db.loadDB()
	if err != nil {
		return UserResponse{}, err
	}

	_, err = db.GetUserByEmail(email)

	if err == nil {
		return UserResponse{}, fmt.Errorf("user with email %s already exists", email)
	}

	id := len(dbStruct.Users) + 1

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

	if err != nil {
		return UserResponse{}, err
	}

	user := User{
		Id:    id,
		Email: email,
		Password: string(hashedPassword),
	}

	dbStruct.Users[id] = user

	err = db.writeDB(dbStruct)
	if err != nil {
		return UserResponse{}, err
	}
	
	userResponse := UserResponse{
		Id:    user.Id,
		Email: user.Email,
	}

	return userResponse, nil
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

func (db *DB) GetUserByEmail(email string) (User, error) {
	dbStruct, err := db.loadDB()
	if err != nil {
		return User{}, err
	}

	for _, user := range dbStruct.Users {
		if user.Email == email {
			return user, nil
		}
	}
	return User{}, fmt.Errorf("User not found")
}
