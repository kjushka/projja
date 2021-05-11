package view

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"projja_telegram/command/current_project/controller"
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

			mes = update.CallbackQuery.Message
			mes.From = update.CallbackQuery.From
		} else if mes.IsCommand() {
			command = mes.Command()
		} else if mes.Text != "" {
			command = mes.Text
		}

		switch command {
		case "change_name":
			msg = ChangeProjectName(botUtil, project)
		case "projects_menu":
			return
		}
		msg = projectmenu.MakeProjectMenu(botUtil.Message, project)
		botUtil.Bot.Send(msg)
	}
}

func ChangeProjectName(botUtil *util.BotUtil, project *model.Project) tgbotapi.MessageConfig {
	text := fmt.Sprintf("Введите новое название для проекта '%s'", project.Name)
	msg := tgbotapi.NewMessage(botUtil.Message.Chat.ID, text)
	botUtil.Bot.Send(msg)

	projectName := ""
	for update := range botUtil.Updates {
		mes := update.Message
		if mes == nil {
			continue
		}

		projectName = mes.Text
		break
	}

	acceptingString := fmt.Sprintf("Вы действительно хотите изменить название проекта на '%s'?", projectName)
	msg = util.GetAcceptingMessage(botUtil.Message, acceptingString)
	botUtil.Bot.Send(msg)

	for update := range botUtil.Updates {
		mes := update.Message
		var command string

		if update.CallbackQuery != nil {
			response := strings.Split(update.CallbackQuery.Data, " ")
			command = response[0]

			mes = update.CallbackQuery.Message
			mes.From = update.CallbackQuery.From
		} else if mes.IsCommand() {
			command = mes.Command()
		} else if mes.Text != "" {
			command = mes.Text
		}

		switch command {
		case "yes":
			text, _ = controller.ChangeProjectName(project, projectName)
			goto LOOP
		case "no":
			text = "Отмена смены названия проекта"
			goto LOOP
		default:
			text = "Неизвестная команда"
			msg = tgbotapi.NewMessage(botUtil.Message.Chat.ID, text)
			botUtil.Bot.Send(msg)

			msg = util.GetAcceptingMessage(botUtil.Message, acceptingString)
			botUtil.Bot.Send(msg)
		}
	}

LOOP:
	msg = tgbotapi.NewMessage(botUtil.Message.Chat.ID, text)
	return msg
}
