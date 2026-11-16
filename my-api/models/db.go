package models

import (
	"database/sql"
	"fmt"
	"log"

	_ "modernc.org/sqlite"
)

var db *sql.DB

func InitDB() {

	var err error

	db, err = sql.Open("sqlite", "database.db?_pragma=foreign_keys(1)&_pragma=busy_timeout(2000)&_pragma=journal_mode(WAL)")

	if err != nil {
		log.Fatalf("Ошибка открытия БД: %v", err)
	}

	if err := db.Ping(); err != nil {
		log.Fatalf("Ошибка подключения к БД: %v", err)
	}

	fmt.Println("Соединение с БД установлено!")
}

func InitSchema() {

	query := `
	 CREATE TABLE IF NOT EXISTS user(
	 id_user 	INTEGER PRIMARY KEY AUTOINCREMENT,
	 user_name 	TEXT NOT NULL,
	 email 		TEXT NOT NULL UNIQUE,
	 password 	TEXT NOT NULL,
	 admin 		BOOLEAN NOT NULL DEFAULT 0,
	 is_blocked BOOLEAN NOT NULL DEFAULT 0
	 );
	 
	 CREATE TABLE IF NOT EXISTS project(
	 id_project				INTEGER PRIMARY KEY AUTOINCREMENT,
	 id_owner				INTEGER NOT NULL REFERENCES user (id_user),
	 name_project			TEXT NOT NULL,
	 description_project	TEXT NULL,
	 deadline				TEXT NOT NULL
	 );

	 CREATE TABLE IF NOT EXISTS users_in_project(
	 id_project		INTEGER NOT NULL REFERENCES project (id_project),
	 id_users		INTEGER NOT NULL REFERENCES user (id_user),
	 PRIMARY KEY (id_project, id_users)
	 );

	 CREATE TABLE IF NOT EXISTS state_tasks(
	 id_state	INTEGER PRIMARY KEY AUTOINCREMENT,
	 name		TEXT NOT NULL UNIQUE
	 );

	 CREATE TABLE IF NOT EXISTS tasks(
	 id_task 		INTEGER PRIMARY KEY AUTOINCREMENT,
	 id_project 	INTEGER NOT NULL REFERENCES project (id_project),
	 id_user 		INTEGER NOT NULL REFERENCES user (id_user),
	 name_task 		TEXT NOT NULL,
	 id_state 		INTEGER NOT NULL REFERENCES state_tasks (id_state),
	 date 			TEXT NOT NULL DEFAULT CURRENT_TIMESTAMP,
	 description 	TEXT
	 ); 
	 `

	_, err := db.Exec(query)

	if err != nil {
		log.Fatalf("Ошибка создания таблицы: %v", err)
	}

	fmt.Println("Таблицы готовы")
}

func SeedStates() {
	_, err := db.Exec("INSERT OR IGNORE INTO state_tasks (name) VALUES (?), (?)", "в работе", "выполнено")
	if err != nil {
		log.Fatalf("Ошибка создания таблицы state_tasks: %v", err)
	}
}
