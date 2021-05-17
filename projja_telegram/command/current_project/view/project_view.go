package view

import (
	projectmenu "projja_telegram/command/current_project/menu"
	"projja_telegram/command/util"
	"projja_telegram/model"
	"strings"
)

func WorkWithProject(botUtil *util.BotUtil, project *model.Project) {
	msg := projectmenu.MakeProjectMenu(botUtil.Message, project)
	botUtil.Bot.Send(msg)

	for update := range botUtil.Updates {
		mes := update.Message
		var command string

		if update.CallbackQuery != nil {
			response := strings.Split(update.CallbackQuery.Data, " ")
			command = response[0]
		} else if mes.IsCommand() {
			command = mes.Command()
		} else if mes.Text != "" {
			command = mes.Text
		}

		switch command {
		case "settings":
			ChangeProjectSetting(botUtil, project)
		case "back_btn":
			return
		default:
			msg = util.GetUnknownMessage(botUtil, command)
			botUtil.Bot.Send(msg)
		}

		msg = projectmenu.MakeProjectMenu(botUtil.Message, project)
		botUtil.Bot.Send(msg)
	}
}
