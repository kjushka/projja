package controller

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/go-martini/martini"
	"github.com/scylladb/go-set"
	"io/ioutil"
	"log"
	"net/http"
	"projja_api/model"
	"strconv"
	"strings"
	"time"
)

func (c *Controller) GetTask(params martini.Params, w http.ResponseWriter) (int, string) {
	tasktId, err := strconv.ParseInt(params["id"], 10, 64)
	if err != nil {
		log.Println("error in parsing taskId", err)
		return 500, err.Error()
	}

	row := c.DB.QueryRow("select t.id, t.description, p.id, p.name, p.ow_id, p.ow_name, p.ow_username, "+
		"p.ow_telegram_id, p.status, t.deadline, t.priority, ts.status, ts.level, t.is_closed, "+
		"e.id, e.name, e.username, e.telegram_id from task t "+
		"left join (select p.id, p.name, u.id ow_id, u.name ow_name, "+
		"u.username ow_username, u.telegram_id ow_telegram_id, p.status "+
		"from project p left join users u on u.id = p.owner) p on p.id = t.project "+
		"left join task_status ts on ts.id = t.status "+
		"left join users e on t.executor = e.id "+
		"where t.id = ?;",
		tasktId,
	)
	task := &model.Task{}
	task.Project = &model.Project{}
	task.Project.Owner = &model.User{}
	task.Status = &model.TaskStatus{}
	task.Executor = &model.User{}
	var deadline time.Time
	isClosed := 0

	err = row.Scan(
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
		&isClosed,
		&task.Executor.Id,
		&task.Executor.Name,
		&task.Executor.Username,
		&task.Executor.TelegramId,
	)

	if err != nil && err != sql.ErrNoRows {
		log.Println("error in scanning task:", err)
		return 500, err.Error()
	}

	task.Deadline = deadline.Format("2006-01-02")
	if isClosed == 1 {
		task.IsClosed = true
	}

	skills, err := c.getSkillsTask(tasktId)
	if err != nil {
		log.Println("error in getting skills:", err)
	}
	task.Skills = skills

	w.Header().Set("Content-Type", "application/json")
	return c.makeContentResponse(200, "current task", task)
}

func (c *Controller) ChangeTaskExecutor(params martini.Params, w http.ResponseWriter, r *http.Request) (int, string) {
	tasktId, err := strconv.ParseInt(params["id"], 10, 64)
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

	result, err := c.DB.Exec(
		"update task set executor = (select id from users where username = ?) where id = ?",
		executor.Username,
		tasktId,
	)
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
	tasktId, err := strconv.ParseInt(params["id"], 10, 64)
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
		tasktId,
	)
	if err != nil {
		log.Println("error in updating description:", err)
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
	tasktId, err := strconv.ParseInt(params["id"], 10, 64)
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
	newSkills := []string{}
	newTaskSkills := make([]string, len(skills.Skills))
	for i, skill := range skills.Skills {
		skill := strings.ToLower(skill)
		if !skillsSet.Has(skill) {
			newSkills = append(newSkills, fmt.Sprintf("('%v')", strings.ToLower(skill)))
		}
		newTaskSkills[i] = fmt.Sprintf("(%v, (select s.id from skill s where s.skill = '%v'))",
			tasktId,
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
		tasktId,
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
	return c.makeContentResponse(202, "skills set", struct {
		Name    string
		Content interface{}
	}{
		Name:    "Rows affected",
		Content: rowsAffected,
	})
}

func (c *Controller) SetPreviousTaskStatus(params martini.Params, w http.ResponseWriter) (int, string) {
	tasktId, err := strconv.ParseInt(params["id"], 10, 64)
	if err != nil {
		log.Println("error in parsing taskId", err)
		return 500, err.Error()
	}

	row := c.DB.QueryRow(
		"select count(*) from task_status ts "+
			"right join (select t.project, ts.level from task t "+
			"left join task_status ts on ts.id = t.status where t.id = ?) "+
			"t on t.project = ts.project where ts.level <= t.level - 1;",
		tasktId,
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
			"right join (select t.project, ts.level from task t "+
			"left join task_status ts on ts.id = t.status where t.id = ?) "+
			"t on t.project = ts.project where ts.level = t.level - 1) where id = ?;",
		tasktId,
		tasktId,
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
	tasktId, err := strconv.ParseInt(params["id"], 10, 64)
	if err != nil {
		log.Println("error in parsing taskId", err)
		return 500, err.Error()
	}

	row := c.DB.QueryRow(
		"select count(*) from task_status ts "+
			"right join (select t.project, ts.level from task t "+
			"left join task_status ts on ts.id = t.status where t.id = ?) "+
			"t on t.project = ts.project where ts.level >= t.level + 1;",
		tasktId,
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
			"right join (select t.project, ts.level from task t "+
			"left join task_status ts on ts.id = t.status where t.id = ?) "+
			"t on t.project = ts.project where ts.level = t.level + 1) where id = ?;",
		tasktId,
		tasktId,
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

func (c *Controller) ChangeTaskPriority(params martini.Params, w http.ResponseWriter, r *http.Request) (int, string) {
	tasktId, err := strconv.ParseInt(params["id"], 10, 64)
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
		tasktId,
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
	tasktId, err := strconv.ParseInt(params["id"], 10, 64)
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

	result, err := c.DB.Exec(
		"update task set deadline = ? where id = ?",
		timeDeadline,
		tasktId,
	)
	if err != nil {
		log.Println("error in updating deadline:", err)
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