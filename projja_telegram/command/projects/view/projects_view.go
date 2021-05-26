package view

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"log"
	"projja_telegram/command/current_project/view"
	"projja_telegram/command/projects/controller"
	projectsmenu "projja_telegram/command/projects/menu"
	"projja_telegram/command/root/menu"
	"projja_telegram/command/util"
	"projja_telegram/model"
	"strconv"
	"strings"
)

func SelectProject(botUtil *util.BotUtil) {
	defer func(message *util.MessageData, bot *tgbotapi.BotAPI) {
		msg := menu.GetRootMenu(message)
		bot.Send(msg)
	}(botUtil.Message, botUtil.Bot)

	projectsCount, status := controller.GetProjectsCount(botUtil.Message.From)
	if !status {
		errorText := "Не удалось получить список проектов\n" +
			"Попробуйте позже"
		msg := tgbotapi.NewMessage(botUtil.Message.Chat.ID, errorText)
		botUtil.Bot.Send(msg)
		return
	}

	page := 1

	msg, projects, status := projectsmenu.MakeProjectsMenu(botUtil.Message, page, projectsCount)
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
		case "back_btn":
			return
		case "create_project":
			page = 1
			msg = CreateProject(botUtil)
			botUtil.Bot.Send(msg)

			projectsCount, status = controller.GetProjectsCount(botUtil.Message.From)
			if !status {
				errorText := "Не удалось получить список проектов\n" +
					"Попробуйте позже"
				msg = tgbotapi.NewMessage(botUtil.Message.Chat.ID, errorText)
				botUtil.Bot.Send(msg)
				return
			}
		case "prev_page":
			page--
		case "next_page":
			page++
		default:
			msg, index, status := IsProjectId(botUtil.Message, command, projects)
			botUtil.Bot.Send(msg)
			if status {
				view.WorkWithProject(botUtil, projects[index])
			}
		}

		msg, projects, status = projectsmenu.MakeProjectsMenu(botUtil.Message, page, projectsCount)
		botUtil.Bot.Send(msg)
		if !status {
			return
		}
	}
}

func CreateProject(botUtil *util.BotUtil) tgbotapi.MessageConfig {
	text := "Введите имя нового проекта"
	msg := tgbotapi.NewMessage(botUtil.Message.Chat.ID, text)

	keyboard := tgbotapi.InlineKeyboardMarkup{}
	row := make([]tgbotapi.InlineKeyboardButton, 0)
	cancelBtn := tgbotapi.NewInlineKeyboardButtonData("Отмена", "cancel_btn")
	row = append(row, cancelBtn)
	keyboard.InlineKeyboard = append(keyboard.InlineKeyboard, row)

	msg.ReplyMarkup = keyboard

	botUtil.Bot.Send(msg)

	projectName := ""
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
		case "cancel_btn":
			text = "Отмена создания проекта"
			msg = tgbotapi.NewMessage(botUtil.Message.Chat.ID, text)
			return msg
		default:
			if command == "" {
				continue
			}
			projectName = command
		}

		if projectName != "" {
			break
		}
	}

	acceptingString := fmt.Sprintf("Вы действительно хотите создать проект с именем '%s'?", projectName)
	msg = util.GetAcceptingMessage(botUtil.Message, acceptingString)
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
			text, _ = controller.CreateNewProject(botUtil.Message, projectName)
			goto LOOP
		case "no_btn":
			text = "Отмена создания проекта"
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

func IsProjectId(message *util.MessageData, command string, projects []*model.Project) (tgbotapi.MessageConfig, int, bool) {
	id, err := strconv.Atoi(command)
	if err != nil {
		log.Println("error in casting command: ", err)
		text := "Вы ввели не номер проекта в списке, а '" + command + "'"
		msg := tgbotapi.NewMessage(message.Chat.ID, text)
		return msg, 0, false
	}
	count := len(projects)
	if id > count || id < 1 {
		log.Println(fmt.Sprintf("id not in range 1-%d", count))
		text := fmt.Sprintf("Номер проекта должен быть в интервале от 1 до %d", count)
		msg := tgbotapi.NewMessage(message.Chat.ID, text)
		return msg, id, false
	}

	text := fmt.Sprintf("Выбран проект '%s'", projects[id-1].Name)
	msg := tgbotapi.NewMessage(message.Chat.ID, text)

	return msg, id - 1, true
}
