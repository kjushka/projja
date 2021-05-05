package controller

import (
	"database/sql"
	"github.com/go-redis/redis/v8"
)

type Controller struct {
	Rds *redis.Client
	DB  *sql.DB
}
