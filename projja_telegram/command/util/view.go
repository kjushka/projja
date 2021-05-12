package util

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

func GetUnknownMessage(botUtil *BotUtil, command string) tgbotapi.MessageConfig {
	text := fmt.Sprintf("Я не знаю команды '%s'", command)
	msg := tgbotapi.NewMessage(botUtil.Message.Chat.ID, text)
	return msg
}
