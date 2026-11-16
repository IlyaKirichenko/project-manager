package main

import (
	"my-api/models"

	"my-api/controllers"

	"github.com/gin-gonic/gin"
)

func main() {

	models.InitDB()
	models.InitSchema()
	models.SeedStates()

	router := gin.Default()

	router.POST("/user", controllers.CreateUserHandler)

	router.POST("/project", controllers.CreateProjectHandler)

	router.POST("/tasks", controllers.CreateTaskHandler)
	router.GET("/tasks/:id", controllers.GetTaskByIDHandler)

	router.Run(":8000")

	/*router.POST("/tasks", func(c *gin.Context) {
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
	})*/
}
