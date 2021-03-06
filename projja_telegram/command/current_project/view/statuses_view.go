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
)

func ChangeProjectStatuses(botUtil *util.BotUtil, project *model.Project) {
	page := 1
	taskStatuses, msg, status := ShowTaskStatusesMenu(botUtil, project, page)
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
		case "Добавить статус":
			msg = CreateTaskStatus(botUtil, project, taskStatuses)
			botUtil.Bot.Send(msg)
		case "Удалить статус":
			msg = RemoveTaskStatus(botUtil, project, taskStatuses)
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

		taskStatuses, msg, status = ShowTaskStatusesMenu(botUtil, project, page)
		botUtil.Bot.Send(msg)
		if !status {
			return
		}
	}
}

func ShowTaskStatusesMenu(botUtil *util.BotUtil, project *model.Project, page int) ([]*model.TaskStatus, tgbotapi.MessageConfig, bool) {
	taskStatuses, status := controller.GetStatuses(project)
	if !status {
		errorText := "Не удалось получить список статусов задач\n" +
			"Попробуйте позже"
		msg := tgbotapi.NewMessage(botUtil.Message.Chat.ID, errorText)
		return nil, msg, false
	}

	count := len(taskStatuses) - (page-1)*4
	if count > 4 {
		count = 4
	}
	msg := menu.MakeTaskStatusesMenu(botUtil.Message, project, taskStatuses, page, count)
	return taskStatuses, msg, true
}

