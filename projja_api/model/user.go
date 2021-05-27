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

type Project struct {
	Id     int64
	Name   string
	Owner  *User
	Status string
}

type Task struct {
	Id          int64
	Description string
	Project     *Project
	Deadline    string
	Priority    string
	Status      *TaskStatus
	IsClosed    bool
	Executor    *User
	Skills      []string
}

type TaskStatus struct {
	Status string
	Level  int
}

type Answer struct {
	Id        int64
	Task      *Task
	Executor  *User
	MessageId int
	ChatId    int64
	Status    string
	SentAt    time.Time
}
