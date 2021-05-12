package controller

import (
	"bytes"
	"encoding/json"
	"fmt"
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
		log.Println("error in request for changing project name: ", err)
		return errorText, false
	}

	if resp.StatusCode == http.StatusInternalServerError {
		log.Println("error in request for changing project name")
		return errorText, false
	}

	project.Status = newStatus
	return "Статус проекта успешно изменен", true
}
