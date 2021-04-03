package betypes

import (
	"encoding/json"
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
}

// TODO Тут стоит дабвить проверку на ошибки
func (obj *GetUserAnswer) UnmarshalJSON(b []byte) error {

	// Преобразуем входное значение к мапе
	var f interface{}
	json.Unmarshal(b, &f)
	m := f.(map[string]interface{})

	//получаем контент
	foomap := m["Content"]
	v := foomap.(map[string]interface{})

	// TODO хз Go -забеал
	skills := []string{}
  for _, value := range v["Skills"].(map[string]interface{}) {
		skills = append(skills, value.(string))
  }

	// TODO хз, почему-то из базы ID Возвращается в виде float 64
	obj.Content = &User{
		Id: int64(v["Id"].(float64)),
		Name: v["Name"].(string),
		Username: v["Username"].(string),
		TelegramId: int(v["TelegramId"].(float64)),
		Skills: skills, 
	}

	// var userAns *User
	// if err := json.Unmarshal([]byte(v), &userAns); err != nil {
  //   logger.ForError(err)
	// }

	// obj.Content = userAns["Content"].(*User) 

	return nil
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

type Skills struct {
	Skills []string
}

