package models

import (
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	Id_user    int    `json:"id_user"`
	User_name  string `json:"user_name"`
	Email      string `json:"email"`
	Password   string `json:"-"`
	Admin      bool   `json:"admin"`
	Is_blocked bool   `json:"is_blocked"`
}

func CreateUser(userName, email, password string) (User, error) {

	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

	if err != nil {
		return User{}, err
	}

	passwordHash := string(hash)

	insertDB, err := db.Exec("INSERT INTO user (user_name, email, password) VALUES (?, ?, ?)", userName, email, passwordHash)

	if err != nil {
		return User{}, err
	}

	newUserId, err := insertDB.LastInsertId()

	if err != nil {
		return User{}, err
	}

	newUser := User{
		Id_user:   int(newUserId),
		User_name: userName,
		Email:     email,
		Password:  passwordHash,
	}

	return newUser, nil
}
