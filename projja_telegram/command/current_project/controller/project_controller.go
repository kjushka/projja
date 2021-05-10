package controller

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"projja_telegram/command/projects/menu"
	"projja_telegram/command/util"
	"projja_telegram/model"
)

func WorkWithProject(message *util.MessageData,
	bot *tgbotapi.BotAPI,
	updates tgbotapi.UpdatesChannel,
	project *model.Project) {

	defer func(message *util.MessageData, bot *tgbotapi.BotAPI) {
		msg, _, _ := menu.MakeProjectsMenu(message, 1, 10)
		bot.Send(msg)
	}(message, bot)
	return
}
