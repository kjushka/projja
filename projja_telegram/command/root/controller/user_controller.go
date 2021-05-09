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

func GetUser(username string) *model.User {
	if username == "" {
		return nil
	}

	response, err := http.Get(config.GetAPIAddr() + "/user/get/" + username)
	if err != nil {
		log.Println("error in getting user by username: ", err)
		return nil
	}
	if response.StatusCode == http.StatusInternalServerError ||
		response.StatusCode == http.StatusNotFound {
		return nil
	}

	jsonBody, err := ioutil.ReadAll(response.Body)
	defer response.Body.Close()
	if err != nil {
		log.Println("error in reading response body: ", err)
		return nil
	}
	responseStruct := &struct {
		Description string
		Content     *model.User
	}{}

	err = json.Unmarshal(jsonBody, responseStruct)
	if err != nil {
		log.Println("error in unmarshalling: ", err)
		return nil
	}

	return responseStruct.Content
}

func RegisterUser(tgUser *tgbotapi.User) (bool, string) {
	if tgUser.UserName == "" {
		return false, "У вас не установлен username\n" +
			"Пожалуйста, задайте его в настройках"
	}
	user := util.TgUserToModelUser(tgUser)
	jsonUser, err := json.Marshal(user)

	returnText := "Возникла ошибка, попробуйте ещё раз через некоторое время"

	if err != nil {
		log.Println("error in marshalling user: ", err)
		return false, returnText
	}

	response, err := http.Post(config.GetAPIAddr()+"/user/register",
		"application/json",
		bytes.NewBuffer(jsonUser),
	)

	if err != nil {
		log.Println("error in sending register request: ", err)
		return false, returnText
	}

	if response.StatusCode == http.StatusInternalServerError {
		log.Println("error in register")
		return false, returnText
	}

	return true, "Ваш профиль был успешно создан"
}

func SetSkills(username string, skills []string) bool {
	skillsStruct := &struct {
		Skills []string
	}{skills}

	jsonSkills, err := json.Marshal(skillsStruct)
	if err != nil {
		log.Println("error in marshalling skills: ", err)
		return false
	}

	resp, err := http.Post(config.GetAPIAddr()+fmt.Sprintf("/user/%s/skills", username),
		"application/json",
		bytes.NewBuffer(jsonSkills),
	)
	if err != nil {
		log.Println("error in sending set skills request: ", err)
		return false
	}

	if resp.StatusCode == http.StatusInternalServerError {
		log.Println("error in setting skills")
		return false
	}

	return true
}

func UpdateUserData(tgUser *tgbotapi.User) (bool, string) {
	user := util.TgUserToModelUser(tgUser)
	jsonUser, err := json.Marshal(user)

	returnText := "Возникла ошибка, попробуйте ещё раз через некоторое время"

	if err != nil {
		log.Println("error in marshalling user: ", err)
		return false, returnText
	}

	response, err := http.Post(config.GetAPIAddr()+fmt.Sprintf("/user/%s/update", user.TelegramId),
		"application/json",
		bytes.NewBuffer(jsonUser),
	)

	if err != nil {
		log.Println("error in sending update request: ", err)
		return false, returnText
	}

	if response.StatusCode == http.StatusInternalServerError {
		log.Println("error in update")
		return false, returnText
	}

	return true, "Данные профиля были успешно обновлены"
}
