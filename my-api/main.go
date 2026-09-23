package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"

	"strconv"

	"github.com/gin-gonic/gin"
	_ "github.com/microsoft/go-mssqldb"
)

type Task struct {
	ID    int    `json:"id"`
	Title string `json:"title"`
	Done  bool   `json:"done"`
}

var db *sql.DB

func initDB() {
	connString := "server=.\\SQLEXPRESS;database=CourseDB;Integrated Security=true;encrypt=disable"

	var err error
	db, err = sql.Open("sqlserver", connString)
	if err != nil {
		log.Fatalf("Ошибка открытия БД: %v", err)
	}
	if err := db.Ping(); err != nil {
		log.Fatalf("Ошибка подключения к БД: %v", err)
	}

	fmt.Println("Подключено к SQL Server!")
}

func main() {

	initDB()

	router := gin.Default()

	router.POST("/tasks", func(c *gin.Context) {
		var newTask Task

		if err := c.ShouldBindJSON(&newTask); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		query := "INSERT INTO tasks (title, done) OUTPUT INSERTED.id VALUES (@p1, @p2)"

		var newID int

		err := db.QueryRow(query, newTask.Title, newTask.Done).Scan(&newID)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка записи в БД: " + err.Error()})
			return
		}

		newTask.ID = newID

		c.JSON(http.StatusCreated, newTask)
	})

	router.GET("/tasks", func(c *gin.Context) {
		doneParam := c.Query("done")
		searchParam := c.Query("search")

		query := "SELECT id, title, done FROM tasks WHERE 1=1"

		var args []interface{}

		paramIndex := 1

		if doneParam != "" {
			isDone, err := strconv.ParseBool(doneParam)
			if err == nil {
				query += fmt.Sprintf(" AND done = @p%d", paramIndex)
				args = append(args, isDone)
				paramIndex++
			}
		}
		if searchParam != "" {
			query += fmt.Sprintf(" AND title LIKE '%%' + @p%d + '%%'", paramIndex)
			args = append(args, searchParam)
			paramIndex++
		}

		rows, err := db.Query(query, args...)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка чтения из БД"})
			return
		}

		defer rows.Close()

		var tasks []Task

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

	router.Run(":8000")
}
