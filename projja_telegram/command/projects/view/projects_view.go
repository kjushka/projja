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
	"strconv"
	"strings"
)

func SelectProject(message *util.MessageData, bot *tgbotapi.BotAPI, updates tgbotapi.UpdatesChannel) {
	defer func(message *util.MessageData, bot *tgbotapi.BotAPI) {
		msg := menu.GetRootMenu(message)
		bot.Send(msg)
	}(message, bot)

	projectsCount, status := controller.GetProjectsCount(message.From)
	if !status {
		errorText := "Не удалось получить список проектов\n" +
			"Попробуйте позже"
		msg := tgbotapi.NewMessage(message.Chat.ID, errorText)
		bot.Send(msg)
		return
	}

	page := 1

	msg, projects, status := projectsmenu.MakeProjectsMenu(message, page, projectsCount)
	bot.Send(msg)
	if !status {
		return
	}

	for update := range updates {
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
		case "root":
			return
		case "create_project":
			page = 1
			msg = CreateProject(message, bot, updates)
			bot.Send(msg)

			projectsCount, status = controller.GetProjectsCount(message.From)
			if !status {
				errorText := "Не удалось получить список проектов\n" +
					"Попробуйте позже"
				msg = tgbotapi.NewMessage(message.Chat.ID, errorText)
				bot.Send(msg)
				return
			}
		case "prev_page":
			page--
		case "next_page":
			page++
		default:
			msg, index, status := IsProjectId(message, command, len(projects))
			bot.Send(msg)
			if status {
				view.WorkWithProject(message, bot, updates, projects[index])
			}
		}

		msg, projects, status = projectsmenu.MakeProjectsMenu(message, page, projectsCount)
		bot.Send(msg)
		if !status {
			return
		}
	}

	log.Println(projects)
}

func CreateProject(message *util.MessageData, bot *tgbotapi.BotAPI, updates tgbotapi.UpdatesChannel) tgbotapi.MessageConfig {
	text := "Введите имя нового проекта"
	msg := tgbotapi.NewMessage(message.Chat.ID, text)
	bot.Send(msg)

	projectName := ""
	for update := range updates {
		mes := update.Message
		if mes == nil {
			continue
		}

		projectName = mes.Text
		break
	}

	acceptingString := fmt.Sprintf("Вы действительно хотите создать проект с именем '%s'?", projectName)
	msg = util.GetAcceptingMessage(message, acceptingString)
	bot.Send(msg)

	for update := range updates {
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
		case "yes":
			text, _ = controller.CreateNewProject(message.From, projectName)
			goto LOOP
		case "no":
			text = "Отмена создания проекта"
			goto LOOP
		default:
			text = "Неизвестная команда"
			msg = tgbotapi.NewMessage(message.Chat.ID, text)
			bot.Send(msg)

			msg = util.GetAcceptingMessage(message, acceptingString)
			bot.Send(msg)
		}
	}

LOOP:
	msg = tgbotapi.NewMessage(message.Chat.ID, text)
	return msg
}

func IsProjectId(message *util.MessageData, command string, projectsCount int) (tgbotapi.MessageConfig, int, bool) {
	id, err := strconv.Atoi(command)
	if err != nil {
		log.Println("error in casting command: ", err)
		text := "Вы ввели не номер проекта в списке, а '" + command + "'"
		msg := tgbotapi.NewMessage(message.Chat.ID, text)
		return msg, 0, false
	}
	if id > projectsCount || id < 1 {
		log.Println(fmt.Sprintf("id not in range 1-%d", projectsCount))
		text := fmt.Sprintf("Номер проекта должен быть в интервале от 1 до %d", projectsCount)
		msg := tgbotapi.NewMessage(message.Chat.ID, text)
		return msg, id, false
	}

	text := fmt.Sprintf("Выбран проект под номером %d", id)
	msg := tgbotapi.NewMessage(message.Chat.ID, text)

	return msg, id - 1, true
}
