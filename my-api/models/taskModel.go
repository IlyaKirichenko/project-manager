package models

import "database/sql"

type Task struct {
	Id_task     int    `json:"id_task"`
	Id_project  int    `json:"id_project"`
	Id_user     int    `json:"id_user"`
	Name_task   string `json:"name_task"`
	Id_state    int    `json:"id_state"`
	Date        string `json:"date"`
	Description string `json:"description"`
}

func CreateTask(idProject, idUser, idState int, nameTask, description string) (Task, error) {
	insertDb, err := db.Exec("INSERT INTO tasks (id_project, id_user, name_task, id_state, description) VALUES (?, ?, ?, ?, ?)", idProject, idUser, nameTask, idState, description)
	if err != nil {
		return Task{}, err
	}

	newTaskId, err := insertDb.LastInsertId()
	if err != nil {
		return Task{}, err
	}

	newTask := Task{
		Id_task:     int(newTaskId),
		Id_project:  idProject,
		Id_user:     idUser,
		Name_task:   nameTask,
		Id_state:    idState,
		Description: description,
	}

	err = db.QueryRow("SELECT date FROM tasks WHERE id_task = ?", newTaskId).Scan(&newTask.Date)

	return newTask, nil
}

func GetTasks() ([]Task, error) {
	tasks := []Task{}

	rows, err := db.Query("SELECT id_task, id_project, id_user, name_task, id_state, date, description FROM tasks")
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		var t Task
		if err := rows.Scan(&t.Id_task, &t.Id_project, &t.Id_user, &t.Name_task, &t.Id_state, &t.Date, &t.Description); err != nil {
			return nil, err
		}
		tasks = append(tasks, t)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return tasks, nil
}

func GetTaskByID(id int) (Task, error) {
	var t Task
	err := db.QueryRow("SELECT id_task, id_project, id_user, name_task, id_state, date, description FROM tasks WHERE id_task = ?", id).Scan(&t.Id_task, &t.Id_project, &t.Id_user, &t.Name_task, &t.Id_state, &t.Date, &t.Description)
	if err == sql.ErrNoRows {
		return Task{}, err
	}
	if err != nil {
		return Task{}, err
	}

	return t, nil
}
