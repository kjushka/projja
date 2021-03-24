package controller

import (
	"fmt"
	"log"
	"net/http"
	"projja-exec/graph"

	"github.com/go-martini/martini"
	"github.com/go-redis/redis/v8"
)

type Controller struct {
	Rds      *redis.Client
	Projects map[int64]*graph.Project
}

func (c *Controller) AddProject(w http.ResponseWriter, r *http.Request) (int, string) {
	contentType := r.Header.Get("Content-Type")
	if contentType != "application/json" {
		err := fmt.Sprintf("Unsupportable Content-Type header: %s", contentType)
		log.Println(err)
		return 500, err
	}

	return 200, ""
}

func (c *Controller) AddTaskToProject(params martini.Params, w http.ResponseWriter, r *http.Request) (int, string) {
	contentType := r.Header.Get("Content-Type")
	if contentType != "application/json" {
		err := fmt.Sprintf("Unsupportable Content-Type header: %s", contentType)
		log.Println(err)
		return 500, err
	}

	return 200, ""
}

/*func (c *Controller) ListenRedisStream(ctx context.Context, done chan<- error) {
	for {
		xStreamSlice := c.Rds.XRead(&redis.XReadArgs{Block: 0, Streams: []string{"ws", "$"}})
		xReadResult, err := xStreamSlice.Result()
		if err != nil {
			if err == redis.Nil {
				continue
			} else {
				panic(err)
			}
		}
		rdsMessages := xReadResult[len(xReadResult)-1].Messages
		for _, rdsMessage := range rdsMessages {
			rdsMap := rdsMessage.Values
			if val, ok := rdsMap["message"]; ok {
				go c.unpackRdsMessage(val, "message")
			}
		}
	}
}

func (c *Controller) unpackRdsMessage(val interface{}, mesType string) {
}*/
