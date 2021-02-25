package controller

import (
	"encoding/json"
	"github.com/go-martini/martini"
	"io/ioutil"
	"log"
	"net/http"
	"projja_api/model"
	"strconv"
)

func (c *Controller) CreateProject(w http.ResponseWriter, r *http.Request) (int, string) {
	jsonProject, err := ioutil.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		log.Println("error in reading body", err)
		return 500, err.Error()
	}
	project := &model.Project{}
	err = json.Unmarshal(jsonProject, project)
	if err != nil {
		log.Println("error in unmarshalling")
		return 500, err.Error()
	}

	result, err := c.DB.Exec(
		"insert into project (name, owner, status) values (?, (select id from users where username = ?), ?)",
		project.Name,
		project.Owner.Username,
		"opened",
	)
	if err != nil {
		log.Println("error in creating project:", err)
		return 500, err.Error()
	}
	rowsAffected, _ := result.RowsAffected()

	w.Header().Set("Content-Type", "application/json")
	return c.makeContentResponse(202, "Project created", struct {
		Name    string
		Content interface{}
	}{
		Name:    "Rows affected",
		Content: rowsAffected,
	})
}

func (c *Controller) ChangeProjectName(params martini.Params, w http.ResponseWriter, r *http.Request) (int, string) {
	projectId, err := strconv.ParseInt(params["id"], 10, 64)
	if err != nil {
		log.Println("error in parsing projectId", err)
		return 500, err.Error()
	}

	jsonProjectName, err := ioutil.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		log.Println("error in reading body", err)
		return 500, err.Error()
	}
	projectName := &struct {
		Name string
	}{}
	err = json.Unmarshal(jsonProjectName, projectName)
	if err != nil {
		log.Println("error in unmarshalling")
		return 500, err.Error()
	}

	result, err := c.DB.Exec(
		"update project set name = ? where id = ?",
		projectName.Name,
		projectId,
	)
	if err != nil {
		log.Println("error in updating name:", err)
		return 500, err.Error()
	}

	rowsAffected, _ := result.RowsAffected()

	w.Header().Set("Content-Type", "application/json")
	return c.makeContentResponse(200, "Project name updated", struct {
		Name    string
		Content interface{}
	}{
		Name:    "Rows affected",
		Content: rowsAffected,
	})
}

func (c *Controller) ChangeProjectStatus(params martini.Params, w http.ResponseWriter, r *http.Request) (int, string) {
	projectId, err := strconv.ParseInt(params["id"], 10, 64)
	if err != nil {
		log.Println("error in parsing projectId", err)
		return 500, err.Error()
	}

	jsonStatus, err := ioutil.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		log.Println("error in reading body", err)
		return 500, err.Error()
	}
	projectStatus := &struct {
		Status string
	}{}
	err = json.Unmarshal(jsonStatus, projectStatus)
	if err != nil {
		log.Println("error in unmarshalling")
		return 500, err.Error()
	}

	result, err := c.DB.Exec(
		"update project set status = ? where id = ?",
		projectStatus.Status,
		projectId,
	)
	if err != nil {
		log.Println("error in updating status:", err)
		return 500, err.Error()
	}

	rowsAffected, _ := result.RowsAffected()

	w.Header().Set("Content-Type", "application/json")
	return c.makeContentResponse(200, "Project status updated", struct {
		Name    string
		Content interface{}
	}{
		Name:    "Rows affected",
		Content: rowsAffected,
	})
}
