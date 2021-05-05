package main

import (
	"context"
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
			Addr:     "redis:6379",
			Password: "",
			DB:       0,
		})

	m := martini.Classic()
	m.Use(c.CheckContentType)
	m.Group("/exec", func(r martini.Router) {
		r.Get("/get/:id", c.GetRedisData)
		r.Post("/project/:id/calc/task", c.CalculateTaskExecutor)
	})

	ctx := context.Background()
	go c.ListenExecStream(ctx)
	go c.ListenProjectStream(ctx)
	go c.ListenTaskStream(ctx)

	m.RunOnAddr(Addr)
}
