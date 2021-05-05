package model

import "time"

type User struct {
	Id         int64
	Name       string
	Username   string
	TelegramId string
	Skills     []string
}

type Task struct {
	Id          int64
	Description string
	Deadline    time.Time
	Executor    *User
	Skills      []string
}

type Project struct {
	Id     int64
	Name   string
	Owner  *User
	Status string
}
