package controller

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/go-martini/martini"
	"log"
	"projja_api/model"
	"strings"
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

func (c *Controller) getUserByUsername(params martini.Params) (*model.User, error) {
	username := params["uname"]
	row := c.DB.QueryRow("select id, name, username, telegram_id from users where username = ?", username)
	user := &model.User{}
	err := row.Scan(&user.Id, &user.Name, &user.Username, &user.TelegramId)
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
		queryArr[i] = fmt.Sprintf("(%v, (select id from skill where skill = %v))", id, skills[i])
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
