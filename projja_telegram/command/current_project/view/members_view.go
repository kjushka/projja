package view

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"log"
	"projja_telegram/command/current_project/controller"
	"projja_telegram/command/current_project/menu"
	"projja_telegram/command/util"
	"projja_telegram/model"
	"strconv"
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
		case "add_member":
			msg = AddMember(botUtil, project, members)
			botUtil.Bot.Send(msg)
		case "remove_member":
			msg = RemoveMember(botUtil, project, members)
			botUtil.Bot.Send(msg)
		case "prev_page":
			page--
		case "next_page":
			page++
		case "back_btn":
			return
		default:
			msg = util.GetUnknownMessage(botUtil, command)
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

func RemoveMember(botUtil *util.BotUtil, project *model.Project, members []*model.User) tgbotapi.MessageConfig {
	page := 1
	count := len(members) - (page-1)*10
	if count > 10 {
		count = 10
	}
	msg := menu.MakeMembersRemovingMenu(botUtil.Message, project, members, page, count)
	botUtil.Bot.Send(msg)

	memberIndex := -1

	for update := range botUtil.Updates {
		mes := update.Message
		var command string

		exit := false

		if update.CallbackQuery != nil {
			response := strings.Split(update.CallbackQuery.Data, " ")
			command = response[0]
		} else if mes.IsCommand() {
			command = mes.Command()
		} else if mes.Text != "" {
			command = mes.Text
		}

		switch command {
		case "cancel_btn":
			text := "Отмена удаления участника"
			msg = tgbotapi.NewMessage(botUtil.Message.Chat.ID, text)
			return msg
		case "prev_page":
			page--
		case "next_page":
			page++
		default:
			text, index, status := IsMemberId(command, len(members), page)
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
		case "yes_btn":
			text, _ = controller.RemoveMember(project, member)
			goto LOOP
		case "no_btn":
			text = "Отмена удаления участника"
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

func IsMemberId(command string, count int, page int) (string, int, bool) {
	id, err := strconv.Atoi(command)
	if err != nil {
		log.Println("error in casting command: ", err)
		text := "Вы ввели не номер участника в списке, а '" + command + "'"
		return text, -1, false
	}
	if id > count || id < 1 {
		log.Println(fmt.Sprintf("id not in range 1-%d", count))
		text := fmt.Sprintf("Номер участника должен быть в интервале от 1 до %d", count)
		return text, -1, false
	}

	id = (page-1)*10 + id

	return "", id - 1, true
}
