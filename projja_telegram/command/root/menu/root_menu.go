package menu

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"projja_telegram/command/util"
)

func GetRootMenu(message *util.MessageData) tgbotapi.MessageConfig {
	text := fmt.Sprintf("%s, что вы хотите сделать?\n", message.From.UserName)
	msg := tgbotapi.NewMessage(message.Chat.ID, text)

	keyboard := tgbotapi.InlineKeyboardMarkup{}

	var row1 []tgbotapi.InlineKeyboardButton
	var row2 []tgbotapi.InlineKeyboardButton
	skillsBtn := tgbotapi.NewInlineKeyboardButtonData("Изменить навыки", "set_skills")
	updateBtn := tgbotapi.NewInlineKeyboardButtonData("Обновить данные профиля", "update_data")
	projectsManageBtn := tgbotapi.NewInlineKeyboardButtonData("Управлять проектами", "project_management")
	tasksBtn := tgbotapi.NewInlineKeyboardButtonData("Ваши задачи", "check_tasks")

	row1 = append(row1, skillsBtn)
	row1 = append(row1, updateBtn)
	row2 = append(row2, projectsManageBtn)
	row2 = append(row2, tasksBtn)
	keyboard.InlineKeyboard = append(keyboard.InlineKeyboard, row1)
	keyboard.InlineKeyboard = append(keyboard.InlineKeyboard, row2)

	msg.ReplyMarkup = keyboard
	return msg
}
