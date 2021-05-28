package view

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"projja_telegram/command/current_project/view"
	"projja_telegram/command/projects/controller"
	projectsmenu "projja_telegram/command/projects/menu"
	"projja_telegram/command/util"
	"projja_telegram/model"
)

func SelectProject(botUtil *util.BotUtil) {
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
		command := ""

		if mes.Text != "" {
			command = mes.Text
		}

		switch command {
		case "Назад":
			return
		case "Создать новый проект":
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
		case "Предыдущая страница":
			page--
		case "Следующая страница":
			page++
		default:
			msg, index, status := IsProjectName(botUtil.Message, command, projects)
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

	row := make([]tgbotapi.KeyboardButton, 0)
	cancelBtn := tgbotapi.NewKeyboardButton("Отмена")
	row = append(row, cancelBtn)

	keyboard := tgbotapi.NewReplyKeyboard(row)
	msg.ReplyMarkup = keyboard

	botUtil.Bot.Send(msg)

	projectName := ""
	for update := range botUtil.Updates {
		mes := update.Message
		command := ""

		if mes.Text != "" {
			command = mes.Text
		}

		switch command {
		case "Отмена":
			text = "Отмена создания проекта"
			msg = tgbotapi.NewMessage(botUtil.Message.Chat.ID, text)
			return msg
		default:
			if command == "" {
				botUtil.Bot.Send(msg)
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
		command := ""

		if mes.Text != "" {
			command = mes.Text
		}

		switch command {
		case "Да":
			text, _ = controller.CreateNewProject(botUtil.Message, projectName)
			goto LOOP
		case "Нет":
			text = "Отмена создания проекта"
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

func IsProjectName(message *util.MessageData, command string, projects []*model.Project) (tgbotapi.MessageConfig, int, bool) {
	if command == "" {
		text := "Проекта с таким названием не существует"
		msg := tgbotapi.NewMessage(message.Chat.ID, text)
		return msg, -1, false
	}

	index := -1
	found := false
	for i, p := range projects {
		if p.Name == command {
			found = true
			index = i
			break
		}
	}

	if !found {
		text := "Проекта с таким названием не существует"
		msg := tgbotapi.NewMessage(message.Chat.ID, text)
		return msg, index, found
	}

	text := fmt.Sprintf("Выбран проект '%s'", projects[index].Name)
	msg := tgbotapi.NewMessage(message.Chat.ID, text)

	return msg, index, found
}
