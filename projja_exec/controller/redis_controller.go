package controller

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"projja-exec/graph"
	"projja-exec/model"
	"strconv"

	"github.com/go-redis/redis/v8"
)

func (c *controller) saveNewProject(newProject *model.Project) error {
	project := graph.MakeNewProject(newProject)

	err := c.writeProjectToRedis(project)
	if err != nil {
		return err
	}
	return nil
}

func (c *controller) writeProjectToRedis(project *graph.Project) error {
	ctx := context.Background()

	byteGraph, err := json.Marshal(project.Graph)
	if err != nil {
		return err
	}

	status := c.Rds.Set(ctx, strconv.FormatInt(project.Id, 10), string(byteGraph), 0)
	if status.Err() != nil {
		return status.Err()
	}
	return nil
}

func (c *controller) ReadData(id int64) (*graph.Project, error) {
	val, err := c.Rds.Get(context.Background(), strconv.FormatInt(id, 10)).Result()
	switch {
	case err == redis.Nil:
		log.Println("key does not exist")
		return nil, err
	case err != nil:
		log.Println("Get failed", err)
		return nil, err
	case val == "":
		err = errors.New("value is empty")
		log.Println(err)
		return nil, err
	}

	log.Println(val)
	g := &graph.Graph{}
	err = json.Unmarshal([]byte(val), g)
	if err != nil {
		return nil, err
	}
	return &graph.Project{
		Id:    id,
		Graph: g,
	}, nil
}
