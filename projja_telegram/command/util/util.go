package util

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"projja_telegram/model"
	"strconv"
)

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
