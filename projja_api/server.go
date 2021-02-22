package main

import (
	"database/sql"
	"fmt"
	"github.com/go-martini/martini"
	_ "github.com/go-sql-driver/mysql"
	"log"
	"projja_api/controller"

	_ "github.com/lib/pq"
)

const (
	DSN = "root:Password123#@!@tcp(localhost:3306)/projja?&charset=utf8&interpolateParams=true&parseTime=true"
	addr = ":8080"
)

func main() {
	db, err := sql.Open("mysql", DSN)
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
			r.Get("/get/:username", c.GetUserByUsername)
			r.Post("/:uname/skills", c.SetSkillsToUser)
			r.Get("/:uname/owner/open", c.GetOpenUserProjects)
			r.Get("/:uname/owner/all", c.GetAllUserProjects)
			r.Post("/:uname/change", c.ChangeUserName)
		})
	})
	m.RunOnAddr(addr)
}
