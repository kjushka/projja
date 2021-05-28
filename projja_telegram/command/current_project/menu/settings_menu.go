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

	row1 := make([]tgbotapi.KeyboardButton, 0)
	changeNameBtn := tgbotapi.NewKeyboardButton("Сменить название")
	changeStatusBtn := tgbotapi.NewKeyboardButton("Открыть/закрыть проект")

	row2 := make([]tgbotapi.KeyboardButton, 0)
	changeMembersBtn := tgbotapi.NewKeyboardButton("Участники проекта")
	changeTaskStatusesBtn := tgbotapi.NewKeyboardButton("Статусы задач")

	row3 := make([]tgbotapi.KeyboardButton, 0)
	projectMenuBtn := tgbotapi.NewKeyboardButton("Назад")

	row1 = append(row1, changeNameBtn)
	row1 = append(row1, changeStatusBtn)
	row2 = append(row2, changeMembersBtn)
	row2 = append(row2, changeTaskStatusesBtn)
	row3 = append(row3, projectMenuBtn)

	keyboard := tgbotapi.NewReplyKeyboard(row1, row2, row3)
	msg.ReplyMarkup = keyboard

	return msg
}
