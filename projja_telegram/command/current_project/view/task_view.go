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

func ManageTask(botUtil *util.BotUtil, project *model.Project, task *model.Task) {
	msg := projectmenu.MakeTaskMenu(botUtil.Message, task)
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
		case "description":
			msg = ChangeTaskDescription(botUtil, task)
			botUtil.Bot.Send(msg)
		case "deadline":
			msg = ChangeTaskDeadline(botUtil, task)
			botUtil.Bot.Send(msg)
		case "priority":
			msg = ChangeTaskPriority(botUtil, task)
			botUtil.Bot.Send(msg)
		case "executor":
			msg = ChangeTaskExecutor(botUtil, project, task)
			botUtil.Bot.Send(msg)
		case "close_task":
			msg = CloseTask(botUtil, task)
			botUtil.Bot.Send(msg)
			return
		case "back_btn":
			return
		default:
			msg = util.GetUnknownMessage(botUtil, command)
			botUtil.Bot.Send(msg)
		}

		msg = projectmenu.MakeTaskMenu(botUtil.Message, task)
		botUtil.Bot.Send(msg)
	}
}

func ChangeTaskDescription(botUtil *util.BotUtil, task *model.Task) tgbotapi.MessageConfig {
	text := fmt.Sprintf("Введите новое описание для задачи '%s'", task.Description)
	msg := tgbotapi.NewMessage(botUtil.Message.Chat.ID, text)

	keyboard := tgbotapi.InlineKeyboardMarkup{}
	row := make([]tgbotapi.InlineKeyboardButton, 0)
	cancelBtn := tgbotapi.NewInlineKeyboardButtonData("Отмена", "cancel_btn")
	row = append(row, cancelBtn)
	keyboard.InlineKeyboard = append(keyboard.InlineKeyboard, row)

	msg.ReplyMarkup = keyboard

	botUtil.Bot.Send(msg)

	taskDescription := ""
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
			text = "Отмена смены описания задачи"
			msg := tgbotapi.NewMessage(botUtil.Message.Chat.ID, text)
			return msg
		default:
			if command == "" {
				continue
			}
			taskDescription = command
		}

		if taskDescription != "" {
			break
		}
	}

	acceptingString := fmt.Sprintf("Вы действительно хотите изменить описание задачи на '%s'?", taskDescription)
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
			text, _ = controller.ChangeTaskDescription(task, taskDescription)
			goto LOOP
		case "no_btn":
			text = "Отмена смены описания задачи"
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

func ChangeTaskDeadline(botUtil *util.BotUtil, task *model.Task) tgbotapi.MessageConfig {
	deadline, cancelStatus := listenForDeadline(botUtil)
	if !cancelStatus {
		text := "Отмена изменения дедлайна задачи"
		msg := tgbotapi.NewMessage(botUtil.Message.Chat.ID, text)
		return msg
	}

	acceptingString := fmt.Sprintf(
		"Вы действительно хотите изменить дедлайн задачи на '%s'?",
		deadline.Format("2006-01-02"),
	)
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
			text, _ = controller.ChangeTaskDeadline(task, deadline)
			goto LOOP
		case "no_btn":
			text = "Отмена изменения дедлайна задачи"
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

func ChangeTaskPriority(botUtil *util.BotUtil, task *model.Task) tgbotapi.MessageConfig {
	priority, cancelStatus := listenForPriority(botUtil)
	if !cancelStatus {
		text := "Отмена изменения приоритета задачи"
		msg := tgbotapi.NewMessage(botUtil.Message.Chat.ID, text)
		return msg
	}

	acceptingString := fmt.Sprintf(
		"Вы действительно хотите изменить приоритет задачи на '%s'?",
		priority,
	)
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
			text, _ = controller.ChangeTaskPriority(task, priority)
			goto LOOP
		case "no_btn":
			text = "Отмена изменения приоритета задачи"
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

func ChangeTaskExecutor(botUtil *util.BotUtil, project *model.Project, task *model.Task) tgbotapi.MessageConfig {
	executor, errorText := listenForExecutor(botUtil, project, "Отмена изменения исполнителя задачи")
	if executor == nil {
		text := errorText
		msg := tgbotapi.NewMessage(botUtil.Message.Chat.ID, text)
		return msg
	}

	acceptingString := fmt.Sprintf(
		"Вы действительно хотите изменить исполнителя задачи на:\n"+
			"Имя: %s\nUsername: %s\n",
		executor.Name,
		executor.Username,
	)
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
			text, _ = controller.ChangeTaskExecutor(task, executor)
			goto LOOP
		case "no_btn":
			text = "Отмена изменения исполнителя задачи"
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

func CloseTask(botUtil *util.BotUtil, task *model.Task) tgbotapi.MessageConfig {
	acceptingString := fmt.Sprintf(
		"Вы действительно хотите закрыть задачу '%s'",
		task.Description,
	)
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
			text, _ = controller.CloseTask(task)
			goto LOOP
		case "no_btn":
			text = "Отмена закрытия задачи"
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
