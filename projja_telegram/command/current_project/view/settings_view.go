package view

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"projja_telegram/command/current_project/controller"
	projectmenu "projja_telegram/command/current_project/menu"
	"projja_telegram/command/util"
	"projja_telegram/model"
)

func ChangeProjectSetting(botUtil *util.BotUtil, project *model.Project) {
	msg := projectmenu.MakeSettingsMenu(botUtil.Message, project)
	botUtil.Bot.Send(msg)

	for update := range botUtil.Updates {
		mes := update.Message
		command := ""

		if mes.Text != "" {
			command = mes.Text
		}

		switch command {
		case "Сменить название":
			msg = ChangeProjectName(botUtil, project)
			botUtil.Bot.Send(msg)
		case "Открыть/закрыть проект":
			msg = ChangeProjectStatus(botUtil, project)
			botUtil.Bot.Send(msg)
		case "Участники проекта":
			ChangeProjectMembers(botUtil, project)
		case "Статусы задач":
			ChangeProjectStatuses(botUtil, project)
		case "Назад":
			return
		default:
			msg = util.GetUnknownMessage(botUtil)
			botUtil.Bot.Send(msg)
		}

		msg = projectmenu.MakeSettingsMenu(botUtil.Message, project)
		botUtil.Bot.Send(msg)
	}
}

func ChangeProjectName(botUtil *util.BotUtil, project *model.Project) tgbotapi.MessageConfig {
	text := fmt.Sprintf("Введите новое название для проекта '%s'", project.Name)
	msg := tgbotapi.NewMessage(botUtil.Message.Chat.ID, text)

	row := make([]tgbotapi.KeyboardButton, 0)
	cancelBtn := tgbotapi.NewKeyboardButton("Отмена")
	row = append(row, cancelBtn)

	keyboard := tgbotapi.NewReplyKeyboard(row)
	msg.ReplyMarkup = keyboard

	botUtil.Bot.Send(msg)

	projectName := ""
	for update := range botUtil.Updates {
		mes := update.Message
		command := ""

		if mes.Text != "" {
			command = mes.Text
		}

		switch command {
		case "Отмена":
			text = "Отмена смены названия проекта"
			msg := tgbotapi.NewMessage(botUtil.Message.Chat.ID, text)
			return msg
		default:
			if command == "" {
				botUtil.Bot.Send(msg)
				continue
			}
			projectName = command
		}

		if projectName != "" {
			break
		}
	}

	acceptingString := fmt.Sprintf("Вы действительно хотите изменить название проекта на '%s'?", projectName)
	msg = util.GetAcceptingMessage(botUtil.Message, acceptingString)

	botUtil.Bot.Send(msg)

	for update := range botUtil.Updates {
		mes := update.Message
		command := ""

		if mes.Text != "" {
			command = mes.Text
		}

		switch command {
		case "Да":
			text, _ = controller.ChangeProjectName(project, projectName)
			goto LOOP
		case "Нет":
			text = "Отмена смены названия проекта"
			goto LOOP
		default:
			text = "Пожалуйста, выберите один из вариантов"
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

func ChangeProjectStatus(botUtil *util.BotUtil, project *model.Project) tgbotapi.MessageConfig {
	var newStatus string
	if project.Status == "opened" {
		newStatus = "closed"
	} else {
		newStatus = "opened"
	}

	acceptingString := fmt.Sprintf("Вы хотите поменять статус проекта с '%s' на '%s'", project.Status, newStatus)
	msg := util.GetAcceptingMessage(botUtil.Message, acceptingString)

	botUtil.Bot.Send(msg)

	var text string
	for update := range botUtil.Updates {
		mes := update.Message
		command := ""

		if mes.Text != "" {
			command = mes.Text
		}

		switch command {
		case "Да":
			text, _ = controller.ChangeProjectStatus(project, newStatus)
			goto LOOP
		case "Нет":
			text = "Отмена смены статуса проекта"
			goto LOOP
		default:
			text = "Пожалуйста, выберите один из вариантов"
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
