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
	"strings"
)

func (c *Controller) Register(w http.ResponseWriter, r *http.Request) (int, string) {
	if r.Header.Get("Content-Type") != "application/json" {
		return 500, "unsupportable content-type"
	}
	jsonUser, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Println("error during reading body:", err)
		return 500, err.Error()
	}
	newUser := &model.User{}
	err = json.Unmarshal(jsonUser, newUser)
	if err != nil {
		log.Println("error during unmarshalling:", err)
		return 500, err.Error()
	}
	result, err := c.DB.Exec(
		"insert into users (`name`, `username`, `telegram_id`) values (?, ?, ?)",
		newUser.Name,
		newUser.Username,
		newUser.TelegramId,
	)
	if err != nil {
		log.Println("error during create user:", err)
		return 500, err.Error()
	}
	rowsAffected, _ := result.RowsAffected()
	w.Header().Set("Content-Type", "application/json")
	return c.makeContentResponse(202, "user registered", struct {
		Name    string
		Content interface{}
	}{
		Name:    "Rows affected",
		Content: rowsAffected,
	})
}

func (c *Controller) GetUserByUsername(params martini.Params, w http.ResponseWriter) (int, string) {
	username := params["uname"]
	user := &model.User{}
	row := c.DB.QueryRow(
		"select u.name, u.username, u.telegram_id from users u where username = ?",
		username,
	)
	err := row.Scan(
		&user.Name,
		&user.Username,
		&user.TelegramId,
	)
	if err != nil && err != sql.ErrNoRows {
		log.Println("error in getting user by username:", err)
		return 500, err.Error()
	}

	skills, err := c.getSkillsByUser(username)
	if err != nil {
		log.Println("error in getting skills:", err)
		return 500, err.Error()
	}
	user.Skills = skills

	w.Header().Set("Content-Type", "application/json")
	return c.makeContentResponse(200, "Current user:", user)
}

func (c *Controller) SetSkillsToUser(params martini.Params, r *http.Request, w http.ResponseWriter) (int, string) {
	if r.Header.Get("Content-Type") != "application/json" {
		return 500, "unsupportable content-type"
	}
	jsonSkills, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Println("error during reading body:", err)
		return 500, err.Error()
	}
	skills := &struct {
		Skills []string
	}{}
	err = json.Unmarshal(jsonSkills, skills)
	if err != nil {
		log.Println("error during unmarshalling:", err)
		return 500, err.Error()
	}
	username := params["uname"]

	row := c.DB.QueryRow("select id from users where username = ?", username)
	var userId int
	err = row.Scan(&userId)
	if err != nil && err != sql.ErrNoRows {
		log.Println("error in getting user id:", err)
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
	newUserSkills := make([]string, len(skills.Skills))
	for i, skill := range skills.Skills {
		skill := strings.ToLower(skill)
		if !skillsSet.Has(skill) {
			newSkills = append(newSkills, fmt.Sprintf("('%v')", strings.ToLower(skill)))
		}
		newUserSkills[i] = fmt.Sprintf("(%v, (select s.id from skill s where s.skill = '%v'))",
			userId,
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
		"delete from users_skill where users = ?",
		userId,
	)
	if err != nil {
		log.Println("error in deleting skills:", err)
		return 500, err.Error()
	}

	result, err := c.DB.Exec(
		"insert into users_skill (users, skill) values " +
			strings.Join(newUserSkills, ", "),
	)
	if err != nil {
		log.Println("error in creating new users_skill:", err)
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

func (c *Controller) GetOpenUserProjects(params martini.Params, w http.ResponseWriter) (int, string) {
	user, err := c.getUserByUsername(params)
	if err != nil {
		log.Println("error in getting user:", err)
		return 500, err.Error()
	}

	rows, err := c.DB.Query(
		"select p.id, p.name, p.status from project p where p.owner = ? and p.status = ?",
		user.Id,
		"opened",
	)

	if err != nil {
		log.Println("error in getting opened projects:", err)
		return 500, err.Error()
	}

	projects, err := c.scanProjects(rows, user)
	if err != nil {
		log.Println("error in scanning rows:", err)
		return 500, err.Error()
	}

	w.Header().Set("Content-Type", "application/json")
	return c.makeContentResponse(200, "projects", projects)
}

func (c *Controller) GetAllUserProjects(params martini.Params, w http.ResponseWriter) (int, string) {
	user, err := c.getUserByUsername(params)
	if err != nil {
		log.Println("error in getting user:", err)
		return 500, err.Error()
	}

	rows, err := c.DB.Query(
		"select p.id, p.name, p.status from project p where p.owner = ?",
		user.Id,
	)

	if err != nil {
		log.Println("error in getting opened projects:", err)
		return 500, err.Error()
	}

	projects, err := c.scanProjects(rows, user)
	if err != nil {
		log.Println("error in scanning rows:", err)
		return 500, err.Error()
	}

	w.Header().Set("Content-Type", "application/json")
	return c.makeContentResponse(200, "projects", projects)
}

func (c *Controller) ChangeUserName(params martini.Params, w http.ResponseWriter, r *http.Request) (int, string) {
	username := params["uname"]
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Println("error in reading body:", err)
		return 500, err.Error()
	}
	defer r.Body.Close()

	nameStruct := &struct {
		Name string
	}{}
	err = json.Unmarshal(body, nameStruct)
	if err != nil {
		log.Println("error in unmarshalling:", err)
		return 500, err.Error()
	}

	result, err := c.DB.Exec("update users set name = ? where username = ?", nameStruct.Name, username)
	if err != nil {
		log.Println("error in updating name:", err)
		return 500, err.Error()
	}

	rowsAffected, _ := result.RowsAffected()

	w.Header().Set("Content-Type", "application/json")
	return c.makeContentResponse(200, "Name changed", struct {
		Name    string
		Content interface{}
	}{
		Name:    "Rows affected",
		Content: rowsAffected,
	})
}

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
