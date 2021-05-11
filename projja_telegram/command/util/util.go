package util

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"projja_telegram/model"
	"strconv"
)

type BotUtil struct {
	Message *MessageData
	Bot     *tgbotapi.BotAPI
	Updates tgbotapi.UpdatesChannel
}

type MessageData struct {
	From *tgbotapi.User
	Chat *tgbotapi.Chat
}

func MessageToMessageData(message *tgbotapi.Message) *MessageData {
	return &MessageData{
		From: message.From,
		Chat: message.Chat,
	}
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
