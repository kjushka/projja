package bot_commands

import (
	"strings"
)

func CreateProject(args string) string {
	if args == "" {
		return "Вы не указали название и владельца проекта!"
	}

	var argsArr[] string = strings.Split(args, " ")
	if (len(argsArr) == 1){
		return "Вы не указали владельца проекта!"
	}

	// получить владельца проекта

	// создать проект

	// userSkills := &betypes.Skills {
	// 	Skills: skills,
	// }

	return "Что-то пошло не так..."
}