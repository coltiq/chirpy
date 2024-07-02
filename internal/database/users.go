package database

import "errors"

type User struct {
	ID             int    `json:"id"`
	Email          string `json:"email"`
	HashedPassword string `json:"hashed_password"`
}

var ErrAlreadyExists = errors.New("already exists")

func (db *DB) CreateUser(email, password string) (User, error) {
	dbStructure, err := db.loadDB()
	if err != nil {
		return User{}, err
	}

	for _, user := range dbStructure.Users {
		if user.Email == email {
			return User{}, errors.New("email already in database")
		}
	}

	newID := len(dbStructure.Users) + 1
	newUser := User{
		ID:             newID,
		Email:          email,
		HashedPassword: string(password),
	}
	dbStructure.Users[newID] = newUser

	err = db.writeDB(dbStructure)
	if err != nil {
		return User{}, err
	}

	return newUser, nil
}

func (db *DB) GetUsers() ([]User, error) {
	dbStructure, err := db.loadDB()
	if err != nil {
		return nil, err
	}

	users := make([]User, 0, len(dbStructure.Users))
	for _, user := range dbStructure.Users {
		users = append(users, user)
	}

	return users, nil
}

func (db *DB) GetUserByEmail(email string) (User, error) {

}
