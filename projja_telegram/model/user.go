package model

import "time"

type User struct {
	Id         int64    `json:"id,omitempty"`
	Name       string   `json:"name"`
	Username   string   `json:"username"`
	TelegramId string   `json:"telegramId"`
	ChatId     int64    `json:"chatId"`
	Skills     []string `json:"skills,omitempty"`
}

type Project struct {
	Id     int64 `json:"id,omitempty"`
	Name   string
	Owner  *User
	Status string
}

type Task struct {
	Id          int64       `json:"id,omitempty"`
	Description string      `json:"description"`
	Project     *Project    `json:"project,omitempty"`
	Deadline    string      `json:"deadline"`
	Priority    string      `json:"priority,omitempty"`
	Status      *TaskStatus `json:"status,omitempty"`
	IsClosed    bool        `json:"isClosed,omitempty"`
	Executor    *User       `json:"executor,omitempty"`
	Skills      []string    `json:"skills,omitempty"`
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
