package menu

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"projja_telegram/command/util"
	"projja_telegram/model"
)

func MakeSettingsMenu(message *util.MessageData, project *model.Project) tgbotapi.MessageConfig {
	text := fmt.Sprintf("Настройки проекта: '%s'", project.Name)
	msg := tgbotapi.NewMessage(message.Chat.ID, text)

	keyboard := tgbotapi.InlineKeyboardMarkup{}

	row := make([]tgbotapi.InlineKeyboardButton, 0)
	changeNameBtn := tgbotapi.NewInlineKeyboardButtonData("Сменить название", "change_name")
	changeStatusBtn := tgbotapi.NewInlineKeyboardButtonData("Открыть/закрыть проект", "change_status")

	row2 := make([]tgbotapi.InlineKeyboardButton, 0)
	projectMenuBtn := tgbotapi.NewInlineKeyboardButtonData("Меню управления проектом", "project_menu")

	row = append(row, changeNameBtn)
	row = append(row, changeStatusBtn)
	row2 = append(row2, projectMenuBtn)

	keyboard.InlineKeyboard = append(keyboard.InlineKeyboard, row)
	keyboard.InlineKeyboard = append(keyboard.InlineKeyboard, row2)

	msg.ReplyMarkup = keyboard

	return msg
}
