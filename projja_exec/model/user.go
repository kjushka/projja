package model

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
	Deadline    string
	Executor    *User
	Skills      []string
}
