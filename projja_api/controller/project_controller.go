package controller

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"projja_api/model"
	"strconv"
	"time"

	"github.com/go-martini/martini"
)

func (c *Controller) CreateProject(w http.ResponseWriter, r *http.Request) (int, string) {
	contentType := r.Header.Get("Content-Type")
	if contentType != "application/json" {
		err := fmt.Sprintf("Unsupportable Content-Type header: %s", contentType)
		log.Println(err)
		return 500, err
	}
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
		"insert into task_status (status, status_level, project) values (?, ?, ?)",
		"new",
		1,
		lastInsertId,
	)
	if err != nil {
		log.Println("error in creating task status 'new':", err)
		return 500, err.Error()
	}

	_, err = c.DB.Exec(
		"insert into member (project, users) values (?, (select id from users where username = ?))",
		lastInsertId,
		project.Owner.Username,
	)
	if err != nil {
		log.Println("error in adding owner as project member:", err)
		return 500, err.Error()
	}

	row := c.DB.QueryRow("select id, name, telegram_id, chat_id from users where username = ?", project.Owner.Username)
	var ownerId int64
	var ownerName string
	var ownerTelegramId string
	var ownerChatId int64
	err = row.Scan(&ownerId, &ownerName, &ownerTelegramId, &ownerChatId)
	if err != nil {
		log.Println("error in getting owner id: ", err)
		return 500, err.Error()
	}

	skills, err := c.getSkillsByUser(project.Owner.Username)
	if err != nil {
		log.Println("error in getting owner skills: ", err)
		return 500, err.Error()
	}

	project.Id = lastInsertId
	project.Owner.Id = ownerId
	project.Owner.Name = ownerName
	project.Owner.TelegramId = ownerTelegramId
	project.Owner.ChatId = ownerChatId
	project.Owner.Skills = skills
	_, err = c.sendDataToStream("project", "new", project)
	if err != nil {
		log.Println(err)
		return 500, err.Error()
	}

	rowsAffected, _ := result.RowsAffected()

	w.Header().Set("Content-Type", "application/json")
	return c.makeContentResponse(201, "Project created", struct {
		Name    string
		Content interface{}
	}{
		Name:    "Rows affected",
		Content: rowsAffected,
	})
}

