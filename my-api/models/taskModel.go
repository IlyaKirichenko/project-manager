package models

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
