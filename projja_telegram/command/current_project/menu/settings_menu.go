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
	changeMembersBtn := tgbotapi.NewInlineKeyboardButtonData("Участники проекта", "change_members")
	changeTaskStatusesBtn := tgbotapi.NewInlineKeyboardButtonData("Статусы задач", "change_statuses")

	row3 := make([]tgbotapi.InlineKeyboardButton, 0)
	projectMenuBtn := tgbotapi.NewInlineKeyboardButtonData("Назад", "back_btn")

	row = append(row, changeNameBtn)
	row = append(row, changeStatusBtn)
	row2 = append(row2, changeMembersBtn)
	row2 = append(row2, changeTaskStatusesBtn)
	row3 = append(row3, projectMenuBtn)

	keyboard.InlineKeyboard = append(keyboard.InlineKeyboard, row)
	keyboard.InlineKeyboard = append(keyboard.InlineKeyboard, row2)
	keyboard.InlineKeyboard = append(keyboard.InlineKeyboard, row3)

	msg.ReplyMarkup = keyboard

	return msg
}
