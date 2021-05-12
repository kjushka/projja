package view

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"projja_telegram/command/current_project/controller"
	"projja_telegram/command/current_project/menu"
	"projja_telegram/command/util"
	"projja_telegram/model"
	"strings"
)

func ChangeProjectMembers(botUtil *util.BotUtil, project *model.Project) {
	page := 1
	msg, status := ShowMemberMenu(botUtil, project, page)
	botUtil.Bot.Send(msg)
	if !status {
		return
	}

	for update := range botUtil.Updates {
		mes := update.Message
		var command string

		if update.CallbackQuery != nil {
			response := strings.Split(update.CallbackQuery.Data, " ")
			command = response[0]

			mes = update.CallbackQuery.Message
			mes.From = update.CallbackQuery.From
		} else if mes.IsCommand() {
			command = mes.Command()
		} else if mes.Text != "" {
			command = mes.Text
		}

		switch command {
		case "prev_page":
			page--
		case "next_page":
			page++
		case "project_menu":
			return
		default:
			msg = util.GetUnknownMessage(botUtil, command)
			botUtil.Bot.Send(msg)
		}

		msg, status = ShowMemberMenu(botUtil, project, page)
		botUtil.Bot.Send(msg)
		if !status {
			return
		}
	}
}

func ShowMemberMenu(botUtil *util.BotUtil, project *model.Project, page int) (tgbotapi.MessageConfig, bool) {
	members, status := controller.GetMembers(project)
	if !status {
		errorText := "Не удалось получить список участников\n" +
			"Попробуйте позже"
		msg := tgbotapi.NewMessage(botUtil.Message.Chat.ID, errorText)
		return msg, false
	}

	count := len(members) - (page-1)*10
	if count > 10 {
		count = 10
	}
	msg := menu.MakeMembersMenu(botUtil.Message, project, members, page, count)
	return msg, true
}