func (c *Controller) ChangeProjectName(params martini.Params, w http.ResponseWriter, r *http.Request) (int, string) {
	contentType := r.Header.Get("Content-Type")
	if contentType != "application/json" {
		err := fmt.Sprintf("Unsupportable Content-Type header: %s", contentType)
		log.Println(err)
		return 500, err
	}
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

	row := c.DB.QueryRow(
		"select count(name) from project where id = ?",
		projectId,
	)
	var count int
	err = row.Scan(&count)
	if err != nil && err != sql.ErrNoRows {
		log.Println("error in checking project being: ", err)
		return 500, err.Error()
	}

	if count == 0 {
		err = fmt.Errorf("no such project with id %v", projectId)
		log.Println(err)
		return 404, err.Error()
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

func (c *Controller) CloseProject(params martini.Params, w http.ResponseWriter) (int, string) {
	projectId, err := strconv.ParseInt(params["id"], 10, 64)
	if err != nil {
		log.Println("error in parsing projectId", err)
		return 500, err.Error()
	}

	rowsAffected, err := c.changeProjectStatus(projectId, "closed")
	if err != nil {
		return 500, err.Error()
	}

	w.Header().Set("Content-Type", "application/json")
	return c.makeContentResponse(200, "Project closed", struct {
		Name    string
		Content interface{}
	}{
		Name:    "Rows affected",
		Content: rowsAffected,
	})
}

func (c *Controller) OpenProject(params martini.Params, w http.ResponseWriter) (int, string) {
	projectId, err := strconv.ParseInt(params["id"], 10, 64)
	if err != nil {
		log.Println("error in parsing projectId", err)
		return 500, err.Error()
	}

	rowsAffected, err := c.changeProjectStatus(projectId, "opened")
	if err != nil {
		return 500, err.Error()
	}

	w.Header().Set("Content-Type", "application/json")
	return c.makeContentResponse(200, "Project closed", struct {
		Name    string
		Content interface{}
	}{
		Name:    "Rows affected",
		Content: rowsAffected,
	})
}

func (c *Controller) changeProjectStatus(id int64, status string) (int64, error) {
	result, err := c.DB.Exec(
		"update project set status = ? where id = ?",
		status,
		id,
	)
	if err != nil {
		log.Println("error in updating status:", err)
		return 0, err
	}

	rowsAffected, _ := result.RowsAffected()

	return rowsAffected, err
}

func (c *Controller) GetProjectMembers(params martini.Params, w http.ResponseWriter) (int, string) {
	projectId, err := strconv.ParseInt(params["id"], 10, 64)
	if err != nil {
		log.Println("error in parsing projectId", err)
		return 500, err.Error()
	}

	rows, err := c.DB.Query(
		"select u.id, u.name, u.username, u.telegram_id, u.chat_id from member m "+
			"right join project p on p.id = m.project "+
			"inner join users u on u.id = m.users where p.id = ? order by u.name asc",
		projectId,
	)
	if err != nil {
		log.Println("error in getting members:", err)
		return 500, err.Error()
	}

	members := make([]*model.User, 0)
	for rows.Next() {
		member := &model.User{}
		err = rows.Scan(&member.Id, &member.Name, &member.Username, &member.TelegramId, &member.ChatId)
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

	row := c.DB.QueryRow(
		"select count(id) from users where username = ?",
		memberUsername,
	)
	var count int
	err = row.Scan(&count)
	if err != nil && err != sql.ErrNoRows {
		log.Println("error in checking user being: ", err)
		return 500, err.Error()
	}

	if count == 0 {
		err = fmt.Errorf("no such user with username %v", memberUsername)
		log.Println(err)
		return 404, err.Error()
	}

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

	row = c.DB.QueryRow("select id, name, username, telegram_id, chat_id from users where username = ?", memberUsername)
	member := &model.User{}
	err = row.Scan(&member.Id, &member.Name, &member.Username, &member.TelegramId, &member.ChatId)
	if err != nil {
		log.Println("error in getting new member info: ", err)
		return 500, err.Error()
	}
	skills, err := c.getSkillsByUser(memberUsername)
	if err != nil {
		log.Println("error in getting new member skills: ", err)
		return 500, err.Error()
	}
	member.Skills = skills

	_, err = c.sendDataToStream("project", "add-member", struct {
		ProjectId int64
		Member    *model.User
	}{projectId,
		member,
	})
	if err != nil {
		log.Println(err)
		return 500, err.Error()
	}

	w.Header().Set("Content-Type", "application/json")
	return c.makeContentResponse(201, "Member added", struct {
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

	_, err = c.DB.Exec(
		"delete from task where executor = (select id from users where username = ?)",
		memberUsername,
	)
	if err != nil {
		log.Println("error in deleting member:", err)
		return 500, err.Error()
	}

	_, err = c.sendDataToStream("project", "remove-member", struct {
		ProjectId      int64
		MemberUsername string
	}{
		projectId,
		memberUsername,
	})
	if err != nil {
		log.Println(err)
		return 500, err.Error()
	}

	rowsAffected, _ := result.RowsAffected()

	w.Header().Set("Content-Type", "application/json")
	return c.makeContentResponse(200, "Member deleted", struct {
		Name    string
		Content interface{}
	}{
		Name:    "Rows affected",
		Content: rowsAffected,
	})
}

func (c *Controller) CreateProjectTaskStatus(params martini.Params, w http.ResponseWriter, r *http.Request) (int, string) {
	contentType := r.Header.Get("Content-Type")
	if contentType != "application/json" {
		err := fmt.Sprintf("Unsupportable Content-Type header: %s", contentType)
		log.Println(err)
		return 500, err
	}
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

	_, err = c.DB.Exec(
		"update task_status set status_level = status_level + 1 where status_level >= ? and project = ?",
		taskStatus.Level,
		projectId,
	)
	if err != nil {
		log.Println("error in updating status:", err)
		return 500, err.Error()
	}
	result, err := c.DB.Exec(
		"insert into task_status (status, status_level, project) values (?, ?, ?)",
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
	return c.makeContentResponse(201, "Project status updated", struct {
		Name    string
		Content interface{}
	}{
		Name:    "Rows affected",
		Content: rowsAffected,
	})
}

func (c *Controller) RemoveStatusFromProject(params martini.Params, w http.ResponseWriter, r *http.Request) (int, string) {
	contentType := r.Header.Get("Content-Type")
	if contentType != "application/json" {
		err := fmt.Sprintf("Unsupportable Content-Type header: %s", contentType)
		log.Println(err)
		return 500, err
	}
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
		"delete from task_status where project = ? and status = ? and status_level = ?",
		projectId,
		taskStatus.Status,
		taskStatus.Level,
	)
	if err != nil {
		log.Println("error in removing task_status:", err)
		return 500, err.Error()
	}

	_, err = c.DB.Exec(
		"update task_status set status_level = status_level-1 "+
			"where status_level > ? and project = ? and status != ?",
		taskStatus.Level,
		projectId,
		taskStatus.Status,
	)
	if err != nil {
		log.Println("error in updating task_statuses:", err)
		return 500, err.Error()
	}

	row := c.DB.QueryRow(
		"select count(id) from task_status where status_level = ? and project = ?",
		taskStatus.Level,
		projectId,
	)
	var countOfStatuses int
	err = row.Scan(&countOfStatuses)

	if err != nil {
		log.Println("error in getting new task statuses count:", err)
		return 500, err.Error()
	}

	var newLevel int
	if countOfStatuses == 0 {
		newLevel = taskStatus.Level - 1
	} else {
		newLevel = taskStatus.Level
	}

	_, err = c.DB.Exec(
		"update task set status = (select id from task_status where status_level = ? and project = ?) "+
			"where status IS NULL and project = ?",
		newLevel,
		projectId,
		projectId,
	)
	if err != nil {
		log.Println("error in updating task statuses:", err)
		return 500, err.Error()
	}

	rowsAffected, _ := result.RowsAffected()

	w.Header().Set("Content-Type", "application/json")
	return c.makeContentResponse(200, "Task status deleted", struct {
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
		"select status, status_level from task_status "+
			"where project = ? order by status_level",
		projectId,
	)
	if err != nil {
		log.Println("error in getting members:", err)
		return 500, err.Error()
	}

	taskStatuses := make([]*model.TaskStatus, 0)
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
	contentType := r.Header.Get("Content-Type")
	if contentType != "application/json" {
		err := fmt.Sprintf("Unsupportable Content-Type header: %s", contentType)
		log.Println(err)
		return 500, err
	}
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
			"values (?, ?, ?, ?, (select id from task_status where status_level = ? and project = ?), "+
			"?, (select id from users where username = ?))",
		task.Description,
		projectId,
		deadline,
		task.Priority,
		1,
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
			log.Println("error in adding skills: ", err)
		}
	}

	row := c.DB.QueryRow("select id from users where username = ?", task.Executor.Username)
	var executorId int64
	err = row.Scan(&executorId)
	if err != nil {
		log.Println("error in getting executor id: ", err)
		return 500, err.Error()
	}

	task.Id = lastInsertId
	task.Executor.Id = executorId
	_, err = c.sendDataToStream("project", "task", struct {
		ProjectId int64
		Task      *model.Task
	}{
		projectId,
		task,
	})

	if err != nil {
		log.Println(err)
		return 500, err.Error()
	}

	rowsAffected, _ := result.RowsAffected()

	w.Header().Set("Content-Type", "application/json")
	return c.makeContentResponse(201, "New task created successfully", struct {
		Name    string
		Content interface{}
	}{
		Name:    "Rows affected",
		Content: rowsAffected,
	})
}

func (c *Controller) GetAllProjectTasks(params martini.Params, w http.ResponseWriter) (int, string) {
	projectId, err := strconv.ParseInt(params["id"], 10, 64)
	if err != nil {
		log.Println("error in parsing projectId", err)
		return 500, err.Error()
	}

	rows, err := c.DB.Query(
		"select t.id, t.description, p.id, p.name, p.ow_id, p.ow_name, p.ow_username, "+
			"p.ow_telegram_id, p.ow_chat_id, p.status, t.deadline, t.priority, ts.status, ts.status_level, "+
			"e.id, e.name, e.username, e.telegram_id, e.chat_id from task t "+
			"left join (select p.id, p.name, u.id ow_id, u.name ow_name, "+
			"u.username ow_username, u.telegram_id ow_telegram_id, u.chat_id ow_chat_id, p.status "+
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
			&task.Project.Owner.ChatId,
			&task.Project.Status,
			&deadline,
			&task.Priority,
			&task.Status.Status,
			&task.Status.Level,
			&task.Executor.Id,
			&task.Executor.Name,
			&task.Executor.Username,
			&task.Executor.TelegramId,
			&task.Executor.ChatId,
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

func (c *Controller) GetProcessProjectTasks(params martini.Params, w http.ResponseWriter) (int, string) {
	projectId, err := strconv.ParseInt(params["id"], 10, 64)
	if err != nil {
		log.Println("error in parsing projectId", err)
		return 500, err.Error()
	}

	rows, err := c.DB.Query(
		"select t.id, t.description, p.id, p.name, p.ow_id, p.ow_name, p.ow_username, "+
			"p.ow_telegram_id, p.ow_chat_id, p.status, t.deadline, t.priority, ts.status, ts.status_level, "+
			"e.id, e.name, e.username, e.telegram_id, e.chat_id from task t "+
			"left join (select p.id, p.name, u.id ow_id, u.name ow_name, "+
			"u.username ow_username, u.telegram_id ow_telegram_id, u.chat_id ow_chat_id, p.status "+
			"from project p left join users u on u.id = p.owner) p on p.id = t.project "+
			"left join task_status ts on ts.id = t.status "+
			"left join users e on t.executor = e.id "+
			"where t.project = ? and t.is_closed = 0 and t.is_closed <> true",
		projectId,
	)
	if err != nil && err != sql.ErrNoRows {
		log.Println("error in getting tasks:", err)
		return 500, err.Error()
	}

	tasks := make([]*model.Task, 0)

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
			&task.Project.Owner.ChatId,
			&task.Project.Status,
			&deadline,
			&task.Priority,
			&task.Status.Status,
			&task.Status.Level,
			&task.Executor.Id,
			&task.Executor.Name,
			&task.Executor.Username,
			&task.Executor.TelegramId,
			&task.Executor.ChatId,
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
