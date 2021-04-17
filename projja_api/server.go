package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"projja_api/controller"

	"github.com/go-martini/martini"
	_ "github.com/go-sql-driver/mysql"

	_ "github.com/lib/pq"
)

const (
	DSN  = "root:Password123#@!@tcp(localhost:3306)/projja?&charset=utf8&interpolateParams=true&parseTime=true"
	addr = ":8080"
)

func main() {
	host := os.Getenv("DATABASE_HOST")
	name := os.Getenv("DATABASE_NAME")
	user := os.Getenv("DATABASE_USER")
	pass := os.Getenv("DATABASE_PASS")
	log.Println(host, name, user, pass)
	dsn := fmt.Sprintf(
		"%v:%v@tcp(%v:3306)/%v?&charset=utf8&interpolateParams=true&parseTime=true",
		user,
		pass,
		host,
		name,
	)

	db, err := sql.Open("mysql", dsn)
	if err != nil {
		log.Fatal(err)
	}
	db.Ping()
	defer db.Close()
	fmt.Println("Connected to db")

	m := martini.Classic()
	c := &controller.Controller{
		DB: db,
	}
	m.Group("/api", func(router martini.Router) {
		router.Group("/user", func(r martini.Router) {
			r.Post("/register", c.Register)
			r.Get("/get/:uname", c.GetUserByUsername)
			r.Post("/:uname/skills", c.SetSkillsToUser)
			r.Get("/:uname/owner/open", c.GetOpenUserProjects)
			r.Get("/:uname/owner/all", c.GetAllUserProjects)
			r.Post("/:uname/update", c.UpdateUserData)
			r.Get("/:uname/member/opened", c.GetOpenProjectsWhereMember)
			r.Get("/:uname/member/all", c.GetAllProjectsWhereMember)
			r.Get("/:uname/executor", c.GetExecuteTasks)
		})
		router.Group("/project", func(r martini.Router) {
			r.Post("/create", c.CreateProject)
			r.Post("/:id/change/name", c.ChangeProjectName)
			r.Post("/:id/change/status", c.ChangeProjectStatus)
			r.Get("/:id/members", c.GetProjectMembers)
			r.Get("/:id/add/member/:uname", c.AddMemberToProject)
			r.Get("/:id/remove/member/:uname", c.RemoveMemberFromProject)
			r.Post("/:id/create/status", c.CreateProjectTaskStatus)
			r.Post("/:id/remove/status", c.RemoveStatusFromProject)
			r.Get("/:id/statuses", c.GetProjectStatuses)
			r.Post("/:id/create/task", c.CreateTask)
			r.Get("/:id/get/tasks/all", c.GetAllProjectTasks)
			r.Get("/:id/get/tasks/process", c.GetProcessProjectTasks)
		})
		router.Group("/task", func(r martini.Router) {
			r.Get("/get/:id", c.GetTask)
			r.Post("/:id/change/executor", c.ChangeTaskExecutor)
			r.Post("/:id/change/description", c.ChangeTaskDescription)
			r.Post("/:id/set/skills", c.SetSkillsToTask)
			r.Get("/:id/change/status/previous", c.SetPreviousTaskStatus)
			r.Get("/:id/change/status/next", c.SetNextTaskStatus)
			r.Post("/:id/change/priority", c.ChangeTaskPriority)
			r.Post("/:id/change/deadline", c.ChangeTaskDeadline)
		})
	})
	m.RunOnAddr(addr)
}
