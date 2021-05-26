package util

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"projja_telegram/model"
	"strconv"
)

type BotUtil struct {
	Message *MessageData
	Bot     *tgbotapi.BotAPI
	Updates chan tgbotapi.Update
}

type MessageData struct {
	From *tgbotapi.User
	Chat *tgbotapi.Chat
}

func TgUserToModelUser(data *MessageData) *model.User {
	name := data.From.FirstName
	if data.From.LastName != "" {
		name += " " + data.From.LastName
	}
	return &model.User{
		Name:       name,
		Username:   data.From.UserName,
		TelegramId: strconv.Itoa(data.From.ID),
		ChatId:     data.Chat.ID,
	}
}
