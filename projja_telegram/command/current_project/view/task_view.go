package view

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"projja_telegram/command/current_project/controller"
	projectmenu "projja_telegram/command/current_project/menu"
	"projja_telegram/command/util"
	"projja_telegram/model"
)

func ManageTask(botUtil *util.BotUtil, project *model.Project, task *model.Task) {
	msg := projectmenu.MakeTaskMenu(botUtil.Message, task)
	botUtil.Bot.Send(msg)

	for update := range botUtil.Updates {
		mes := update.Message
		command := ""

		if mes.Text != "" {
			command = mes.Text
		}

		switch command {
		case "Изменить описание":
			msg = ChangeTaskDescription(botUtil, task)
			botUtil.Bot.Send(msg)
		case "Изменить дедлайн":
			msg = ChangeTaskDeadline(botUtil, task)
			botUtil.Bot.Send(msg)
		case "Изменить приоритет":
			msg = ChangeTaskPriority(botUtil, task)
			botUtil.Bot.Send(msg)
		case "Изменить исполнителя":
			msg = ChangeTaskExecutor(botUtil, project, task)
			botUtil.Bot.Send(msg)
		case "Закрыть задачу":
			msg = CloseTask(botUtil, task)
			botUtil.Bot.Send(msg)
			return
		case "Назад":
			return
		default:
			msg = util.GetUnknownMessage(botUtil)
			botUtil.Bot.Send(msg)
		}

		msg = projectmenu.MakeTaskMenu(botUtil.Message, task)
		botUtil.Bot.Send(msg)
	}
}

func ChangeTaskDescription(botUtil *util.BotUtil, task *model.Task) tgbotapi.MessageConfig {
	text := fmt.Sprintf("Введите новое описание для задачи '%s'", task.Description)
	msg := tgbotapi.NewMessage(botUtil.Message.Chat.ID, text)

	row := make([]tgbotapi.KeyboardButton, 0)
	cancelBtn := tgbotapi.NewKeyboardButton("Отмена")
	row = append(row, cancelBtn)

	keyboard := tgbotapi.NewReplyKeyboard(row)
	msg.ReplyMarkup = keyboard

	botUtil.Bot.Send(msg)

	taskDescription := ""
	for update := range botUtil.Updates {
		mes := update.Message
		command := ""

		if mes.Text != "" {
			command = mes.Text
		}

		switch command {
		case "Отмена":
			text = "Отмена смены описания задачи"
			msg := tgbotapi.NewMessage(botUtil.Message.Chat.ID, text)
			return msg
		default:
			if command == "" {
				botUtil.Bot.Send(msg)
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
		command := ""

		if mes.Text != "" {
			command = mes.Text
		}

		switch command {
		case "Да":
			text, _ = controller.ChangeTaskDescription(task, taskDescription)
			goto LOOP
		case "Нет":
			text = "Отмена смены описания задачи"
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
		command := ""

		if mes.Text != "" {
			command = mes.Text
		}

		switch command {
		case "Да":
			text, _ = controller.ChangeTaskDeadline(task, deadline)
			goto LOOP
		case "Нет":
			text = "Отмена изменения дедлайна задачи"
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

func ChangeTaskPriority(botUtil *util.BotUtil, task *model.Task) tgbotapi.MessageConfig {
	priority, cancelStatus := listenForPriority(botUtil)
	if !cancelStatus {
		text := "Отмена изменения приоритета задачи"
		msg := tgbotapi.NewMessage(botUtil.Message.Chat.ID, text)
		return msg
	}

	var newPriority string
	switch priority {
	case "critical":
		newPriority = "критический"
	case "high":
		newPriority = "высокий"
	case "medium":
		newPriority = "средний"
	case "low":
		newPriority = "низкий"
	}
	acceptingString := fmt.Sprintf(
		"Вы действительно хотите изменить приоритет задачи на '%s'?",
		newPriority,
	)
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
			text, _ = controller.ChangeTaskPriority(task, priority)
			goto LOOP
		case "Нет":
			text = "Отмена изменения приоритета задачи"
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
		command := ""

		if mes.Text != "" {
			command = mes.Text
		}

		switch command {
		case "Да":
			prevExecutor := task.Executor
			changeText, status := controller.ChangeTaskExecutor(task, executor)
			text = changeText

			if status {
				notification := fmt.Sprintf("Вы больше не работаете над задачей '%s'", task.Description)
				msg := tgbotapi.NewMessage(prevExecutor.ChatId, notification)
				botUtil.Bot.Send(msg)

				notification = fmt.Sprintf("Вам назначена новая задача '%s'", task.Description)
				msg = tgbotapi.NewMessage(executor.ChatId, notification)
				botUtil.Bot.Send(msg)
			}

			goto LOOP
		case "Нет":
			text = "Отмена изменения исполнителя задачи"
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
		command := ""

		if mes.Text != "" {
			command = mes.Text
		}

		switch command {
		case "Да":
			text, _ = controller.CloseTask(task)
			goto LOOP
		case "Нет":
			text = "Отмена закрытия задачи"
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
