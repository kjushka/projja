package model

import "time"

type User struct {
	Id         int64
	Name       string
	Username   string
	TelegramId string
	ChatId     int64
	Skills     []string
}

type Task struct {
	Id          int64
	Description string
	Deadline    time.Time
	Priority    string
	Executor    *User
	Skills      []string
}

type Project struct {
	Id     int64
	Name   string
	Owner  *User
	Status string
}
