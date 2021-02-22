package model

type User struct {
	Id         int64
	Name       string
	Username   string
	TelegramId string
	Skills     []string
}

type Project struct {
	Name   string
	Owner  *User
	Status string
}
