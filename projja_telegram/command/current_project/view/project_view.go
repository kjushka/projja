package view

import (
	projectmenu "projja_telegram/command/current_project/menu"
	"projja_telegram/command/util"
	"projja_telegram/model"
)

func WorkWithProject(botUtil *util.BotUtil, project *model.Project) {
	msg := projectmenu.MakeProjectMenu(botUtil.Message, project)
	botUtil.Bot.Send(msg)

	for update := range botUtil.Updates {
		mes := update.Message
		command := ""

		if mes.Text != "" {
			command = mes.Text
		}

		switch command {
		case "Настройки проекта":
			ChangeProjectSetting(botUtil, project)
		case "Управление задачами":
			ManageProjectTasks(botUtil, project)
		case "Ответы на задачи":
		case "Назад":
			return
		default:
			msg = util.GetUnknownMessage(botUtil)
			botUtil.Bot.Send(msg)
		}

		msg = projectmenu.MakeProjectMenu(botUtil.Message, project)
		botUtil.Bot.Send(msg)
	}
}
