package controller

// TODO: не уверен, что это хорошее решение,
// т.к. telegram-bot-api приходится подключать
// в каждом отдельном файле, где используются библиотеги этого api
import (
	"net/http"
	"encoding/json"
	"io/ioutil"
	"bytes"
	"github.com/go-telegram-bot-api/telegram-bot-api"
	"projja_bot/betypes"
	"projja_bot/logger"
	"fmt"
	"strings"
	"strconv"
	// "strconv"
)

func RegiserUser(from *tgbotapi.User) string {
	// Телега возвращает id типа int
	user := &betypes.User{
		Name: from.FirstName + " " + from.LastName,
		Username: from.UserName,
		TelegramId: strconv.Itoa(from.ID), 
	}

	messageBytes, err := json.Marshal(user)
	logger.LogCommandResult(string(messageBytes))
	logger.ForError(err)

	resp, err := http.Post(betypes.GetPathToMySQl("http") + "api/user/register", "application/json", bytes.NewBuffer(messageBytes))
	logger.ForError(err)

	// 500 ошибка может возвращаться, если ты пытаешься зарегать юзера, который уже есть в бд
	fmt.Println(resp.Status)
	logger.LogCommandResult(resp.Status)

	if(resp.StatusCode >= 500) {
		jsonUser, err := ioutil.ReadAll(resp.Body)
		defer resp.Body.Close()
		logger.ForError(err)

		var duplicateUser bool = strings.HasPrefix(string(jsonUser), "Error 1062:")
		if (duplicateUser) {
			return "Такой пользователь уже зарегистрирован!"
		}

		return "Неизвестная ошибка"
	} else if (resp.StatusCode < 300) {
		return fmt.Sprintf("Пользователь %s был успешно зарегистрирован!", from.UserName);
	}

	logger.LogCommandResult("Non-standard situation during registration.")
	return "Что-то пошло не так..." 
}

func GetUser(userName string) (*betypes.User) {
	if(userName == "") {
		return nil  
	}
	fmt.Println(userName)

	resp, err := http.Get(betypes.GetPathToMySQl("http") + "api/user/get/" + userName);
	logger.ForError(err)

	if resp.StatusCode == 500 {
		return nil 
	}

	getUserInfo, err := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()
	logger.ForError(err)
	logger.LogCommandResult(string(getUserInfo));

	var userAns *betypes.GetUserAnswer
	if err := json.Unmarshal(getUserInfo, &userAns); err != nil {
    logger.ForError(err)
	}
	
	return userAns.Content
} 

func SetSkills(userName string, args string) string {	
	if (userName == "") {
		return "Вы не указали имя пользователя!"
	} 
	if (args == "") {
		return "Вы не указали навыки пользователя!"
	}
	
	var skills[] string = strings.Split(args, " ")
	userSkills := &betypes.Skills {
		Skills: skills,
	}

	skillsBytes, err := json.Marshal(userSkills)
	logger.LogCommandResult(string(skillsBytes))
	logger.ForError(err)

	resp, err :=  http.Post(betypes.GetPathToMySQl("http") + fmt.Sprintf("api/user/%s/skills", userName), "application/json", bytes.NewBuffer(skillsBytes) )
	logger.ForError(err)
	logger.LogCommandResult(resp.Status)

	if resp.StatusCode == 404 || resp.StatusCode == 500 {
		return fmt.Sprintf("Пользователь с именем %s не зарегистрирован!", userName)
	}
	if resp.StatusCode == 202 {
		return "Навыки были успешно установлены!"
	}
	return "Что-то пошло не так..."
}


// Нужна ли эта функция
// func ChangeName()
