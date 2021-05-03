package controller

import (
	"context"
	"encoding/json"
	"github.com/go-redis/redis/v8"
	"log"
)

func (c *Controller) sendDataToStream(stream string, key string, data interface{}) (string, error) {
	jsonData, err := json.Marshal(data)
	if err != nil {
		log.Println("error during marshalling data for redis:", err)
		return "", err
	}
	result := c.Rds.XAdd(context.Background(), &redis.XAddArgs{
		Stream: stream,
		Values: map[string]interface{}{key: string(jsonData)},
	})
	insertId, err := result.Result()
	if err != nil {
		log.Println("error during getting xadd result:", err)
		return "", err
	}
	return insertId, nil
}