func CreateTaskStatus(botUtil *util.BotUtil, project *model.Project, taskStatuses []*model.TaskStatus) tgbotapi.MessageConfig {
	text, cancelStatus := util.ListenForText(botUtil,
		"Введите название нового статуса",
		"Отмена создания статуса задач",
	)
	if !cancelStatus {
		msg := tgbotapi.NewMessage(botUtil.Message.Chat.ID, text)
		return msg
	}

	newTaskStatusStatus := text

	for _, taskStatus := range taskStatuses {
		if taskStatus.Status == newTaskStatusStatus {
			text = "Данный статус уже добавлен в проект"
			msg := tgbotapi.NewMessage(botUtil.Message.Chat.ID, text)
			return msg
		}
	}

	text, cancelStatus = listenForLevel(botUtil, taskStatuses)

	if !cancelStatus {
		msg := tgbotapi.NewMessage(botUtil.Message.Chat.ID, text)
		return msg
	}

	newTaskStatusLevel, err := strconv.Atoi(text)
	if err != nil {
		log.Println("error in casting task status level: ", err)
		errorText := "Во время создания статуса произошла ошибка\nПопробуйте позже"
		msg := tgbotapi.NewMessage(botUtil.Message.Chat.ID, errorText)
		return msg
	}

	newTaskStatus := &model.TaskStatus{
		Status: newTaskStatusStatus,
		Level:  newTaskStatusLevel,
	}

	acceptingString := fmt.Sprintf("Вы хотите создать:\n"+
		"Статус: %s\n"+
		"Level: %d\n",
		newTaskStatus.Status,
		newTaskStatus.Level,
	)
	msg := util.GetAcceptingMessage(botUtil.Message, acceptingString)

	botUtil.Bot.Send(msg)

	for update := range botUtil.Updates {
		mes := update.Message
		command := ""

		if mes.Text != "" {
			command = mes.Text
		}

		switch command {
		case "Да":
			text, _ = controller.CreateTaskStatus(project, newTaskStatus)
			goto LOOP
		case "Нет":
			text = "Отмена создания статуса задач"
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

func listenForLevel(botUtil *util.BotUtil, taskStatuses []*model.TaskStatus) (string, bool) {
	mesText := "Выберите уровень для нового статуса"
	msg := tgbotapi.NewMessage(botUtil.Message.Chat.ID, mesText)

	rows := make([][]tgbotapi.KeyboardButton, 0)
	i := 0
	for i < len(taskStatuses)+1 {
		statusesRow := make([]tgbotapi.KeyboardButton, 0)
		firstRowStatusBtn := tgbotapi.NewKeyboardButton(strconv.Itoa(i + 1))
		statusesRow = append(statusesRow, firstRowStatusBtn)
		i++

		if i != len(taskStatuses)+1 {
			secondRowStatusBtn := tgbotapi.NewKeyboardButton(strconv.Itoa(i + 1))
			statusesRow = append(statusesRow, secondRowStatusBtn)
			i++
		}

		rows = append(rows, statusesRow)
	}

	row := make([]tgbotapi.KeyboardButton, 0)
	cancelBtn := tgbotapi.NewKeyboardButton("Отмена")
	row = append(row, cancelBtn)
	rows = append(rows, row)

	keyboard := tgbotapi.NewReplyKeyboard(rows...)
	msg.ReplyMarkup = keyboard

	botUtil.Bot.Send(msg)

	cancelText := "Отмена создания статуса задач"

	resultText := ""
	for update := range botUtil.Updates {
		mes := update.Message
		command := ""

		if mes.Text != "" {
			command = mes.Text
		}

		switch command {
		case "Отмена":
			resultText = cancelText
			return resultText, false
		default:
			if command == "" {
				botUtil.Bot.Send(msg)
				continue
			}

			resultLevel, err := strconv.Atoi(command)
			if err != nil {
				text := fmt.Sprintf("Вы ввели не уровень, а '%s'", command)
				errorMsg := tgbotapi.NewMessage(botUtil.Message.Chat.ID, text)
				botUtil.Bot.Send(errorMsg)
				botUtil.Bot.Send(msg)
				continue
			}

			if resultLevel < 1 || resultLevel > i {
				text := "Вы ввели неверный номер уровеня"
				errorMsg := tgbotapi.NewMessage(botUtil.Message.Chat.ID, text)
				botUtil.Bot.Send(errorMsg)
				botUtil.Bot.Send(msg)
				continue
			}

			resultText = command
		}

		if resultText != "" {
			break
		}
	}

	return resultText, true
}

func RemoveTaskStatus(botUtil *util.BotUtil, project *model.Project, taskStatuses []*model.TaskStatus) tgbotapi.MessageConfig {
	page := 1
	count := len(taskStatuses) - (page-1)*4
	if count > 4 {
		count = 4
	}
	msg := menu.MakeTaskStatusesRemovingMenu(botUtil.Message, taskStatuses, page, count)
	botUtil.Bot.Send(msg)

	statusIndex := -1

	for update := range botUtil.Updates {
		mes := update.Message
		command := ""

		exit := false

		if mes.Text != "" {
			command = mes.Text
		}

		switch command {
		case "Отмена":
			text := "Отмена удаления статуса задач"
			msg = tgbotapi.NewMessage(botUtil.Message.Chat.ID, text)
			return msg
		case "Предыдущая страница":
			page--
		case "Следующая страница":
			page++
		default:
			text, index, status := IsTaskStatusStatus(taskStatuses, command)
			statusIndex = index
			if !status {
				msg := tgbotapi.NewMessage(botUtil.Message.Chat.ID, text)
				botUtil.Bot.Send(msg)
			} else {
				exit = true
			}
		}

		if exit && statusIndex != -1 {
			break
		}
		botUtil.Bot.Send(msg)
	}

	taskStatus := taskStatuses[statusIndex]

	acceptingString := fmt.Sprintf("Вы хотите удалить статус задач '%s'", taskStatus.Status)

	msg = util.GetAcceptingMessage(botUtil.Message, acceptingString)
	botUtil.Bot.Send(msg)

	var text string
	for update := range botUtil.Updates {
		mes := update.Message
		command := ""

		if mes.Text != "" {
			command = mes.Text
		}

		switch command {
		case "Да":
			text, _ = controller.RemoveTaskStatus(project, taskStatus)
			goto LOOP
		case "Нет":
			text = "Отмена удаления статуса задач"
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

func IsTaskStatusStatus(statuses []*model.TaskStatus, command string) (string, int, bool) {
	if command == "" {
		text := "Статуса с таким названием не существует"
		return text, -1, false
	}

	index := -1
	found := false
	for i, s := range statuses {
		if s.Status == command {
			found = true
			index = i
			break
		}
	}

	if !found {
		text := "Статуса с таким названием не существует"
		return text, index, found
	}

	return "", index, found
}
