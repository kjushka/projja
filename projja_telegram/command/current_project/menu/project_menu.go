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

	rows := make([][]tgbotapi.KeyboardButton, 0)
	row1 := make([]tgbotapi.KeyboardButton, 0)
	settingsBtn := tgbotapi.NewKeyboardButton("Настройки проекта")

	row1 = append(row1, settingsBtn)

	projectsMenuBtn := tgbotapi.NewKeyboardButton("Назад")

	if project.Status == "opened" {
		addTaskBtn := tgbotapi.NewKeyboardButton("Управление задачами")
		checkAnswersBtn := tgbotapi.NewKeyboardButton("Ответы на задачи")

		row2 := make([]tgbotapi.KeyboardButton, 0)

		row1 = append(row1, addTaskBtn)
		row2 = append(row2, checkAnswersBtn)
		row2 = append(row2, projectsMenuBtn)

		rows = append(rows, row1, row2)
	} else {
		row1 = append(row1, projectsMenuBtn)
		rows = append(rows, row1)
	}

	keyboard := tgbotapi.NewReplyKeyboard(rows...)
	msg.ReplyMarkup = keyboard
	return msg
}
