package view

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	controller2 "projja_telegram/command/current_project/controller"
	"projja_telegram/command/current_project/view"
	"projja_telegram/command/execute/controller"
	"projja_telegram/command/execute/menu"
	"projja_telegram/command/util"
	"projja_telegram/model"
	"strings"
)

func ExecuteTasks(botUtil *util.BotUtil) {
	page := 1
	tasks, msg, status := ShowExecutedTasksMenu(botUtil, page)
	botUtil.Bot.Send(msg)
	if !status {
		return
	}

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
		case "prev_page":
			page--
		case "next_page":
			page++
		case "back_btn":
			return
		default:
			text, index, status := view.IsTaskId(command, len(tasks), page)
			msg = tgbotapi.NewMessage(botUtil.Message.Chat.ID, text)
			botUtil.Bot.Send(msg)
			if status {
				ManageExecutorAnswers(botUtil, tasks[index])
			}
		}

		tasks, msg, status = ShowExecutedTasksMenu(botUtil, page)
		botUtil.Bot.Send(msg)
		if !status {
			return
		}
	}
}

func ShowExecutedTasksMenu(botUtil *util.BotUtil, page int) ([]*model.Task, tgbotapi.MessageConfig, bool) {
	tasks, status := controller.GetExecutedTasks(botUtil.Message.From)
	if !status {
		errorText := "Не удалось получить список задач\n" +
			"Попробуйте позже"
		msg := tgbotapi.NewMessage(botUtil.Message.Chat.ID, errorText)
		return nil, msg, false
	}

	count := len(tasks) - (page-1)*10
	if count > 10 {
		count = 10
	}
	msg := menu.MakeExecutedTasksMenu(botUtil.Message, tasks, page, count)
	return tasks, msg, true
}

func ManageExecutorAnswers(botUtil *util.BotUtil, task *model.Task) {
	answerMenu, status := MakeAddAnswerMenu(botUtil, task)
	botUtil.Bot.Send(answerMenu)
	if !status {
		return
	}

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
		case "back_btn":
			return
		case "add_answer":
			msg := AddAnswer(botUtil, task)
			botUtil.Bot.Send(msg)
		default:
			msg := util.GetUnknownMessage(botUtil, command)
			botUtil.Bot.Send(msg)
		}

		answerMenu, status = MakeAddAnswerMenu(botUtil, task)
		botUtil.Bot.Send(answerMenu)
		if !status {
			return
		}
	}
}

func MakeAddAnswerMenu(botUtil *util.BotUtil, task *model.Task) (tgbotapi.MessageConfig, bool) {
	answer, status := controller.GetLastAnswer(botUtil.Message.From, task)
	if !status {
		msg := tgbotapi.NewMessage(
			botUtil.Message.Chat.ID,
			"При получении последнего решения произошла ошибка\n"+
				"Попробуйте позже ещё раз",
		)
		return msg, false
	}

	var answerAsString string
	if answer == nil {
		answerAsString = "Вы ещё не отправляли решение для этой задачи"
	} else {
		msg := tgbotapi.NewForward(botUtil.Message.Chat.ID, answer.ChatId, answer.MessageId)
		botUtil.Bot.Send(msg)
		answerAsString = fmt.Sprintf("Ваш последний ответ: %s", answer.Status)
	}
	answerMenu := tgbotapi.NewMessage(botUtil.Message.Chat.ID, answerAsString)
	keyboard := tgbotapi.InlineKeyboardMarkup{}
	row := make([]tgbotapi.InlineKeyboardButton, 0)
	addAnswerBtn := tgbotapi.NewInlineKeyboardButtonData("Добавить ответ", "add_answer")
	rootBtn := tgbotapi.NewInlineKeyboardButtonData("Назад", "back_btn")
	row = append(row, addAnswerBtn)
	row = append(row, rootBtn)
	keyboard.InlineKeyboard = append(keyboard.InlineKeyboard, row)
	answerMenu.ReplyMarkup = keyboard

	return answerMenu, true
}

func AddAnswer(botUtil *util.BotUtil, task *model.Task) tgbotapi.MessageConfig {
	text := "Загрузите файл или напишите сообщение"
	msg := tgbotapi.NewMessage(botUtil.Message.Chat.ID, text)

	keyboard := tgbotapi.InlineKeyboardMarkup{}
	row := make([]tgbotapi.InlineKeyboardButton, 0)
	cancelBtn := tgbotapi.NewInlineKeyboardButtonData("Отмена", "cancel_btn")
	row = append(row, cancelBtn)
	keyboard.InlineKeyboard = append(keyboard.InlineKeyboard, row)
	msg.ReplyMarkup = keyboard

	botUtil.Bot.Send(msg)

	ready := false

	var messageId int
	for update := range botUtil.Updates {
		mes := update.Message
		var command string

		if update.CallbackQuery != nil {
			response := strings.Split(update.CallbackQuery.Data, " ")
			command = response[0]
		} else if mes.Text != "" {
			command = mes.Text
		}

		if mes != nil {
			messageId = mes.MessageID
			command = "answer_entered"
		}

		switch command {
		case "cancel_btn":
			text = "Отмена добавления ответа"
			msg = tgbotapi.NewMessage(botUtil.Message.Chat.ID, text)
			return msg
		case "answer_entered":
			ready = true
		default:
			msg := util.GetUnknownMessage(botUtil, command)
			botUtil.Bot.Send(msg)
		}

		if ready {
			break
		}
	}

	executor, text := controller2.GetUser(botUtil.Message.From.UserName)
	if executor == nil {
		msg = tgbotapi.NewMessage(botUtil.Message.Chat.ID, text)
		return msg
	}

	answer := &model.Answer{
		Task:      task,
		Executor:  executor,
		MessageId: messageId,
		ChatId:    botUtil.Message.Chat.ID,
		Status:    "not checked",
	}

	forward := tgbotapi.NewForward(answer.ChatId, answer.ChatId, answer.MessageId)
	botUtil.Bot.Send(forward)
	acceptingString := "Вы действительно хотите создать ответ:\n"
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
			text, _ = controller.AddAnswer(answer)
			goto LOOP
		case "no_btn":
			text = "Отмена добавления ответа"
			goto LOOP
		default:
			text = "Неизвестная команда"
			msg = tgbotapi.NewMessage(botUtil.Message.Chat.ID, text)
			botUtil.Bot.Send(msg)

			msg = util.GetAcceptingMessage(botUtil.Message, acceptingString)
			botUtil.Bot.Send(msg)
			forward = tgbotapi.NewForward(botUtil.Message.Chat.ID, botUtil.Message.Chat.ID, messageId)
			botUtil.Bot.Send(forward)
		}
	}

LOOP:
	msg = tgbotapi.NewMessage(botUtil.Message.Chat.ID, text)
	return msg
}
