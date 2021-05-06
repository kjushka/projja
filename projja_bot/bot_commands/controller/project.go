package controller

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"projja_bot/services/memcached"
	"projja_bot/betypes"
	"projja_bot/logger"
	"strconv"
	"strings"

	"github.com/go-telegram-bot-api/telegram-bot-api"
)

// TODO тут должна быть проверка на то, что такой проект уже создан

func CreateProject(userName string, projectName string) string {
	if userName == "" {
		return "Вы не указали владельца проекта!"
	}

	if projectName == "" {
		return "Вы не указали название проекта!" 
	}

	user := GetUser(userName)
	if	user == nil {
		return fmt.Sprintf("Пользоватль с именем %s не зарегистрирован!", userName)
	}

	project := &betypes.Project{
		Name:   projectName,
		Owner: user,
		Status: "opened",
	}

	projectBytes, err := json.Marshal(project)
	logger.LogCommandResult(string(projectBytes))
	logger.ForError(err)

	resp, err := http.Post(betypes.GetPathToMySQl("http") + "api/project/create", "application/json", bytes.NewBuffer(projectBytes))
	logger.ForError(err)
	logger.LogCommandResult(resp.Status)
	fmt.Println(resp.Status)

	if resp.StatusCode >= 200 && resp.StatusCode < 300 {
		return fmt.Sprintf("Проект %s, с владельцем %s был создан", projectName, userName)
	}

	return "Что-то пошло не так..."
}

// Возвращаемые параметры ссылка на клавиатуру и количетво сообщений
func GetAllProjects(userName string) (tgbotapi.InlineKeyboardMarkup, int) {
	resp, err := http.Get(betypes.GetPathToMySQl("http") + fmt.Sprintf("api/user/%s/owner/all", userName))
	logger.ForError(err)
	fmt.Println(resp.Status)

	gettingProjects, err := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()
	logger.ForError(err)
	
	var projects *betypes.ProjectsList
	if err := json.Unmarshal(gettingProjects, &projects); err != nil {
    logger.ForError(err)
	}
	keyboard := tgbotapi.InlineKeyboardMarkup{}

	for i := 0; i < len(projects.Content); i++ {
		text := fmt.Sprintf("select_project %s %s", strconv.FormatInt(projects.Content[i].Id, 10), projects.Content[i].Name)

		var row []tgbotapi.InlineKeyboardButton
		btn := tgbotapi.NewInlineKeyboardButtonData(projects.Content[i].Name, text)

		row = append(row, btn)
		keyboard.InlineKeyboard = append(keyboard.InlineKeyboard, row)
	}

	return keyboard, len(projects.Content)
}

func AddMemberToProject(userName string) string {
	projectId, projectName, err := memcached.GetSelectedProject(userName)
	if err != nil {
		logger.ForError(err)
		return "Истекло время ожидания, заново выберете проект и пользователя!"
	}

	member, err := memcached.GetSelectedMember(userName)
	if err != nil {
		logger.ForError(err)
		return "Истекло время ожидания, заново выберете проект и пользователя!"
	}

	resp, err := http.Get(betypes.GetPathToMySQl("http") + fmt.Sprintf("api/project/%s/add/member/%s", projectId, member))
	logger.ForError(err)

	if(resp.StatusCode >= 500) {
		jsonUser, err := ioutil.ReadAll(resp.Body)
		defer resp.Body.Close()
		logger.ForError(err)

		var duplicateUser bool = strings.HasPrefix(string(jsonUser), "Error 1062:")
		if (duplicateUser) {
			return fmt.Sprintf("Пользователь %s уже является участником проекта %s!", member, projectName)
		}

		return "Неизвестная ошибка"
	}

	return fmt.Sprintf("Пользователь %s добавлен в проект %s!", member, projectName)
}

func GetProjectMembers(projectId string) (string, int) {
	resp, err := http.Get(betypes.GetPathToMySQl("http") + fmt.Sprintf("api/project/%s/members", projectId))
	logger.ForError(err)

	gettingMembers, err := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()
	logger.ForError(err)

	var members *betypes.MembersList
	if err := json.Unmarshal(gettingMembers, &members); err != nil {
    logger.ForError(err)
	}

	answer := "№ RealName UserName\n"
	for i := 0; i < len(members.Content); i++ {
		answer += fmt.Sprintf("%d. %s %s\n", i+1, members.Content[i].Name, members.Content[i].Username)
	}

	return answer, len(members.Content)
}

func RemoveMemberFromProject(projectOwner string, projectExecuter string) string {
	projectId, projectName, err := memcached.GetSelectedProject(projectOwner)
	if err != nil {
		logger.ForError(err)
		return "Истекло время ожидания, заново выберете проект и пользователя!"
	}
	logger.LogCommandResult(fmt.Sprintf("Remove %s from %s", projectExecuter, projectName) );

	fmt.Println(projectExecuter)
	resp, err := http.Get(betypes.GetPathToMySQl("http") + fmt.Sprintf("api/project/%s/remove/member/%s", projectId, projectExecuter))
	logger.ForError(err)
	fmt.Println(resp.Status)

	text := fmt.Sprintf("Пользователь %s был удален из проекта %s", projectExecuter, projectName)
	return text
}


func ChangeProjectName(projectOwner string) string {	
	projectId, _, err := memcached.GetSelectedProject(projectOwner)
	if err != nil {
		logger.ForError(err)
		return "Истекло время ожидания, заново выберете проект и пользователя!"
	}

	project := &betypes.Project{
		Name:   "new",
	}
	projectBytes, err := json.Marshal(project)
	logger.LogCommandResult(string(projectBytes))
	logger.ForError(err)

	fmt.Println(projectId)
	requestUrl := fmt.Sprintf("%sapi/project/14/change/name", betypes.GetPathToMySQl("http"))
	fmt.Println(requestUrl)
	resp, err := http.Post(requestUrl, "application/json", bytes.NewBuffer(projectBytes))
	fmt.Println(resp.Status)


	return "Изменено название проекта!"
}