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

	db, err = sql.Open("sqlite", "database.db")

	if err != nil {
		log.Fatalf("Ошибка открытия БД: %v", err)
	}
	if err := db.Ping(); err != nil {
		log.Fatalf("Ошибка подключения к БД: %v", err)
	}

	fmt.Println("Подключено к SQLite!")
}

func initSchema() {
	query := `CREATE TABLE IF NOT EXISTS tasks (
	id		INTEGER PRIMARY KEY AUTOINCREMENT,
	title	TEXT NOT NULL,
	done 	BOOLEAN NOT NULL DEFAULT 0
	 );`

	_, err := db.Exec(query)

	if err != nil {
		log.Fatalf("Ошибка создания таблицы: %v", err)
	}

	fmt.Println("Таблица tasks готова")
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

	router.Run(":8000")
}
