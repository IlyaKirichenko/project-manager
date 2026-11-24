package models

import "database/sql"

type Project struct {
	Id_project   int    `json:"id_project"`
	Id_owner     int    `json:"id_owner"`
	Name_project string `json:"name_project"`
	Description  string `json:"description_project"`
	Deadline     string `json:"deadline"`
}

func CreateProject(IdOwner int, nameProject, description, deadline string) (Project, error) {
	tx, err := db.Begin()
	if err != nil {
		return Project{}, err
	}
	defer tx.Rollback()

	insertDb, err := tx.Exec("INSERT INTO project (id_owner, name_project, description_project, deadline) VALUES (?, ?, ?, ?)", IdOwner, nameProject, description, deadline)
	if err != nil {
		return Project{}, err
	}

	newProjectId, err := insertDb.LastInsertId()
	if err != nil {
		return Project{}, err
	}

	_, err = tx.Exec("INSERT INTO users_in_project (id_project, id_users) VALUES (?, ?)", newProjectId, IdOwner)
	if err != nil {
		return Project{}, err
	}

	newProject := Project{
		Id_project:   int(newProjectId),
		Id_owner:     IdOwner,
		Name_project: nameProject,
		Description:  description,
		Deadline:     deadline,
	}

	if err := tx.Commit(); err != nil {
		return Project{}, err
	}

	return newProject, nil
}

func GetProjectsByUserID(userID int) ([]Project, error) {
	projects := []Project{}

	rows, err := db.Query("SELECT p.id_project, p.id_owner, p.name_project, p.description_project, p.deadline FROM project p JOIN users_in_project uip ON p.id_project = uip.id_project WHERE uip.id_users = ?", userID)
	if err == sql.ErrNoRows {
		return nil, err
	}
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		var p Project
		if err := rows.Scan(&p.Id_project, &p.Id_owner, &p.Name_project, &p.Description, &p.Deadline); err != nil {
			return nil, err
		}
		projects = append(projects, p)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return projects, nil
}
