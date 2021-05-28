package util

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

func GetUnknownMessage(botUtil *BotUtil) tgbotapi.MessageConfig {
	text := "Пожалуйста, выберите один из вариантов"
	msg := tgbotapi.NewMessage(botUtil.Message.Chat.ID, text)
	return msg
}

func ListenForText(botUtil *BotUtil, mesText string, cancelText string) (string, bool) {
	msg := tgbotapi.NewMessage(botUtil.Message.Chat.ID, mesText)

	row := make([]tgbotapi.KeyboardButton, 0)
	cancelBtn := tgbotapi.NewKeyboardButton("Отмена")
	row = append(row, cancelBtn)

	keyboard := tgbotapi.NewReplyKeyboard(row)
	msg.ReplyMarkup = keyboard

	botUtil.Bot.Send(msg)

	resultText := ""
	for update := range botUtil.Updates {
		mes := update.Message
		command := ""

		if mes.Text != "" {
			command = mes.Text
		}

		switch command {
		case "Отмена":
			resultText = cancelText
			return resultText, false
		default:
			if command == "" {
				botUtil.Bot.Send(msg)
				continue
			}
			resultText = command
		}

		if resultText != "" {
			break
		}
	}

	return resultText, true
}
