package controller

import (
	"bytes"
	"encoding/json"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"io/ioutil"
	"log"
	"net/http"
	"projja_telegram/command/util"
	"projja_telegram/config"
	"projja_telegram/model"
)

func GetProjects(user *tgbotapi.User, page int, count int) ([]*model.Project, bool) {
	projectsCount := count - (page-1)*4
	if projectsCount > 4 {
		projectsCount = 4
	}
	resp, err := http.Get(config.GetAPIAddr() +
		fmt.Sprintf(
			"/user/%s/owner/all?page=%d&count=%d",
			user.UserName,
			page,
			projectsCount,
		),
	)

	if err != nil {
		log.Println("error in getting projects: ", err)
		return nil, false
	}
	if resp.StatusCode == http.StatusInternalServerError {
		log.Println("error in getting projects")
		return nil, false
	}

	respData := &struct {
		Description string
		Content     []*model.Project
	}{}
	jsonBody, err := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()
	if err != nil {
		log.Println("error in reading response body: ", err)
		return nil, false
	}

	err = json.Unmarshal(jsonBody, respData)
	if err != nil {
		log.Println("error in unmarshalling projects: ", err)
		return nil, false
	}

	return respData.Content, true
}

func GetProjectsCount(user *tgbotapi.User) (int, bool) {
	resp, err := http.Get(config.GetAPIAddr() + fmt.Sprintf("/user/%s/owner/count", user.UserName))
	if err != nil {
		log.Println("error in getting projects count: ", err.Error())
		return 0, false
	}

	if resp.StatusCode == http.StatusInternalServerError {
		log.Println("error in getting projects count")
		return 0, false
	}

	respData := &struct {
		Description string
		Content     int
	}{}

	jsonBody, err := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()
	if err != nil {
		log.Println("error in reading response body: ", err)
		return 0, false
	}

	err = json.Unmarshal(jsonBody, respData)
	if err != nil {
		log.Println("error in unmarshalling projects: ", err)
		return 0, false
	}

	return respData.Content, true
}

func CreateNewProject(data *util.MessageData, projectName string) (string, bool) {
	user := util.TgUserToModelUser(data)
	project := &model.Project{
		Name:   projectName,
		Owner:  user,
		Status: "opened",
	}
	jsonProject, err := json.Marshal(project)

	returnText := "Возникла ошибка, попробуйте ещё раз через некоторое время"

	if err != nil {
		log.Println("error in marshalling project: ", err)
		return returnText, false
	}

	response, err := http.Post(config.GetAPIAddr()+"/project/create",
		"application/json",
		bytes.NewBuffer(jsonProject),
	)

	if err != nil {
		log.Println("error in sending creating project request: ", err)
		return returnText, false
	}

	if response.StatusCode == http.StatusInternalServerError {
		log.Println("error in register")
		return returnText, false
	}

	return fmt.Sprintf("Проект %s был успешно создан", projectName), true
}
