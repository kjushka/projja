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

	lastInsertId, _ := result.LastInsertId()
	_, err = c.DB.Exec(
		"insert into task_status (status, level, project) values (?, ?, ?)",
		"new",
		0,
		lastInsertId,
	)
	if err != nil {
		log.Println("error in creating task status 'new':", err)
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

func (c *Controller) GetProjectMembers(params martini.Params, w http.ResponseWriter) (int, string) {
	projectId, err := strconv.ParseInt(params["id"], 10, 64)
	if err != nil {
		log.Println("error in parsing projectId", err)
		return 500, err.Error()
	}

	rows, err := c.DB.Query(
		"select u.id, u.name, u.username, u.telegram_id from member m "+
			"right join project p on p.id = m.project "+
			"inner join users u on u.id = m.users where p.id = ?;",
		projectId,
	)
	if err != nil {
		log.Println("error in getting members:", err)
		return 500, err.Error()
	}

	members := []*model.User{}
	for rows.Next() {
		member := &model.User{}
		err = rows.Scan(&member.Id, &member.Name, &member.Username, &member.TelegramId)
		if err != nil {
			log.Println("error in scanning members:", err)
			return 500, err.Error()
		}
		members = append(members, member)
	}

	w.Header().Set("Content-Type", "application/json")
	return c.makeContentResponse(200, "Current project members", members)
}

func (c *Controller) AddMemberToProject(params martini.Params, w http.ResponseWriter) (int, string) {
	projectId, err := strconv.ParseInt(params["id"], 10, 64)
	if err != nil {
		log.Println("error in parsing projectId", err)
		return 500, err.Error()
	}

	memberUsername := params["uname"]

	result, err := c.DB.Exec(
		"insert into member (project, users) values (?, (select id from users where username = ?))",
		projectId,
		memberUsername,
	)
	if err != nil {
		log.Println("error in adding member:", err)
		return 500, err.Error()
	}

	rowsAffected, _ := result.RowsAffected()

	w.Header().Set("Content-Type", "application/json")
	return c.makeContentResponse(202, "Member added", struct {
		Name    string
		Content interface{}
	}{
		Name:    "Rows affected",
		Content: rowsAffected,
	})
}

func (c *Controller) RemoveMemberFromProject(params martini.Params, w http.ResponseWriter) (int, string) {
	projectId, err := strconv.ParseInt(params["id"], 10, 64)
	if err != nil {
		log.Println("error in parsing projectId", err)
		return 500, err.Error()
	}

	memberUsername := params["uname"]

	result, err := c.DB.Exec(
		"delete from member where project = ? and users = (select id from users where username = ?)",
		projectId,
		memberUsername,
	)
	if err != nil {
		log.Println("error in deleting member:", err)
		return 500, err.Error()
	}

	rowsAffected, _ := result.RowsAffected()

	w.Header().Set("Content-Type", "application/json")
	return c.makeContentResponse(202, "Member deleted", struct {
		Name    string
		Content interface{}
	}{
		Name:    "Rows affected",
		Content: rowsAffected,
	})
}

func (c *Controller) CreateProjectTaskStatus(params martini.Params, w http.ResponseWriter, r *http.Request) (int, string) {
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
	taskStatus := &model.TaskStatus{}
	err = json.Unmarshal(jsonStatus, taskStatus)
	if err != nil {
		log.Println("error in unmarshalling task status")
		return 500, err.Error()
	}

	result, err := c.DB.Exec(
		"update task_status set level = level + 1 where level >= ? and project = ?",
		taskStatus.Level,
		projectId,
	)
	if err != nil {
		log.Println("error in updating status:", err)
		return 500, err.Error()
	}
	result, err = c.DB.Exec(
		"insert into task_status (status, level, project) values (?, ?, ?)",
		taskStatus.Status,
		taskStatus.Level,
		projectId,
	)
	if err != nil {
		log.Println("error in creating task status:", err)
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

func (c *Controller) RemoveStatusFromProject(params martini.Params, w http.ResponseWriter, r *http.Request) (int, string) {
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
	taskStatus := &model.TaskStatus{}
	err = json.Unmarshal(jsonStatus, taskStatus)
	if err != nil {
		log.Println("error in unmarshalling task status")
		return 500, err.Error()
	}

	result, err := c.DB.Exec(
		"delete from task_status where project = ? and status = ? and level = ?",
		projectId,
		taskStatus.Status,
		taskStatus.Level,
	)
	if err != nil {
		log.Println("error in removing task_status:", err)
		return 500, err.Error()
	}

	_, err = c.DB.Exec(
		"update task_status set level = level-1 "+
			"where level > ? and project = ? and status != ?",
		taskStatus.Level,
		projectId,
		taskStatus.Status,
	)
	if err != nil {
		log.Println("error in updating task statuses:", err)
		return 500, err.Error()
	}

	rowsAffected, _ := result.RowsAffected()

	w.Header().Set("Content-Type", "application/json")
	return c.makeContentResponse(202, "Task status deleted", struct {
		Name    string
		Content interface{}
	}{
		Name:    "Rows affected",
		Content: rowsAffected,
	})
}

func (c *Controller) GetProjectStatuses(params martini.Params, w http.ResponseWriter) (int, string) {
	projectId, err := strconv.ParseInt(params["id"], 10, 64)
	if err != nil {
		log.Println("error in parsing projectId", err)
		return 500, err.Error()
	}

	rows, err := c.DB.Query(
		"select status, level from task_status "+
			"where project = ? order by level",
		projectId,
	)
	if err != nil {
		log.Println("error in getting members:", err)
		return 500, err.Error()
	}

	taskStatuses := []*model.TaskStatus{}
	for rows.Next() {
		taskStatus := &model.TaskStatus{}
		err = rows.Scan(&taskStatus.Status, &taskStatus.Level)
		if err != nil {
			log.Println("error in scanning task status:", err)
			return 500, err.Error()
		}
		taskStatuses = append(taskStatuses, taskStatus)
	}

	w.Header().Set("Content-Type", "application/json")
	return c.makeContentResponse(200, "Current project task statuses", taskStatuses)
}

/*func (c *Controller) CreateTask(params martini.Params, w http.ResponseWriter, r *http.Request) (int, string) {
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
}*/
