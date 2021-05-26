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
	"strings"
	"time"

	"github.com/go-martini/martini"
	"github.com/scylladb/go-set"
)

func (c *Controller) GetTask(params martini.Params, w http.ResponseWriter) (int, string) {
	taskId, err := strconv.ParseInt(params["id"], 10, 64)
	if err != nil {
		log.Println("error in parsing taskId", err)
		return 500, err.Error()
	}

	task, err := c.getTaskById(taskId)
	if err != nil {
		return 500, err.Error()
	}

	w.Header().Set("Content-Type", "application/json")
	return c.makeContentResponse(200, "Current task", task)
}

func (c *Controller) ChangeTaskExecutor(params martini.Params, w http.ResponseWriter, r *http.Request) (int, string) {
	contentType := r.Header.Get("Content-Type")
	if contentType != "application/json" {
		err := fmt.Sprintf("Unsupportable Content-Type header: %s", contentType)
		log.Println(err)
		return 500, err
	}
	taskId, err := strconv.ParseInt(params["id"], 10, 64)
	if err != nil {
		log.Println("error in parsing taskId", err)
		return 500, err.Error()
	}

	jsonExecutor, err := ioutil.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		log.Println("error in reading body", err)
		return 500, err.Error()
	}
	executor := &model.User{}
	err = json.Unmarshal(jsonExecutor, executor)
	if err != nil {
		log.Println("error in unmarshalling")
		return 500, err.Error()
	}

	row := c.DB.QueryRow("select executor from task where id = ?", taskId)
	var oldUserId int64
	err = row.Scan(&oldUserId)
	if err != nil {
		log.Println("error in getting userId: ", err)
		return 500, err.Error()
	}

	result, err := c.DB.Exec(
		"update task set executor = (select id from users where username = ?) where id = ?",
		executor.Username,
		taskId,
	)
	if err != nil {
		log.Println("error in updating task executor: ", err)
		return 500, err.Error()
	}

	row = c.DB.QueryRow("select id from users where username = ?", executor.Username)
	var newUserId int64
	err = row.Scan(&newUserId)
	if err != nil {
		log.Println("error in getting userId: ", err)
		return 500, err.Error()
	}

	row = c.DB.QueryRow("select project from task where id = ?", taskId)
	var projectId int64
	err = row.Scan(&projectId)
	if err != nil {
		log.Println("error in getting projectId: ", err)
		return 500, err.Error()
	}

	_, err = c.sendDataToStream("task", "executor", struct {
		TaskId    int64
		OldUserId int64
		NewUserId int64
		ProjectId int64
	}{
		taskId,
		oldUserId,
		newUserId,
		projectId,
	})

	if err != nil {
		log.Println(err)
		return 500, err.Error()
	}

	if err != nil {
		log.Println("error in updating executor:", err)
		return 500, err.Error()
	}

	rowsAffected, _ := result.RowsAffected()

	w.Header().Set("Content-Type", "application/json")
	return c.makeContentResponse(200, "Task executor updated", struct {
		Name    string
		Content interface{}
	}{
		Name:    "Rows affected",
		Content: rowsAffected,
	})
}

func (c *Controller) ChangeTaskDescription(params martini.Params, w http.ResponseWriter, r *http.Request) (int, string) {
	contentType := r.Header.Get("Content-Type")
	if contentType != "application/json" {
		err := fmt.Sprintf("Unsupportable Content-Type header: %s", contentType)
		log.Println(err)
		return 500, err
	}
	taskId, err := strconv.ParseInt(params["id"], 10, 64)
	if err != nil {
		log.Println("error in parsing taskId", err)
		return 500, err.Error()
	}

	jsonDescription, err := ioutil.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		log.Println("error in reading body", err)
		return 500, err.Error()
	}
	description := &struct {
		Description string
	}{}
	err = json.Unmarshal(jsonDescription, description)
	if err != nil {
		log.Println("error in unmarshalling")
		return 500, err.Error()
	}

	result, err := c.DB.Exec(
		"update task set description = ? where id = ?",
		description.Description,
		taskId,
	)
	if err != nil {
		log.Println("error in updating description:", err)
		return 500, err.Error()
	}

	row := c.DB.QueryRow("select project from task where id = ?", taskId)
	var projectId int64
	err = row.Scan(&projectId)
	if err != nil {
		log.Println("error in getting projectId: ", err)
		return 500, err.Error()
	}

	_, err = c.sendDataToStream("task", "description", struct {
		TaskId      int64
		Description string
		ProjectId   int64
	}{
		taskId,
		description.Description,
		projectId,
	})
	if err != nil {
		log.Println(err)
		return 500, err.Error()
	}

	rowsAffected, _ := result.RowsAffected()

	w.Header().Set("Content-Type", "application/json")
	return c.makeContentResponse(200, "Task description updated", struct {
		Name    string
		Content interface{}
	}{
		Name:    "Rows affected",
		Content: rowsAffected,
	})
}

