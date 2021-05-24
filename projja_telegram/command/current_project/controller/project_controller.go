package controller

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"projja_telegram/config"
	"projja_telegram/model"
)

func ChangeProjectName(project *model.Project, newName string) (string, bool) {
	nameStruct := &struct {
		Name string
	}{newName}

	errorText := "Во время смены названия проекта произошла ошибка\nПопробуйте позже ещё раз"

	jsonNameStruct, err := json.Marshal(nameStruct)
	if err != nil {
		log.Println("error in marshalling name: ", err)
		return errorText, false
	}

	resp, err := http.Post(config.GetAPIAddr()+fmt.Sprintf("/project/%d/change/name", project.Id),
		"application/json",
		bytes.NewBuffer(jsonNameStruct),
	)
	if err != nil {
		log.Println("error in request for changing project name: ", err)
		return errorText, false
	}

	if resp.StatusCode == http.StatusInternalServerError {
		log.Println("error in request for changing project name")
		return errorText, false
	}

	project.Name = newName
	return "Название проекта успешно изменено", true
}

func ChangeProjectStatus(project *model.Project, newStatus string) (string, bool) {
	var url string
	if newStatus == "opened" {
		url = config.GetAPIAddr() + fmt.Sprintf("/project/%d/open", project.Id)
	} else {
		url = config.GetAPIAddr() + fmt.Sprintf("/project/%d/close", project.Id)
	}

	errorText := "Во время смены статуса проекта произошла ошибка\nПопробуйте позже ещё раз"

	resp, err := http.Get(url)
	if err != nil {
		log.Println("error in request for changing project status: ", err)
		return errorText, false
	}

	if resp.StatusCode == http.StatusInternalServerError {
		log.Println("error in request for changing project status")
		return errorText, false
	}

	project.Status = newStatus
	return "Статус проекта успешно изменен", true
}

func GetMembers(project *model.Project) ([]*model.User, bool) {
	resp, err := http.Get(config.GetAPIAddr() +
		fmt.Sprintf("/project/%d/members", project.Id),
	)

	if err != nil {
		log.Println("error in getting members: ", err)
		return nil, false
	}
	if resp.StatusCode == http.StatusInternalServerError {
		log.Println("error in getting members")
		return nil, false
	}

	respData := &struct {
		Description string
		Content     []*model.User
	}{}
	jsonBody, err := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()
	if err != nil {
		log.Println("error in reading response body: ", err)
		return nil, false
	}

	err = json.Unmarshal(jsonBody, respData)
	if err != nil {
		log.Println("error in unmarshalling members: ", err)
		return nil, false
	}

	return respData.Content, true
}

func GetUser(username string) (*model.User, string) {
	response, err := http.Get(config.GetAPIAddr() + "/user/get/" + username)

	textError := "Возникла ошибка информации об участнике"

	if err != nil {
		log.Println("error in getting user by username: ", err)
		return nil, textError
	}
	if response.StatusCode == http.StatusInternalServerError {
		log.Println("error in getting user by username")
		return nil, textError
	}
	if response.StatusCode == http.StatusNotFound {
		log.Println("no such user with username: ", username)
		return nil, fmt.Sprintf("Пользователь с username '%s' не зарегистрирован", username)
	}

	jsonBody, err := ioutil.ReadAll(response.Body)
	defer response.Body.Close()
	if err != nil {
		log.Println("error in reading response body: ", err)
		return nil, textError
	}
	responseStruct := &struct {
		Description string
		Content     *model.User
	}{}

	err = json.Unmarshal(jsonBody, responseStruct)
	if err != nil {
		log.Println("error in unmarshalling: ", err)
		return nil, textError
	}

	return responseStruct.Content, ""
}

func AddMember(project *model.Project, member *model.User) (string, bool) {
	resp, err := http.Get(config.GetAPIAddr() +
		fmt.Sprintf("/project/%d/add/member/%s", project.Id, member.Username))
	textError := "При добавлении участника возникла ошибка"

	if err != nil {
		log.Println("error in adding member by username: ", err)
		return textError, false
	}

	if resp.StatusCode == http.StatusInternalServerError {
		log.Println("error in adding member by username")
		return textError, false
	}

	return "Участник успешно добавлен в проект", true
}

func RemoveMember(project *model.Project, member *model.User) (string, bool) {
	resp, err := http.Get(config.GetAPIAddr() +
		fmt.Sprintf("/project/%d/remove/member/%s", project.Id, member.Username))
	textError := "При удалении участника возникла ошибка"

	if err != nil {
		log.Println("error in removing member by username: ", err)
		return textError, false
	}

	if resp.StatusCode == http.StatusInternalServerError {
		log.Println("error in removing member by username")
		return textError, false
	}

	return "Участник успешно удален из проекта", true
}

func GetStatuses(project *model.Project) ([]*model.TaskStatus, bool) {
	resp, err := http.Get(config.GetAPIAddr() +
		fmt.Sprintf("/project/%d/statuses", project.Id),
	)

	if err != nil {
		log.Println("error in getting statuses: ", err)
		return nil, false
	}
	if resp.StatusCode == http.StatusInternalServerError {
		log.Println("error in getting statuses")
		return nil, false
	}

	respData := &struct {
		Description string
		Content     []*model.TaskStatus
	}{}
	jsonBody, err := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()
	if err != nil {
		log.Println("error in reading response body: ", err)
		return nil, false
	}

	err = json.Unmarshal(jsonBody, respData)
	if err != nil {
		log.Println("error in unmarshalling statuses: ", err)
		return nil, false
	}

	return respData.Content, true
}

func CreateTaskStatus(project *model.Project, status *model.TaskStatus) (string, bool) {
	errorText := "Во время добавления статуса задач произошла ошибка\nПопробуйте позже ещё раз"

	jsonTaskStatus, err := json.Marshal(status)
	if err != nil {
		log.Println("error in marshalling task status: ", err)
		return errorText, false
	}

	resp, err := http.Post(config.GetAPIAddr()+fmt.Sprintf("/project/%d/create/status", project.Id),
		"application/json",
		bytes.NewBuffer(jsonTaskStatus),
	)
	if err != nil {
		log.Println("error in request for creating task status: ", err)
		return errorText, false
	}

	if resp.StatusCode == http.StatusInternalServerError {
		log.Println("error in request for creating task status")
		return errorText, false
	}

	return "Статус задач успешно создан", true
}

func RemoveTaskStatus(project *model.Project, status *model.TaskStatus) (string, bool) {
	errorText := "Во время удаления статуса задач произошла ошибка\nПопробуйте позже ещё раз"

	jsonTaskStatus, err := json.Marshal(status)
	if err != nil {
		log.Println("error in marshalling task status: ", err)
		return errorText, false
	}

	resp, err := http.Post(config.GetAPIAddr()+fmt.Sprintf("/project/%d/remove/status", project.Id),
		"application/json",
		bytes.NewBuffer(jsonTaskStatus),
	)
	if err != nil {
		log.Println("error in request for removing task status: ", err)
		return errorText, false
	}

	if resp.StatusCode == http.StatusInternalServerError {
		log.Println("error in request for removing task status")
		return errorText, false
	}

	return "Статус задач успешно удален", true
}

func GetProjectTasks(project *model.Project) ([]*model.Task, bool) {
	resp, err := http.Get(config.GetAPIAddr() +
		fmt.Sprintf("/project/%d/get/tasks/process", project.Id),
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
