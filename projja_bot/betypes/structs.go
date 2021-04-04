package betypes

import (
	"encoding/json"
	// "errors"
	// "fmt"
	// "fmt"
	// "projja_bot/logger"
)

type User struct {
	Id         int64
	Name       string `json:"name"`
	Username   string `json:"username"`
	TelegramId int `json:"telegramId"`
	Skills     []string
}

type GetUserAnswer struct {
	Description string
	Content	*User
	IsEmpty bool
}

func (obj *GetUserAnswer) UnmarshalJSON(b []byte) error {

	// Преобразуем входное значение к мапе
	var f interface{}
	json.Unmarshal(b, &f)
	m := f.(map[string]interface{})

	//получаем контент
	foomap := m["Content"]
	v := foomap.(map[string]interface{})

	skills := []string{}
  for _, value := range v["Skills"].([]interface{}) {
		skills = append(skills, value.(string))
  }
	userName := v["Name"].(string)

	// TODO хз, почему-то из базы ID Возвращается в виде float 64
	obj.Content = &User{
		Id: int64(v["Id"].(float64)),
		Name: userName,
		Username: v["Username"].(string),
		TelegramId: int(v["TelegramId"].(float64)),
		Skills: skills, 
	}

	// Т.К. у сервера нет ответа о том, что пользовател уже есть в базе
	// поэтому я сделаю его сам
	if userName == "" {
		obj.IsEmpty = true
		// Либо можно сделать так
		// return errors.New(fmt.Sprintf("500 user named %s does not exist", ))
	} else {
		obj.IsEmpty = false
	}

	return nil
}

type Project struct {
	Id     int64
	Name   string
	Owner  *User
	Status string
}

type ProjectsList struct {
	Content []*Project
}

func (obj *ProjectsList) UnmarshalJSON(b []byte) error {
	var f interface{}
	json.Unmarshal(b, &f)
	m := f.(map[string]interface{})

	foomap := m["Content"]

	v := foomap.([]interface{})

	projects := []*Project{}

  for _, line := range v {
		l := line.(map[string]interface{})

		// TODO для желания можно добавить юзера, хотя он особо не нужен
		projects = append(projects, 
			&Project{
				Name: l["Name"].(string),
				Id: int64(l["Id"].(float64)),
				Status: l["Status"].(string),
				// Name: projName,
			})
  }

	obj.Content = projects

	return nil
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

type Skills struct {
	Skills []string
}

