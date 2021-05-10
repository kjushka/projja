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
	var row2 []tgbotapi.InlineKeyboardButton
	var row3 []tgbotapi.InlineKeyboardButton
	changeNameBtn := tgbotapi.NewInlineKeyboardButtonData("Сменить название проекта", "change_name")
	changeMembersBtn := tgbotapi.NewInlineKeyboardButtonData("Редактировать участников проекта", "change_members")
	changeTaskStatusesBtn := tgbotapi.NewInlineKeyboardButtonData("Редактировать статусы задач", "change_statuses")
	addTaskBtn := tgbotapi.NewInlineKeyboardButtonData("Создать задачу", "add_task")
	checkAnswersBtn := tgbotapi.NewInlineKeyboardButtonData("Проверить ответы на задачи", "check_answers")
	projectsMenuBtn := tgbotapi.NewInlineKeyboardButtonData("Вернуться в меню выбора проектов", "projects_menu")

	row1 = append(row1, changeNameBtn)
	row1 = append(row1, changeMembersBtn)
	row2 = append(row2, changeTaskStatusesBtn)
	row2 = append(row2, addTaskBtn)
	row3 = append(row3, checkAnswersBtn)
	row3 = append(row3, projectsMenuBtn)
	keyboard.InlineKeyboard = append(keyboard.InlineKeyboard, row1)
	keyboard.InlineKeyboard = append(keyboard.InlineKeyboard, row2)
	keyboard.InlineKeyboard = append(keyboard.InlineKeyboard, row3)

	msg.ReplyMarkup = keyboard
	return msg
}
