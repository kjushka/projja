package util

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

func GetAcceptingMessage(message *MessageData, acceptingString string) tgbotapi.MessageConfig {
	msg := tgbotapi.NewMessage(message.Chat.ID, acceptingString)

	keyboard := tgbotapi.InlineKeyboardMarkup{}

	var row []tgbotapi.InlineKeyboardButton
	yesBtn := tgbotapi.NewInlineKeyboardButtonData("Да", "yes_btn")
	noBtn := tgbotapi.NewInlineKeyboardButtonData("Нет", "no_btn")
	row = append(row, yesBtn)
	row = append(row, noBtn)
	keyboard.InlineKeyboard = append(keyboard.InlineKeyboard, row)

	msg.ReplyMarkup = keyboard

	return msg
}
