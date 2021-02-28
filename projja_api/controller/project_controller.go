package controller

import (
	"database/sql"
	"encoding/json"
	"github.com/go-martini/martini"
	"io/ioutil"
	"log"
	"net/http"
	"projja_api/model"
	"strconv"
	"time"
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

func (c *Controller) CreateTask(params martini.Params, w http.ResponseWriter, r *http.Request) (int, string) {
	projectId, err := strconv.ParseInt(params["id"], 10, 64)
	if err != nil {
		log.Println("error in parsing projectId", err)
		return 500, err.Error()
	}

	jsonTask, err := ioutil.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		log.Println("error in reading body", err)
		return 500, err.Error()
	}
	task := &model.Task{}
	err = json.Unmarshal(jsonTask, task)
	if err != nil {
		log.Println("error in unmarshalling")
		return 500, err.Error()
	}

	deadline, err := time.Parse("2006-01-02", task.Deadline)
	if err != nil {
		log.Println("error in parsing time:", err)
		return 500, err.Error()
	}

	result, err := c.DB.Exec(
		"insert into task (description, project, deadline, priority, status, is_closed, executor) "+
			"values (?, ?, ?, ?, (select id from task_status where status = ? and project = ?), "+
			"?, (select id from users where username = ?))",
		task.Description,
		projectId,
		deadline,
		task.Priority,
		"new",
		projectId,
		0,
		task.Executor.Username,
	)
	if err != nil {
		log.Println("error in creating task:", err)
		return 500, err.Error()
	}
	lastInsertId, _ := result.LastInsertId()

	if len(task.Skills) != 0 {
		_, err = c.setSkillsToTask(task.Skills, lastInsertId)
		if err != nil {
			log.Println("error in adding skills")
		}
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

func (c *Controller) GetProjectTasks(params martini.Params, w http.ResponseWriter) (int, string) {
	projectId, err := strconv.ParseInt(params["id"], 10, 64)
	if err != nil {
		log.Println("error in parsing projectId", err)
		return 500, err.Error()
	}

	rows, err := c.DB.Query(
		"select t.id, t.description, p.id, p.name, p.ow_id, p.ow_name, p.ow_username, "+
			"p.ow_telegram_id, p.status, t.deadline, t.priority, ts.status, ts.level, "+
			"e.id, e.name, e.username, e.telegram_id from task t "+
			"left join (select p.id, p.name, u.id ow_id, u.name ow_name, "+
			"u.username ow_username, u.telegram_id ow_telegram_id, p.status "+
			"from project p left join users u on u.id = p.owner) p on p.id = t.project "+
			"left join task_status ts on ts.id = t.status "+
			"left join users e on t.executor = e.id "+
			"where t.project = ? and t.is_closed = 0;",
		projectId,
	)
	if err != nil && err != sql.ErrNoRows {
		log.Println("error in getting tasks:", err)
		return 500, err.Error()
	}

	tasks := []*model.Task{}

	for rows.Next() {
		task := &model.Task{}
		task.Project = &model.Project{}
		task.Project.Owner = &model.User{}
		task.Status = &model.TaskStatus{}
		task.Executor = &model.User{}
		var deadline time.Time

		err = rows.Scan(
			&task.Id,
			&task.Description,
			&task.Project.Id,
			&task.Project.Name,
			&task.Project.Owner.Id,
			&task.Project.Owner.Name,
			&task.Project.Owner.Username,
			&task.Project.Owner.TelegramId,
			&task.Project.Status,
			&deadline,
			&task.Priority,
			&task.Status.Status,
			&task.Status.Level,
			&task.Executor.Id,
			&task.Executor.Name,
			&task.Executor.Username,
			&task.Executor.TelegramId,
		)

		if err != nil {
			log.Println("error in scanning tasks:", err)
			return 500, err.Error()
		}

		task.Deadline = deadline.Format("2006-01-02")
		tasks = append(tasks, task)
	}

	w.Header().Set("Content-Type", "application/json")
	return c.makeContentResponse(200, "project tasks", tasks)
}
