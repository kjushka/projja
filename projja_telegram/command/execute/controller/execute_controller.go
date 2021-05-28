package controller

import (
	"bytes"
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
	response, err := http.Get(config.GetAPIAddr() + fmt.Sprintf("/answer/last/%d/%s", task.Id, user.UserName))

	if err != nil {
		log.Println("error in getting user by username: ", err)
		return nil, false
	}
	if response.StatusCode == http.StatusInternalServerError {
		log.Println("error in getting user by username")
		return nil, false
	}
	if response.StatusCode == http.StatusNotFound {
		log.Println("no any answer for task ", task.Description)
		return nil, true
	}

	jsonBody, err := ioutil.ReadAll(response.Body)
	defer response.Body.Close()
	if err != nil {
		log.Println("error in reading response body: ", err)
		return nil, false
	}
	responseStruct := &struct {
		Description string
		Content     *model.Answer
	}{}

	err = json.Unmarshal(jsonBody, responseStruct)
	if err != nil {
		log.Println("error in unmarshalling: ", err)
		return nil, false
	}

	return responseStruct.Content, true
}

func AddAnswer(answer *model.Answer) (string, bool) {
	errorText := "Во время добавления ответа произошла ошибка\nПопробуйте позже ещё раз"

	jsonAnswer, err := json.Marshal(answer)
	if err != nil {
		log.Println("error in marshalling answer: ", err)
		return errorText, false
	}

	resp, err := http.Post(config.GetAPIAddr()+"/answer/create",
		"application/json",
		bytes.NewBuffer(jsonAnswer),
	)
	if err != nil {
		log.Println("error in request for creating answer: ", err)
		return errorText, false
	}

	if resp.StatusCode == http.StatusInternalServerError {
		log.Println("error in request for creating answer")
		return errorText, false
	}

	return "Ответ успешно добавлен", true
}
