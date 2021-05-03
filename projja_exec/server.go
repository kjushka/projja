package main

import (
	"projja-exec/controller"

	"github.com/go-martini/martini"
	"github.com/go-redis/redis/v8"
)

const (
	Addr = ":8090"
)

func main() {
	c := controller.NewController(
		&redis.Options{
			Addr:     "localhost:6379",
			Password: "",
			DB:       0,
		})

	m := martini.Classic()
	m.Use(c.CheckContentType)
	m.Group("/exec", func(r martini.Router) {
		r.Post("/add/project", c.AddProject)
		r.Get("/get/:id", c.GetRedisData)
		r.Post("/project/:id/add/exec", c.AddExecutorToProject)
		r.Post("/project/:id/add/task", c.AddTaskToProject)
	})

	// done := make(chan error, 1)
	// go c.ListenRedisStream(context.Background(), done)
	// log.Println(<-done)

	m.RunOnAddr(Addr)
}
