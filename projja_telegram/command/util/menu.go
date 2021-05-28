package util

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

func GetAcceptingMessage(message *MessageData, acceptingString string) tgbotapi.MessageConfig {
	msg := tgbotapi.NewMessage(message.Chat.ID, acceptingString)

	row := make([]tgbotapi.KeyboardButton, 0)
	yesBtn := tgbotapi.NewKeyboardButton("Да")
	noBtn := tgbotapi.NewKeyboardButton("Нет")
	row = append(row, yesBtn)
	row = append(row, noBtn)

	keyboard := tgbotapi.NewReplyKeyboard(row)
	msg.ReplyMarkup = keyboard

	return msg
}
