package controller

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"projja_api/model"
	"strings"
	"time"
)

func (c *Controller) getSkillsByUser(username string) ([]string, error) {
	rows, err := c.DB.Query(
		"select s.skill from users_skill us "+
			"left join skill s on s.id = us.skill "+
			"inner join ("+
			"select * from users u where u.username = ?"+
			") u on u.id = us.users",
		username,
	)
	if err != nil && err != sql.ErrNoRows {
		log.Println("error in getting skills by username:", err)
		return nil, err
	}
	skills := []string{}
	for rows.Next() {
		var skill string
		err = rows.Scan(&skill)
		if err != nil {
			log.Println("error in scan skills:", err)
			return nil, err
		}
		skills = append(skills, skill)
	}
	return skills, nil
}

func (c *Controller) getSkillsTask(taskId int64) ([]string, error) {
	rows, err := c.DB.Query(
		"select s.skill from task_skill ts "+
			"left join skill s on s.id = ts.skill "+
			"inner join ("+
			"select * from task t where t.id = ?"+
			") t on t.id = ts.task",
		taskId,
	)
	if err != nil && err != sql.ErrNoRows {
		log.Println("error in getting skills by task id:", err)
		return nil, err
	}
	skills := []string{}
	for rows.Next() {
		var skill string
		err = rows.Scan(&skill)
		if err != nil {
			log.Println("error in scan skills:", err)
			return nil, err
		}
		skills = append(skills, skill)
	}
	return skills, nil
}

func (c *Controller) getUserByUsername(username string) (*model.User, error) {
	row := c.DB.QueryRow("select id, name, username, telegram_id, chat_id from users where username = ?", username)
	user := &model.User{}
	err := row.Scan(&user.Id, &user.Name, &user.Username, &user.TelegramId, &user.ChatId)
	if err != nil && err != sql.ErrNoRows {
		log.Println("error in getting user:", err)
		return nil, err
	}
	return user, err
}

func (c *Controller) scanProjects(rows *sql.Rows, user *model.User) ([]*model.Project, error) {
	projects := make([]*model.Project, 0)
	for rows.Next() {
		project := &model.Project{}
		err := rows.Scan(&project.Id, &project.Name, &project.Status)
		if err != nil {
			log.Println("error in scanning rows:", err)
			return nil, err
		}
		project.Owner = user
		projects = append(projects, project)
	}
	return projects, nil
}

func (c *Controller) makeContentResponse(code int, desc string, content interface{}) (int, string) {
	response := response{
		desc,
		content,
	}
	byteResponse, err := json.Marshal(response)
	if err != nil {
		log.Println("Error during content marshalling:", err.Error())
		return 500, err.Error()
	}
	return code, string(byteResponse)
}

func (c *Controller) setSkillsToTask(skills []string, id int64) (int64, error) {
	query := "insert into task_skill (task, skill) values "
	queryArr := make([]string, len(skills))
	for i := range skills {
		queryArr[i] = fmt.Sprintf("(%v, (select id from skill where skill = '%v'))", id, skills[i])
	}
	query += strings.Join(queryArr, ", ")
	result, err := c.DB.Exec(
		query,
	)
	if err != nil {
		return 0, err
	}
	rowsAffected, _ := result.RowsAffected()
	return rowsAffected, err
}

func (c *Controller) getTaskById(taskId int64) (*model.Task, error) {
	row := c.DB.QueryRow("select t.id, t.description, p.id, p.name, p.ow_id, p.ow_name, p.ow_username, "+
		"p.ow_telegram_id, p.ow_chat_id, p.status, t.deadline, t.priority, ts.status, ts.status_level, "+
		"e.id, e.name, e.username, e.telegram_id, e.chat_id from task t "+
		"left join (select p.id, p.name, u.id ow_id, u.name ow_name, "+
		"u.username ow_username, u.telegram_id ow_telegram_id, u.chat_id ow_chat_id, p.status "+
		"from project p left join users u on u.id = p.owner) p on p.id = t.project "+
		"left join task_status ts on ts.id = t.status "+
		"left join users e on t.executor = e.id "+
		"where t.id = ?;",
		taskId,
	)
	task := &model.Task{}
	task.Project = &model.Project{}
	task.Project.Owner = &model.User{}
	task.Status = &model.TaskStatus{}
	task.Executor = &model.User{}
	var deadline time.Time
	isClosed := 0

	err := row.Scan(
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

	if err != nil && err != sql.ErrNoRows {
		log.Println("error in scanning task:", err)
		return nil, err
	}

	task.Deadline = deadline.Format("2006-01-02")
	if isClosed == 1 {
		task.IsClosed = true
	}

	skills, err := c.getSkillsTask(taskId)
	if err != nil {
		log.Println("error in getting skills:", err)
	}
	task.Skills = skills

	return task, nil
}
