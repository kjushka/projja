package controller

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"projja_bot/betypes"
	"projja_bot/logger"
	"strings"
	"strconv"

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
	_projectName := strings.Split(projectName, " ")[0]

	user := GetUser(userName)
	if	user == nil {
		return fmt.Sprintf("Пользоватль с именем %s не зарегистрирован!", userName)
	}

	project := &betypes.Project{
		Name:   _projectName,
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
		return fmt.Sprintf("Проект %s, с владельцем %s был создан", _projectName, userName)
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
		text := fmt.Sprintf("select_project %s %s", projects.Content[i].Name, strconv.FormatInt(projects.Content[i].Id, 10))

		var row []tgbotapi.InlineKeyboardButton
		btn := tgbotapi.NewInlineKeyboardButtonData(projects.Content[i].Name, text)

		row = append(row, btn)
		keyboard.InlineKeyboard = append(keyboard.InlineKeyboard, row)
	}

	return keyboard, len(projects.Content)
}

func AddMemberToProject(userName string) string {
	addedMember, err := betypes.MemCashed.Get(fmt.Sprintf("%s_member", userName))
	if err != nil {
		logger.ForError(err)
		return "Истекло время ожидания, заново выберете проект и пользователя!"
	}

	projectForAdd, err := betypes.MemCashed.Get(fmt.Sprintf("%s_poject", userName))
	if err != nil {
		logger.ForError(err)
		return "Истекло время ожидания, заново выберете проект и пользователя!"
	}
	
	member := string(addedMember.Value)
	projectId := string(projectForAdd.Value)
	fmt.Println(member)
	fmt.Println(projectId)

	resp, err := http.Get(betypes.GetPathToMySQl("http") + fmt.Sprintf("api/project/0/add/member/kjushka"))

	logger.ForError(err)
	fmt.Println(resp.Status)
	fmt.Println(resp.Body)

	return "test"
}