func (c *Controller) SetSkillsToTask(params martini.Params, w http.ResponseWriter, r *http.Request) (int, string) {
	contentType := r.Header.Get("Content-Type")
	if contentType != "application/json" {
		err := fmt.Sprintf("Unsupportable Content-Type header: %s", contentType)
		log.Println(err)
		return 500, err
	}
	taskId, err := strconv.ParseInt(params["id"], 10, 64)
	if err != nil {
		log.Println("error in parsing taskId", err)
		return 500, err.Error()
	}

	jsonSkills, err := ioutil.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		log.Println("error in reading body", err)
		return 500, err.Error()
	}
	skills := &struct {
		Skills []string
	}{}
	err = json.Unmarshal(jsonSkills, skills)
	if err != nil {
		log.Println("error in unmarshalling")
		return 500, err.Error()
	}

	rows, err := c.DB.Query(
		"select s.skill from skill s",
	)
	if err != nil && err != sql.ErrNoRows {
		log.Println("error in getting skills:", err)
		return 500, err.Error()
	}
	skillsSet := set.NewStringSet()
	for rows.Next() {
		var skill string
		err = rows.Scan(&skill)
		if err != nil {
			log.Println("error in scan skills:", err)
			return 500, err.Error()
		}
		skillsSet.Add(skill)
	}
	newSkills := make([]string, 0)
	newTaskSkills := make([]string, len(skills.Skills))
	for i, skill := range skills.Skills {
		skill := strings.ToLower(skill)
		if !skillsSet.Has(skill) {
			newSkills = append(newSkills, fmt.Sprintf("('%v')", strings.ToLower(skill)))
		}
		newTaskSkills[i] = fmt.Sprintf("(%v, (select s.id from skill s where s.skill = '%v'))",
			taskId,
			skill,
		)
	}

	if len(newSkills) != 0 {
		_, err := c.DB.Exec("insert into skill (skill) values " + strings.Join(newSkills, ", "))
		if err != nil {
			log.Println("error in creating new skills:", err)
			return 500, err.Error()
		}
	}

	_, err = c.DB.Exec(
		"delete from task_skill where task = ?",
		taskId,
	)
	if err != nil {
		log.Println("error in deleting skills:", err)
		return 500, err.Error()
	}

	result, err := c.DB.Exec(
		"insert into task_skill (task, skill) values " +
			strings.Join(newTaskSkills, ", "),
	)
	if err != nil {
		log.Println("error in creating new task_skill:", err)
		return 500, err.Error()
	}

	rowsAffected, _ := result.RowsAffected()
	w.Header().Set("Content-Type", "application/json")
	return c.makeContentResponse(201, "skills set", struct {
		Name    string
		Content interface{}
	}{
		Name:    "Rows affected",
		Content: rowsAffected,
	})
}

