package models

import (
	"database/sql"

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

func GetUserByEmail(email string) (User, error) {
	var u User

	err := db.QueryRow("SELECT id_user, user_name, email, password, admin, is_blocked FROM user WHERE email = ?", email).Scan(&u.Id_user, &u.User_name, &u.Email, &u.Password, &u.Admin, &u.Is_blocked)

	if err == sql.ErrNoRows {
		return User{}, err
	}
	if err != nil {
		return User{}, err
	}

	return u, nil
}

func GetUserByID(id int) (User, error) {
	var u User

	err := db.QueryRow("SELECT id_user, user_name, email, password, admin, is_blocked FROM user WHERE id_user = ?", id).Scan(&u.Id_user, &u.User_name, &u.Email, &u.Password, &u.Admin, &u.Is_blocked)

	if err == sql.ErrNoRows {
		return User{}, err
	}
	if err != nil {
		return User{}, err
	}

	return u, nil
}

func CheckPassword(user User, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
}
