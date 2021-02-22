package model

type User struct {
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
