package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"

	"strconv"

	"github.com/gin-gonic/gin"
	_ "modernc.org/sqlite"
)

type Task struct {
	ID    int    `json:"id"`
	Title string `json:"title"`
	Done  bool   `json:"done"`
}

var db *sql.DB

func initDB() {

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

func initSchema() {

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

func main() {

	initDB()

	initSchema()

	router := gin.Default()

	router.POST("/tasks", func(c *gin.Context) {
		var newTask Task

		if err := c.ShouldBindJSON(&newTask); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		res, err := db.Exec("INSERT INTO tasks (title, done) VALUES (?, ?)", newTask.Title, newTask.Done)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка записи в БД: " + err.Error()})
			return
		}

		newID, err := res.LastInsertId()

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка записи в БД: " + err.Error()})
			return
		}

		newTask.ID = int(newID)

		c.JSON(http.StatusCreated, newTask)
	})

	router.GET("/tasks", func(c *gin.Context) {
		doneParam := c.Query("done")
		searchParam := c.Query("search")

		query := "SELECT id, title, done FROM tasks WHERE 1=1"

		var args []interface{}

		if doneParam != "" {
			isDone, err := strconv.ParseBool(doneParam)
			if err == nil {
				query += " AND done = ?"
				args = append(args, isDone)
			}
		}
		if searchParam != "" {
			query += " AND title LIKE '%' || ? || '%'"
			args = append(args, searchParam)
		}

		rows, err := db.Query(query, args...)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка чтения из БД"})
			return
		}

		defer rows.Close()

		tasks := []Task{}

		for rows.Next() {
			var t Task
			err := rows.Scan(&t.ID, &t.Title, &t.Done)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка сканирования БД"})
				return
			}

			tasks = append(tasks, t)
		}

		if err = rows.Err(); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка при чтении БД"})
			return
		}
		c.JSON(http.StatusOK, tasks)
	})

	router.GET("/tasks/:id", func(c *gin.Context) {
		idStr := c.Param("id")

		id, err := strconv.Atoi(idStr)

		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Неверный id"})
			return
		}

		var t Task
		err = db.QueryRow("SELECT id, title, done FROM tasks WHERE id = ?", id).Scan(&t.ID, &t.Title, &t.Done)
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, gin.H{"error": "Задача не найдена"})
			return
		}
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Ошибка запроса"})
			return
		}

		c.JSON(http.StatusOK, t)
	})

	router.PUT("/tasks/:id", func(c *gin.Context) {
		idStr := c.Param("id")

		id, err := strconv.Atoi(idStr)

		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Неверный id"})
			return
		}

		var updateTask Task
		if err := c.ShouldBindJSON(&updateTask); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		res, err := db.Exec("UPDATE tasks SET title = ?, done = ? WHERE id = ?", updateTask.Title, updateTask.Done, id)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		n, err := res.RowsAffected()

		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		if n == 0 {
			c.JSON(http.StatusNotFound, gin.H{"error": "Затронуто 0 строк"})
			return
		}

		updateTask.ID = id

		c.JSON(http.StatusOK, updateTask)

	})

	router.DELETE("/tasks/:id", func(c *gin.Context) {
		idStr := c.Param("id")

		id, err := strconv.Atoi(idStr)

		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Неверный id"})
			return
		}

		res, err := db.Exec("DELETE FROM tasks WHERE id = ?", id)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		n, err := res.RowsAffected()

		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		if n == 0 {
			c.JSON(http.StatusNotFound, gin.H{"error": "Затронуто 0 строк"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "Задача удалена"})
	})

	router.Run(":8000")
}
