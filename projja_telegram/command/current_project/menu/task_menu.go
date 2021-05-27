package menu

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"projja_telegram/command/util"
	"projja_telegram/model"
)

func MakeTaskMenu(message *util.MessageData, task *model.Task) tgbotapi.MessageConfig {
	text := fmt.Sprintf("Работаем над задачей '%s' до %s, %s", task.Description, task.Deadline, task.Priority)
	msg := tgbotapi.NewMessage(message.Chat.ID, text)

	var row1 []tgbotapi.KeyboardButton
	var row2 []tgbotapi.KeyboardButton
	var row3 []tgbotapi.KeyboardButton
	descriptionBtn := tgbotapi.NewKeyboardButton("Изменить описание")
	deadlineBtn := tgbotapi.NewKeyboardButton("Изменить дедлайн")
	priorityBtn := tgbotapi.NewKeyboardButton("Изменить приоритет")
	executorBtn := tgbotapi.NewKeyboardButton("Изменить исполнителя")
	closeBtn := tgbotapi.NewKeyboardButton("Закрыть задачу")
	backBtn := tgbotapi.NewKeyboardButton("Назад")
	row1 = append(row1, descriptionBtn)
	row1 = append(row1, deadlineBtn)
	row2 = append(row2, priorityBtn)
	row2 = append(row2, executorBtn)
	row3 = append(row3, closeBtn)
	row3 = append(row3, backBtn)

	keyboard := tgbotapi.NewReplyKeyboard(row1, row2, row3)
	msg.ReplyMarkup = keyboard
	return msg
}
