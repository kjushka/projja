package menu

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"projja_telegram/command/util"
	"projja_telegram/model"
)

func MakeProjectMenu(message *util.MessageData, project *model.Project) tgbotapi.MessageConfig {
	text := fmt.Sprintf("Работаем над проектом '%s'", project.Name)
	msg := tgbotapi.NewMessage(message.Chat.ID, text)

	keyboard := tgbotapi.InlineKeyboardMarkup{}

	var row1 []tgbotapi.InlineKeyboardButton
	settingsBtn := tgbotapi.NewInlineKeyboardButtonData("Настройки проекта", "settings")

	row1 = append(row1, settingsBtn)

	projectsMenuBtn := tgbotapi.NewInlineKeyboardButtonData("Назад", "back_btn")

	if project.Status == "opened" {
		addTaskBtn := tgbotapi.NewInlineKeyboardButtonData("Создать задачу", "add_task")
		checkAnswersBtn := tgbotapi.NewInlineKeyboardButtonData("Ответы на задачи", "answers")

		var row2 []tgbotapi.InlineKeyboardButton

		row1 = append(row1, addTaskBtn)
		row2 = append(row2, checkAnswersBtn)
		row2 = append(row2, projectsMenuBtn)

		keyboard.InlineKeyboard = append(keyboard.InlineKeyboard, row1)
		keyboard.InlineKeyboard = append(keyboard.InlineKeyboard, row2)
	} else {
		row1 = append(row1, projectsMenuBtn)
		keyboard.InlineKeyboard = append(keyboard.InlineKeyboard, row1)
	}

	msg.ReplyMarkup = keyboard
	return msg
}
