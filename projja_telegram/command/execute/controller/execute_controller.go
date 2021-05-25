package controller

import (
	"encoding/json"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"io/ioutil"
	"log"
	"net/http"
	"projja_telegram/config"
	"projja_telegram/model"
)

func GetExecutedTasks(user *tgbotapi.User) ([]*model.Task, bool) {
	resp, err := http.Get(config.GetAPIAddr() +
		fmt.Sprintf("/user/%s/executor", user.UserName),
	)

	if err != nil {
		log.Println("error in getting tasks: ", err)
		return nil, false
	}
	if resp.StatusCode == http.StatusInternalServerError {
		log.Println("error in getting tasks")
		return nil, false
	}

	respData := &struct {
		Description string
		Content     []*model.Task
	}{}
	jsonBody, err := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()
	if err != nil {
		log.Println("error in reading response body: ", err)
		return nil, false
	}

	err = json.Unmarshal(jsonBody, respData)
	if err != nil {
		log.Println("error in unmarshalling tasks: ", err)
		return nil, false
	}

	return respData.Content, true
}

func GetLastAnswer(user *tgbotapi.User, task *model.Task) (*model.Answer, bool) {
	return nil, true
}
