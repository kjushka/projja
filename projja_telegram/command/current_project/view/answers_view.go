package view

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"projja_telegram/command/current_project/controller"
	"projja_telegram/command/current_project/menu"
	"projja_telegram/command/util"
	"projja_telegram/model"
)

func CheckProjectAnswers(botUtil *util.BotUtil, project *model.Project) {
	page := 1
	answers, msg, status := ShowAnswersMenu(botUtil, project, page)
	botUtil.Bot.Send(msg)
	if !status {
		return
	}

	for update := range botUtil.Updates {
		mes := update.Message
		command := ""

		if mes.Text != "" {
			command = mes.Text
		}

		switch command {
		case "Обновить данные":
			break
		case "Предыдущая страница":
			page--
		case "Следующая страница":
			page++
		case "Назад":
			return
		default:
			text, index, status := IsTaskDescription(answers, command)
			msg = tgbotapi.NewMessage(botUtil.Message.Chat.ID, text)
			botUtil.Bot.Send(msg)
			if status {
				CheckAnswer(botUtil, answers[index])
			}
		}

		answers, msg, status = ShowAnswersMenu(botUtil, project, page)
		botUtil.Bot.Send(msg)
		if !status {
			return
		}
	}
}

func ShowAnswersMenu(botUtil *util.BotUtil, project *model.Project, page int) ([]*model.Answer, tgbotapi.MessageConfig, bool) {
	answers, status := controller.GetProjectAnswers(project)
	if !status {
		errorText := "Не удалось получить список ответов\n" +
			"Попробуйте позже"
		msg := tgbotapi.NewMessage(botUtil.Message.Chat.ID, errorText)
		return nil, msg, false
	}

	count := len(answers) - (page-1)*4
	if count > 4 {
		count = 4
	}
	msg := menu.MakeProjectAnswersMenu(botUtil.Message, answers, page, count)
	return answers, msg, true
}

func IsTaskDescription(answers []*model.Answer, command string) (string, int, bool) {
	if command == "" {
		text := "Для задачи с таким описанием нет ответа"
		return text, -1, false
	}

	index := -1
	found := false
	for i, answer := range answers {
		if answer.Task.Description == command {
			found = true
			index = i
			break
		}
	}

	if !found {
		text := "Для задачи с таким описанием нет ответа"
		return text, index, found
	}

	return "", index, found
}

func CheckAnswer(botUtil *util.BotUtil, answer *model.Answer) {
	text := fmt.Sprintf("Ответ на задачу '%s'", answer.Task.Description)
	msg := tgbotapi.NewMessage(botUtil.Message.Chat.ID, text)

	row1 := make([]tgbotapi.KeyboardButton, 0)
	row2 := make([]tgbotapi.KeyboardButton, 0)
	acceptBtn := tgbotapi.NewKeyboardButton("Принять")
	declineBtn := tgbotapi.NewKeyboardButton("Отклонить")
	backBtn := tgbotapi.NewKeyboardButton("Назад")
	row1 = append(row1, acceptBtn, declineBtn)
	row2 = append(row2, backBtn)

	keyboard := tgbotapi.NewReplyKeyboard(row1, row2)
	msg.ReplyMarkup = keyboard
	botUtil.Bot.Send(msg)

	forward := tgbotapi.NewForward(botUtil.Message.Chat.ID, answer.ChatId, answer.MessageId)
	botUtil.Bot.Send(forward)

	for update := range botUtil.Updates {
		mes := update.Message
		command := ""

		if mes.Text != "" {
			command = mes.Text
		}

		switch command {
		case "Принять":
			resultText, status, closed := controller.AcceptAnswer(answer)
			msg := tgbotapi.NewMessage(botUtil.Message.Chat.ID, resultText)
			botUtil.Bot.Send(msg)
			if status {
				notification := fmt.Sprintf("Ваше решение задачи '%s' принято. ", answer.Task.Description)
				if closed {
					notification += "Задача закрыта"
				} else {
					notification += "Задача переведена на следующий этап"
				}
				msg := tgbotapi.NewMessage(answer.ChatId, notification)
				botUtil.Bot.Send(msg)
				return
			}
		case "Отклонить":
			resultText, status := controller.DeclineAnswer(answer)
			msg := tgbotapi.NewMessage(botUtil.Message.Chat.ID, resultText)
			botUtil.Bot.Send(msg)
			if status {
				notification := fmt.Sprintf("Ваше решение задачи '%s' отклонено", answer.Task.Description)
				msg := tgbotapi.NewMessage(answer.ChatId, notification)
				botUtil.Bot.Send(msg)
				return
			}
		case "Назад":
			return
		default:
			msg := util.GetUnknownMessage(botUtil)
			botUtil.Bot.Send(msg)
		}

		botUtil.Bot.Send(msg)
		botUtil.Bot.Send(forward)
	}
}
