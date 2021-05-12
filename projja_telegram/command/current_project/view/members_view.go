package view

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"projja_telegram/command/current_project/controller"
	"projja_telegram/command/current_project/menu"
	"projja_telegram/command/util"
	"projja_telegram/model"
	"strings"
)

func ChangeProjectMembers(botUtil *util.BotUtil, project *model.Project) {
	page := 1
	members, msg, status := ShowMemberMenu(botUtil, project, page)
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
		case "add_member":
			msg = AddMember(botUtil, project, members)
			botUtil.Bot.Send(msg)
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

		members, msg, status = ShowMemberMenu(botUtil, project, page)
		botUtil.Bot.Send(msg)
		if !status {
			return
		}
	}
}

func ShowMemberMenu(botUtil *util.BotUtil, project *model.Project, page int) ([]*model.User, tgbotapi.MessageConfig, bool) {
	members, status := controller.GetMembers(project)
	if !status {
		errorText := "Не удалось получить список участников\n" +
			"Попробуйте позже"
		msg := tgbotapi.NewMessage(botUtil.Message.Chat.ID, errorText)
		return nil, msg, false
	}

	count := len(members) - (page-1)*10
	if count > 10 {
		count = 10
	}
	msg := menu.MakeMembersMenu(botUtil.Message, project, members, page, count)
	return members, msg, true
}

func AddMember(botUtil *util.BotUtil, project *model.Project, members []*model.User) tgbotapi.MessageConfig {
	text, cancelStatus := util.ListenForText(botUtil,
		"Введите username нового участника",
		"Отмена добавления участника",
	)
	if !cancelStatus {
		msg := tgbotapi.NewMessage(botUtil.Message.Chat.ID, text)
		return msg
	}

	memberUsername := text

	for _, member := range members {
		if member.Username == memberUsername {
			text = "Данный участник уже добавлен в проект"
			msg := tgbotapi.NewMessage(botUtil.Message.Chat.ID, text)
			return msg
		}
	}

	member, text := controller.GetUser(memberUsername)
	if member == nil {
		msg := tgbotapi.NewMessage(botUtil.Message.Chat.ID, text)
		return msg
	}

	acceptingString := fmt.Sprintf("Вы хотите добавить:\n"+
		"Имя: %s\n"+
		"Username: %s\n",
		member.Name,
		member.Username,
	)
	msg := util.GetAcceptingMessage(botUtil.Message, acceptingString)

	botUtil.Bot.Send(msg)

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
		case "yes_btn":
			text, _ = controller.AddMember(project, member)
			goto LOOP
		case "no_btn":
			text = "Отмена добавления участника"
			goto LOOP
		default:
			text = "Неизвестная команда"
			msg = tgbotapi.NewMessage(botUtil.Message.Chat.ID, text)
			botUtil.Bot.Send(msg)

			msg = util.GetAcceptingMessage(botUtil.Message, acceptingString)
			botUtil.Bot.Send(msg)
		}
	}

LOOP:
	msg = tgbotapi.NewMessage(botUtil.Message.Chat.ID, text)
	return msg
}
