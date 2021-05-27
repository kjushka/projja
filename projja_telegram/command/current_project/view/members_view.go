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
	members, msg, status := ShowMembersMenu(botUtil, project, page)
	botUtil.Bot.Send(msg)
	if !status {
		return
	}

	for update := range botUtil.Updates {
		mes := update.Message
		command := ""

		if mes.Text != "" {
			command = mes.Text
		}

		switch command {
		case "Добавить участника":
			msg = AddMember(botUtil, project, members)
			botUtil.Bot.Send(msg)
		case "Удалить участника":
			msg = RemoveMember(botUtil, project, members)
			botUtil.Bot.Send(msg)
		case "Предыдущая страница":
			page--
		case "Следующая страница":
			page++
		case "Назад":
			return
		default:
			msg = util.GetUnknownMessage(botUtil)
			botUtil.Bot.Send(msg)
		}

		members, msg, status = ShowMembersMenu(botUtil, project, page)
		botUtil.Bot.Send(msg)
		if !status {
			return
		}
	}
}

func ShowMembersMenu(botUtil *util.BotUtil, project *model.Project, page int) ([]*model.User, tgbotapi.MessageConfig, bool) {
	members, status := controller.GetMembers(project)
	if !status {
		errorText := "Не удалось получить список участников\n" +
			"Попробуйте позже"
		msg := tgbotapi.NewMessage(botUtil.Message.Chat.ID, errorText)
		return nil, msg, false
	}

	count := len(members) - (page-1)*4
	if count > 4 {
		count = 4
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
		"Username: %s\n"+
		"Навыки: %s\n",
		member.Name,
		member.Username,
		strings.Join(member.Skills, ", "),
	)
	msg := util.GetAcceptingMessage(botUtil.Message, acceptingString)

	botUtil.Bot.Send(msg)

	for update := range botUtil.Updates {
		mes := update.Message
		var command string

		if update.CallbackQuery != nil {
			response := strings.Split(update.CallbackQuery.Data, " ")
			command = response[0]
		} else if mes.IsCommand() {
			command = mes.Command()
		} else if mes.Text != "" {
			command = mes.Text
		}

		switch command {
		case "Да":
			text, _ = controller.AddMember(project, member)
			goto LOOP
		case "Нет":
			text = "Отмена добавления участника"
			goto LOOP
		default:
			text = "Пожалуйста, выберите один из вариантов"
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

func RemoveMember(botUtil *util.BotUtil, project *model.Project, members []*model.User) tgbotapi.MessageConfig {
	page := 1
	count := len(members) - (page-1)*4
	if count > 4 {
		count = 4
	}
	msg := menu.MakeMembersRemovingMenu(botUtil.Message, project, members, page, count)
	botUtil.Bot.Send(msg)

	memberIndex := -1

	for update := range botUtil.Updates {
		mes := update.Message
		command := ""

		exit := false

		if mes.Text != "" {
			command = mes.Text
		}

		switch command {
		case "Отмена":
			text := "Отмена удаления участника"
			msg = tgbotapi.NewMessage(botUtil.Message.Chat.ID, text)
			return msg
		case "Предыдущая страница":
			page--
		case "Следующая страница":
			page++
		default:
			text, index, status := IsMemberName(members, command)
			memberIndex = index
			if !status {
				msg := tgbotapi.NewMessage(botUtil.Message.Chat.ID, text)
				botUtil.Bot.Send(msg)
			} else {
				exit = true
			}
		}

		if exit && memberIndex != -1 {
			break
		}
		botUtil.Bot.Send(msg)
	}

	member := members[memberIndex]

	acceptingString := fmt.Sprintf("Вы хотите удалить участника '%s'", member.Name)

	msg = util.GetAcceptingMessage(botUtil.Message, acceptingString)
	botUtil.Bot.Send(msg)

	var text string
	for update := range botUtil.Updates {
		mes := update.Message
		var command string

		if update.CallbackQuery != nil {
			response := strings.Split(update.CallbackQuery.Data, " ")
			command = response[0]
		} else if mes.IsCommand() {
			command = mes.Command()
		} else if mes.Text != "" {
			command = mes.Text
		}

		switch command {
		case "Да":
			text, _ = controller.RemoveMember(project, member)
			goto LOOP
		case "Нет":
			text = "Отмена удаления участника"
			goto LOOP
		default:
			text = "Пожалуйста, выберите один из вариантов"
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

func IsMemberName(members []*model.User, command string) (string, int, bool) {
	if command == "" {
		text := "Участника с таким username не существует"
		return text, -1, false
	}

	index := -1
	found := false
	for i, m := range members {
		if m.Username == command {
			found = true
			index = i
			break
		}
	}

	if !found {
		text := "Участника с таким username не существует"
		return text, index, found
	}

	return "", index, found
}
