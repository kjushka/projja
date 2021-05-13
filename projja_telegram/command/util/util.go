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

func TgUserToModelUser(tgUser *tgbotapi.User) *model.User {
	name := tgUser.FirstName
	if tgUser.LastName != "" {
		name += " " + tgUser.LastName
	}
	return &model.User{
		Name:       name,
		Username:   tgUser.UserName,
		TelegramId: strconv.Itoa(tgUser.ID),
	}
}
