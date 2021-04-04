package bot_commands

import (
	"fmt"
	"strings"
	"projja_bot/betypes"
	"projja_bot/logger"
	"bytes"
	"net/http"
	"encoding/json"
)

// TODO тут должна быть проверка на то, что такой проект уже создан

func CreateProject(args string) string {
	if args == "" {
		return "Вы не указали название и владельца проекта!"
	}

	var argsArr[] string = strings.Split(args, " ")
	if (len(argsArr) == 1){
		return "Вы не указали владельца проекта!"
	}
	var projectName string = strings.Split(args, " ")[0]
	var userName string = strings.Split(args, " ")[1]

	ans, user := GetUser(userName)
	if	user == nil {
		return ans
	}

	// TODO переделать ствтус!
	project := &betypes.Project{
		Name:   projectName,
		Owner: user,
		Status: "В работе!",
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