func (c *Controller) SetPreviousTaskStatus(params martini.Params, w http.ResponseWriter) (int, string) {
	taskId, err := strconv.ParseInt(params["id"], 10, 64)
	if err != nil {
		log.Println("error in parsing taskId", err)
		return 500, err.Error()
	}

	row := c.DB.QueryRow(
		"select count(*) from task_status ts "+
			"right join (select t.project, ts.status_level from task t "+
			"left join task_status ts on ts.id = t.status where t.id = ?) "+
			"t on t.project = ts.project where ts.status_level <= t.status_level - 1;",
		taskId,
	)
	count := 0
	err = row.Scan(&count)
	if err != nil && err != sql.ErrNoRows {
		log.Println("error in getting count of previous levels:", err)
		return 500, err.Error()
	}
	if count == 0 {
		return 500, "no such previous status"
	}
	result, err := c.DB.Exec(
		"update task set status = (select ts.id from task_status ts "+
			"right join (select t.project, ts.status_level from task t "+
			"left join task_status ts on ts.id = t.status where t.id = ?) "+
			"t on t.project = ts.project where ts.status_level = t.status_level - 1) where id = ?;",
		taskId,
		taskId,
	)
	if err != nil {
		log.Println("error in updating status:", err)
		return 500, err.Error()
	}

	rowsAffected, _ := result.RowsAffected()

	w.Header().Set("Content-Type", "application/json")
	return c.makeContentResponse(200, "Task status updated", struct {
		Name    string
		Content interface{}
	}{
		Name:    "Rows affected",
		Content: rowsAffected,
	})
}

func (c *Controller) SetNextTaskStatus(params martini.Params, w http.ResponseWriter) (int, string) {
	taskId, err := strconv.ParseInt(params["id"], 10, 64)
	if err != nil {
		log.Println("error in parsing taskId", err)
		return 500, err.Error()
	}

	row := c.DB.QueryRow(
		"select count(*) from task_status ts "+
			"right join (select t.project, ts.status_level from task t "+
			"left join task_status ts on ts.id = t.status where t.id = ?) "+
			"t on t.project = ts.project where ts.status_level >= t.status_level + 1;",
		taskId,
	)
	count := 0
	err = row.Scan(&count)
	if err != nil && err != sql.ErrNoRows {
		log.Println("error in getting count of next levels:", err)
		return 500, err.Error()
	}
	var result sql.Result
	var message string
	if count == 0 {
		result, err = c.DB.Exec(
			"update task set is_closed = ? where id = ?",
			true,
			taskId,
		)

		if err != nil {
			log.Println("error in closing task: ", err)
			return 500, err.Error()
		}

		row := c.DB.QueryRow("select executor from task where id = ?", taskId)
		var executorId int64
		err = row.Scan(&executorId)
		if err != nil {
			log.Println("error in getting userId: ", err)
			return 500, err.Error()
		}

		row = c.DB.QueryRow("select project from task where id = ?", taskId)
		var projectId int64
		err = row.Scan(&projectId)
		if err != nil {
			log.Println("error in getting projectId: ", err)
			return 500, err.Error()
		}

		_, err = c.sendDataToStream("task", "close", struct {
			TaskId     int64
			ExecutorId int64
			ProjectId  int64
		}{
			taskId,
			executorId,
			projectId,
		})
		if err != nil {
			log.Println(err)
			return 500, err.Error()
		}

		message = "Task closed because last status stayed yet"
	} else {
		result, err = c.DB.Exec(
			"update task set status = (select ts.id from task_status ts "+
				"right join (select t.project, ts.status_level from task t "+
				"left join task_status ts on ts.id = t.status where t.id = ?) "+
				"t on t.project = ts.project where ts.status_level = t.status_level + 1) where id = ?",
			taskId,
			taskId,
		)
		if err != nil {
			log.Println("error in updating status:", err)
			return 500, err.Error()
		}

		message = "Task status set to next"
	}

	rowsAffected, _ := result.RowsAffected()

	w.Header().Set("Content-Type", "application/json")
	return c.makeContentResponse(200, message, struct {
		Name    string
		Content interface{}
	}{
		Name:    "Rows affected",
		Content: rowsAffected,
	})
}

