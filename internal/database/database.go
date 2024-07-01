package database

import (
	"encoding/json"
	"log"
	"os"
	"sync"
)

type DB struct {
	path string
	mu   *sync.RWMutex
}

type DBStructure struct {
	Chirps map[int]Chirp `json:"chirps"`
}

type Chirp struct {
	Id   int    `json:"id"`
	Body string `json:"body"`
}

func NewDB(path string) (*DB, error) {
	return &DB{
		path: path,
		mu:   &sync.RWMutex{},
	}, nil
}

func (db *DB) CreateChirp(body string) (Chirp, error) {
	db.ensureDB()
	dbStructure, err := db.loadDB()
	if err != nil {
		log.Printf("Error loading DB during Chirp creation: %s", err)
		return Chirp{}, err
	}

	newID := len(dbStructure.Chirps) + 1

	newChirp := Chirp{
		Id:   newID,
		Body: body,
	}

	dbStructure.Chirps[newID] = newChirp
	err = db.writeDB(dbStructure)
	if err != nil {
		log.Printf("Error writing to DB during Chirp creation: %s", err)
		return Chirp{}, err
	}

	return newChirp, nil
}

func (db *DB) GetChirps() ([]Chirp, error) {
	data, err := db.loadDB()
	if err != nil {
		log.Printf("Error loading chirps from database: %s", err)
		return []Chirp{}, err
	}
	chirps := data.Chirps
	chirpArray := []Chirp{}

	for _, chirp := range chirps {
		chirpArray = append(chirpArray, chirp)
	}

	return chirpArray, nil
}

func (db *DB) ensureDB() error {
	_, err := os.Stat(db.path)
	if os.IsNotExist(err) {
		file, err := os.Create(db.path)
		if err != nil {
			log.Printf("Error creating file: %s", err)
			return err
		}
		file.Close()
	} else if err != nil {
		log.Printf("Error checking file: %s", err)
		return err
	}
	return nil
}

func (db *DB) loadDB() (DBStructure, error) {
	db.mu.Lock()
	defer db.mu.Unlock()

	bytes, err := os.ReadFile(db.path)
	if err != nil {
		log.Printf("Error reading file: %s", err)
		return DBStructure{}, err
	}

	if len(bytes) == 0 {
		return DBStructure{
			Chirps: map[int]Chirp{},
		}, nil
	}

	var data DBStructure
	err = json.Unmarshal(bytes, &data)
	if err != nil {
		log.Printf("Error unmarshalling JSON: %s", err)
		return DBStructure{}, err
	}

	return data, nil
}

func (db *DB) writeDB(dbStructure DBStructure) error {
	dat, err := json.Marshal(dbStructure)
	if err != nil {
		log.Printf("Error marshalling JSON: %s", err)
		return err
	}

	db.mu.Lock()
	defer db.mu.Unlock()

	err = os.WriteFile(db.path, dat, os.FileMode(0644))
	if err != nil {
		log.Printf("Error writing to file: %s", err)
		return err
	}

	return nil
}
