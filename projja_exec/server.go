package main

import (
	"projja-exec/controller"

	"github.com/go-martini/martini"
	"github.com/go-redis/redis/v8"
)

const (
	Addr = ":8080"
)

func main() {
	c := &controller.Controller{
		Rds: redis.NewClient(&redis.Options{
			Addr:     "redis:6379",
			Password: "",
			DB:       0,
		}),
	}

	m := martini.Classic()
	m.Group("/exec", func(r martini.Router) {
		r.Post("/add/project", c.AddProject)
	})

	// done := make(chan error, 1)
	// go c.ListenRedisStream(context.Background(), done)
	// log.Println(<-done)

	m.RunOnAddr(Addr)
}
