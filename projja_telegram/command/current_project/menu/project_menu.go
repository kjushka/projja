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

	projectsMenuBtn := tgbotapi.NewInlineKeyboardButtonData("Меню выбора проектов", "projects_menu")

	if project.Status == "opened" {
		changeMembersBtn := tgbotapi.NewInlineKeyboardButtonData("Участники проекта", "members")
		changeTaskStatusesBtn := tgbotapi.NewInlineKeyboardButtonData("Статусы задач", "statuses")
		addTaskBtn := tgbotapi.NewInlineKeyboardButtonData("Создать задачу", "add_task")
		checkAnswersBtn := tgbotapi.NewInlineKeyboardButtonData("Ответы на задачи", "answers")

		var row2 []tgbotapi.InlineKeyboardButton
		var row3 []tgbotapi.InlineKeyboardButton

		row1 = append(row1, changeMembersBtn)
		row2 = append(row2, changeTaskStatusesBtn)
		row2 = append(row2, addTaskBtn)
		row3 = append(row3, checkAnswersBtn)
		row3 = append(row3, projectsMenuBtn)

		keyboard.InlineKeyboard = append(keyboard.InlineKeyboard, row1)
		keyboard.InlineKeyboard = append(keyboard.InlineKeyboard, row2)
		keyboard.InlineKeyboard = append(keyboard.InlineKeyboard, row3)
	} else {
		row1 = append(row1, projectsMenuBtn)
		keyboard.InlineKeyboard = append(keyboard.InlineKeyboard, row1)
	}

	msg.ReplyMarkup = keyboard
	return msg
}
