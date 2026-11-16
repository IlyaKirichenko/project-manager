package models

type Project struct {
	Id_project   int    `json:"id_project"`
	Id_owner     int    `json:"id_owner"`
	Name_project string `json:"name_project"`
	Description  string `json:"description_project"`
	Deadline     string `json:"deadline"`
}

func CreateProject(IdOwner int, nameProject, descroption, deadline string) (Project, error) {
	inserDb, err := db.Exec("INSERT INTO project (id_owner, name_project, description_project, deadline) VALUES (?, ?, ?, ?)", IdOwner, nameProject, descroption, deadline)
	if err != nil {
		return Project{}, err
	}

	newProjectId, err := inserDb.LastInsertId()
	if err != nil {
		return Project{}, err
	}

	newProject := Project{
		Id_project:   int(newProjectId),
		Id_owner:     IdOwner,
		Name_project: nameProject,
		Description:  descroption,
		Deadline:     deadline,
	}

	return newProject, nil
}
