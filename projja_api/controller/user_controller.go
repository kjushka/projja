package controller

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"projja_api/model"
	"strings"
	"time"

	"github.com/go-martini/martini"
	"github.com/scylladb/go-set"
)

func (c *Controller) Register(w http.ResponseWriter, r *http.Request) (int, string) {
	contentType := r.Header.Get("Content-Type")
	if contentType != "application/json" {
		err := fmt.Sprintf("Unsupportable Content-Type header: %s", contentType)
		log.Println(err)
		return 500, err
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
		"select u.id, u.name, u.username, u.telegram_id from users u where username = ?",
		username,
	)
	err := row.Scan(
		&user.Id,
		&user.Name,
		&user.Username,
		&user.TelegramId,
	)
	if err != nil && err != sql.ErrNoRows {
		log.Println("error in getting user by username:", err)
		return 500, err.Error()
	}

	if err == sql.ErrNoRows {
		noUserErr := fmt.Errorf("no such user with username %s", username)
		log.Println(noUserErr)
		return 500, noUserErr.Error()
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
	contentType := r.Header.Get("Content-Type")
	if contentType != "application/json" {
		err := fmt.Sprintf("Unsupportable Content-Type header: %s", contentType)
		log.Println(err)
		return 500, err
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

	uniqueSkills := make([]string, 0)
	for _, s := range skills.Skills {
		skip := false
		for _, u := range uniqueSkills {
			if s == u {
				skip = true
				break
			}
		}
		if !skip {
			uniqueSkills = append(uniqueSkills, s)
		}
	}

	skills.Skills = uniqueSkills

	username := params["uname"]

	row := c.DB.QueryRow("select id from users where username = ?", username)
	var userId int64
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
	newSkills := make([]string, 0)
	newUserSkills := make([]string, len(skills.Skills))
	for i, skill := range skills.Skills {
		skill := strings.ToLower(skill)
		if !skillsSet.Has(skill) {
			newSkills = append(newSkills, fmt.Sprintf("('%v')", skill))
		}
		newUserSkills[i] = fmt.Sprintf("(%v, (select s.id from skill s where s.skill = '%v'))",
			userId,
			skill,
		)
	}

	if len(newSkills) != 0 {
		_, err := c.DB.Exec("insert ignore into skill (skill) values " +
			strings.Join(newSkills, ", "),
		)
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

	rows, err = c.DB.Query(
		"select p.id, p.name, p.status from project p "+
			"left join (select project, users from member) m on m.project = p.id "+
			"where m.users = ? and p.status = ?",
		userId,
		"opened",
	)

	if err != nil {
		log.Println("error in getting opened projects:", err)
		return 500, err.Error()
	}

	projects, err := c.scanProjects(rows, nil)
	if err != nil {
		log.Println("error in scanning rows:", err)
		return 500, err.Error()
	}

	projectsIds := make([]int64, len(projects))
	for i, v := range projects {
		projectsIds[i] = v.Id
	}

	_, err = c.sendDataToStream("exec", "skills", struct {
		UserId      int64
		Skills      []string
		ProjectsIds []int64
	}{
		userId,
		skills.Skills,
		projectsIds,
	})
	if err != nil {
		log.Println(err)
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

func (c *Controller) UpdateUserData(params martini.Params, w http.ResponseWriter, r *http.Request) (int, string) {
	contentType := r.Header.Get("Content-Type")
	if contentType != "application/json" {
		err := fmt.Sprintf("Unsupportable Content-Type header: %s", contentType)
		log.Println(err)
		return 500, err
	}
	username := params["uname"]
	userDataJson, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Println("error in reading body:", err)
		return 500, err.Error()
	}
	defer r.Body.Close()

	newUserInfo := &model.User{}
	err = json.Unmarshal(userDataJson, newUserInfo)
	if err != nil {
		log.Println("error in unmarshalling:", err)
		return 500, err.Error()
	}

	row := c.DB.QueryRow("select id from user where username = ?", username)
	var userId int64
	err = row.Scan(&userId)

	result, err := c.DB.Exec("update users set name = ?, username = ?, telegram_id = ? where id = ?",
		newUserInfo.Name, newUserInfo.Username, newUserInfo.TelegramId, userId)
	if err != nil {
		log.Println("error in updating info:", err)
		return 500, err.Error()
	}

	rows, err := c.DB.Query(
		"select p.id, p.name, p.status from project p "+
			"left join (select project, users from member) m on m.project = p.id "+
			"where m.users = ? and p.status = ?",
		userId,
		"opened",
	)

	if err != nil {
		log.Println("error in getting opened projects:", err)
		return 500, err.Error()
	}

	newUserInfo.Id = userId

	projects, err := c.scanProjects(rows, nil)
	if err != nil {
		log.Println("error in scanning rows:", err)
		return 500, err.Error()
	}

	projectsIds := make([]int64, len(projects))
	for i, v := range projects {
		projectsIds[i] = v.Id
	}

	_, err = c.sendDataToStream("exec", "info", struct {
		NewUserInfo *model.User
		ProjectsIds []int64
	}{
		newUserInfo,
		projectsIds,
	})

	rowsAffected, _ := result.RowsAffected()

	w.Header().Set("Content-Type", "application/json")
	return c.makeContentResponse(200, "Info updated", struct {
		Name    string
		Content interface{}
	}{
		Name:    "Rows affected",
		Content: rowsAffected,
	})
}

func (c *Controller) GetOpenProjectsWhereMember(params martini.Params, w http.ResponseWriter) (int, string) {
	user, err := c.getUserByUsername(params)
	if err != nil {
		log.Println("error in getting user:", err)
		return 500, err.Error()
	}

	rows, err := c.DB.Query(
		"select p.id, p.name, p.status from project p "+
			"left join (select project, users from member) m on m.project = p.id "+
			"where m.users = ? and p.status = ?",
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

func (c *Controller) GetAllProjectsWhereMember(params martini.Params, w http.ResponseWriter) (int, string) {
	user, err := c.getUserByUsername(params)
	if err != nil {
		log.Println("error in getting user:", err)
		return 500, err.Error()
	}

	rows, err := c.DB.Query(
		"select p.id, p.name, p.status from project p "+
			"left join (select project, users from member) m on m.project = p.id "+
			"where m.users = ?",
		user.Id,
	)

	if err != nil {
		log.Println("error in getting all projects:", err)
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

func (c *Controller) GetExecuteTasks(params martini.Params, w http.ResponseWriter) (int, string) {
	user, err := c.getUserByUsername(params)
	if err != nil {
		log.Println("error in getting user:", err)
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
			"where t.executor = ? and t.is_closed = 0;",
		user.Id,
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
	return c.makeContentResponse(200, "executor tasks", tasks)
}
