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
	projectsManageBtn := tgbotapi.NewInlineKeyboardButtonData("Управлять проектами", "project_management")
	updateBtn := tgbotapi.NewInlineKeyboardButtonData("Обновить данные профиля", "update_data")

	row1 = append(row1, skillsBtn)
	row1 = append(row1, projectsManageBtn)
	row2 = append(row2, updateBtn)
	keyboard.InlineKeyboard = append(keyboard.InlineKeyboard, row2)
	keyboard.InlineKeyboard = append(keyboard.InlineKeyboard, row1)

	msg.ReplyMarkup = keyboard
	return msg
}
