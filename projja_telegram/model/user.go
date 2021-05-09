package model

type User struct {
	Id         int64    `json:"id,omitempty"`
	Name       string   `json:"name"`
	Username   string   `json:"username"`
	TelegramId string   `json:"telegramId"`
	Skills     []string `json:"skills,omitempty"`
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