func (c *Controller) ChangeTaskPriority(params martini.Params, w http.ResponseWriter, r *http.Request) (int, string) {
	contentType := r.Header.Get("Content-Type")
	if contentType != "application/json" {
		err := fmt.Sprintf("Unsupportable Content-Type header: %s", contentType)
		log.Println(err)
		return 500, err
	}
	taskId, err := strconv.ParseInt(params["id"], 10, 64)
	if err != nil {
		log.Println("error in parsing taskId", err)
		return 500, err.Error()
	}

	jsonPriority, err := ioutil.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		log.Println("error in reading body", err)
		return 500, err.Error()
	}
	priority := &struct {
		Priority string
	}{}
	err = json.Unmarshal(jsonPriority, priority)
	if err != nil {
		log.Println("error in unmarshalling")
		return 500, err.Error()
	}

	result, err := c.DB.Exec(
		"update task set priority = ? where id = ?",
		priority.Priority,
		taskId,
	)
	if err != nil {
		log.Println("error in updating priority:", err)
		return 500, err.Error()
	}

	rowsAffected, _ := result.RowsAffected()

	w.Header().Set("Content-Type", "application/json")
	return c.makeContentResponse(200, "Task priority updated", struct {
		Name    string
		Content interface{}
	}{
		Name:    "Rows affected",
		Content: rowsAffected,
	})
}

func (c *Controller) ChangeTaskDeadline(params martini.Params, w http.ResponseWriter, r *http.Request) (int, string) {
	contentType := r.Header.Get("Content-Type")
	if contentType != "application/json" {
		err := fmt.Sprintf("Unsupportable Content-Type header: %s", contentType)
		log.Println(err)
		return 500, err
	}
	taskId, err := strconv.ParseInt(params["id"], 10, 64)
	if err != nil {
		log.Println("error in parsing taskId", err)
		return 500, err.Error()
	}

	jsonDeadline, err := ioutil.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		log.Println("error in reading body", err)
		return 500, err.Error()
	}
	deadline := &struct {
		Deadline string
	}{}
	err = json.Unmarshal(jsonDeadline, deadline)
	if err != nil {
		log.Println("error in unmarshalling")
		return 500, err.Error()
	}
	timeDeadline, err := time.Parse("2006-01-02", deadline.Deadline)
	if err != nil {
		log.Println("error in parsing deadline: ", err)
		return 500, err.Error()
	}

	result, err := c.DB.Exec(
		"update task set deadline = ? where id = ?",
		timeDeadline,
		taskId,
	)
	if err != nil {
		log.Println("error in updating deadline:", err)
		return 500, err.Error()
	}

	row := c.DB.QueryRow("select project from task where id = ?", taskId)
	var projectId int64
	err = row.Scan(&projectId)
	if err != nil {
		log.Println("error in getting projectId: ", err)
		return 500, err.Error()
	}

	_, err = c.sendDataToStream("task", "deadline", struct {
		TaskId    int64
		Deadline  string
		ProjectId int64
	}{
		taskId,
		deadline.Deadline,
		projectId,
	})
	if err != nil {
		log.Println(err)
		return 500, err.Error()
	}

	rowsAffected, _ := result.RowsAffected()

	w.Header().Set("Content-Type", "application/json")
	return c.makeContentResponse(200, "Task deadline updated", struct {
		Name    string
		Content interface{}
	}{
		Name:    "Rows affected",
		Content: rowsAffected,
	})
}

func (c *Controller) CloseTask(params martini.Params, w http.ResponseWriter) (int, string) {
	taskId, err := strconv.ParseInt(params["id"], 10, 64)
	if err != nil {
		log.Println("error in parsing taskId", err)
		return 500, err.Error()
	}

	result, err := c.DB.Exec(
		"update task set is_closed = ? where id = ?",
		true,
		taskId,
	)
	if err != nil {
		log.Println("error in closing task")
		return 500, err.Error()
	}
	rowsAffected, _ := result.RowsAffected()

	row := c.DB.QueryRow("select executor from task where id = ?", taskId)
	var executorId int64
	err = row.Scan(&executorId)
	if err != nil {
		log.Println("error in getting userId: ", err)
		return 500, err.Error()
	}

	row = c.DB.QueryRow("select project from task where id = ?", taskId)
	var projectId int64
	err = row.Scan(&projectId)
	if err != nil {
		log.Println("error in getting projectId: ", err)
		return 500, err.Error()
	}

	_, err = c.sendDataToStream("task", "close", struct {
		TaskId     int64
		ExecutorId int64
		ProjectId  int64
	}{
		taskId,
		executorId,
		projectId,
	})
	if err != nil {
		log.Println(err)
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
