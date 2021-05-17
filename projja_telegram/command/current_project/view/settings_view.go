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

func ChangeProjectSetting(botUtil *util.BotUtil, project *model.Project) {
	msg := projectmenu.MakeSettingsMenu(botUtil.Message, project)
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
		case "change_name":
			msg = ChangeProjectName(botUtil, project)
			botUtil.Bot.Send(msg)
		case "change_status":
			msg = ChangeProjectStatus(botUtil, project)
			botUtil.Bot.Send(msg)
		case "change_members":
			ChangeProjectMembers(botUtil, project)
		case "change_statuses":
			ChangeProjectStatuses(botUtil, project)
		case "back_btn":
			return
		}

		msg = projectmenu.MakeSettingsMenu(botUtil.Message, project)
		botUtil.Bot.Send(msg)
	}
}

func ChangeProjectName(botUtil *util.BotUtil, project *model.Project) tgbotapi.MessageConfig {
	text := fmt.Sprintf("Введите новое название для проекта '%s'", project.Name)
	msg := tgbotapi.NewMessage(botUtil.Message.Chat.ID, text)

	keyboard := tgbotapi.InlineKeyboardMarkup{}
	row := make([]tgbotapi.InlineKeyboardButton, 0)
	cancelBtn := tgbotapi.NewInlineKeyboardButtonData("Отмена", "cancel_btn")
	row = append(row, cancelBtn)
	keyboard.InlineKeyboard = append(keyboard.InlineKeyboard, row)

	msg.ReplyMarkup = keyboard

	botUtil.Bot.Send(msg)

	projectName := ""
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
		case "cancel_btn":
			text = "Отмена смены названия проекта"
			msg := tgbotapi.NewMessage(botUtil.Message.Chat.ID, text)
			return msg
		default:
			if command == "" {
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
		case "yes_btn":
			text, _ = controller.ChangeProjectName(project, projectName)
			goto LOOP
		case "no_btn":
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
		case "yes_btn":
			text, _ = controller.ChangeProjectStatus(project, newStatus)
			goto LOOP
		case "no_btn":
			text = "Отмена смены статуса проекта"
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
