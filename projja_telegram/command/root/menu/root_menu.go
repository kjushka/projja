package menu

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"projja_telegram/command/util"
)

func GetRootMenu(message *util.MessageData) tgbotapi.MessageConfig {
	text := fmt.Sprintf("%s, что вы хотите сделать?\n", message.From.UserName)
	msg := tgbotapi.NewMessage(message.Chat.ID, text)

	var row1 []tgbotapi.KeyboardButton
	var row2 []tgbotapi.KeyboardButton
	skillsBtn := tgbotapi.NewKeyboardButton("Изменить навыки")
	updateBtn := tgbotapi.NewKeyboardButton("Обновить данные профиля")
	projectsManageBtn := tgbotapi.NewKeyboardButton("Управлять проектами")
	tasksBtn := tgbotapi.NewKeyboardButton("Ваши задачи")

	row1 = append(row1, skillsBtn)
	row1 = append(row1, updateBtn)
	row2 = append(row2, projectsManageBtn)
	row2 = append(row2, tasksBtn)

	keyboard := tgbotapi.NewReplyKeyboard(row1, row2)

	msg.ReplyMarkup = keyboard
	return msg
}
