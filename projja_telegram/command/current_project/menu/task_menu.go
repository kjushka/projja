package menu

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"projja_telegram/command/util"
	"projja_telegram/model"
)

func MakeTaskMenu(message *util.MessageData, task *model.Task) tgbotapi.MessageConfig {
	text := fmt.Sprintf("Работаем над задачей '%s до %s, %s'", task.Description, task.Deadline, task.Priority)
	msg := tgbotapi.NewMessage(message.Chat.ID, text)

	keyboard := tgbotapi.InlineKeyboardMarkup{}

	var row1 []tgbotapi.InlineKeyboardButton
	var row2 []tgbotapi.InlineKeyboardButton
	var row3 []tgbotapi.InlineKeyboardButton
	descriptionBtn := tgbotapi.NewInlineKeyboardButtonData("Изменить описание", "description")
	deadlineBtn := tgbotapi.NewInlineKeyboardButtonData("Изменить дедлайн", "deadline")
	priorityBtn := tgbotapi.NewInlineKeyboardButtonData("Изменить приоритет", "priority")
	executorBtn := tgbotapi.NewInlineKeyboardButtonData("Изменить исполнителя", "executor")
	closeBtn := tgbotapi.NewInlineKeyboardButtonData("Закрыть задачу", "close_task")
	backBtn := tgbotapi.NewInlineKeyboardButtonData("Назад", "back_btn")
	row1 = append(row1, descriptionBtn)
	row1 = append(row1, deadlineBtn)
	row2 = append(row2, priorityBtn)
	row2 = append(row2, executorBtn)
	row3 = append(row3, closeBtn)
	row3 = append(row3, backBtn)
	keyboard.InlineKeyboard = append(keyboard.InlineKeyboard, row1)
	keyboard.InlineKeyboard = append(keyboard.InlineKeyboard, row2)
	keyboard.InlineKeyboard = append(keyboard.InlineKeyboard, row3)

	msg.ReplyMarkup = keyboard
	return msg
}
