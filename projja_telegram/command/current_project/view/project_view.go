package view

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"projja_telegram/command/current_project/controller"
	projectmenu "projja_telegram/command/current_project/menu"
	"projja_telegram/command/projects/menu"
	"projja_telegram/command/util"
	"projja_telegram/model"
	"strings"
)

func WorkWithProject(message *util.MessageData,
	bot *tgbotapi.BotAPI,
	updates tgbotapi.UpdatesChannel,
	project *model.Project,
) {
	defer func(message *util.MessageData, bot *tgbotapi.BotAPI) {
		msg, _, _ := menu.MakeProjectsMenu(message, 1, 10)
		bot.Send(msg)
	}(message, bot)

	msg := projectmenu.MakeProjectMenu(message, project)
	bot.Send(msg)

	for update := range updates {
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
			msg = ChangeProjectName(message, bot, updates, project)
		}

		msg = projectmenu.MakeProjectMenu(message, project)
		bot.Send(msg)
	}

	return
}

func ChangeProjectName(message *util.MessageData,
	bot *tgbotapi.BotAPI,
	updates tgbotapi.UpdatesChannel,
	project *model.Project,
) tgbotapi.MessageConfig {
	text := fmt.Sprintf("Введите новое название для проекта '%s'", project.Name)
	msg := tgbotapi.NewMessage(message.Chat.ID, text)
	bot.Send(msg)

	projectName := ""
	for update := range updates {
		mes := update.Message
		if mes == nil {
			continue
		}

		projectName = mes.Text
		break
	}

	acceptingString := fmt.Sprintf("Вы действительно хотите изменить название проекта на '%s'?", projectName)
	msg = util.GetAcceptingMessage(message, acceptingString)
	bot.Send(msg)

	for update := range updates {
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
			msg = tgbotapi.NewMessage(message.Chat.ID, text)
			bot.Send(msg)

			msg = util.GetAcceptingMessage(message, acceptingString)
			bot.Send(msg)
		}
	}

LOOP:
	msg = tgbotapi.NewMessage(message.Chat.ID, text)
	return msg
}
