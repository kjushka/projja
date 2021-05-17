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

func ChangeProjectStatuses(botUtil *util.BotUtil, project *model.Project) {
	page := 1
	taskStatuses, msg, status := ShowTaskStatusesMenu(botUtil, project, page)
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
		case "add_status":
			msg = CreateTaskStatus(botUtil, project, taskStatuses)
			botUtil.Bot.Send(msg)
		case "remove_status":
			msg = RemoveTaskStatus(botUtil, project, taskStatuses)
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

	count := len(taskStatuses) - (page-1)*10
	if count > 10 {
		count = 10
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
			text, _ = controller.CreateTaskStatus(project, newTaskStatus)
			goto LOOP
		case "no_btn":
			text = "Отмена создания статуса задач"
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

func listenForLevel(botUtil *util.BotUtil, taskStatuses []*model.TaskStatus) (string, bool) {
	mesText := "Выберите уровень для нового статуса"
	msg := tgbotapi.NewMessage(botUtil.Message.Chat.ID, mesText)

	keyboard := tgbotapi.InlineKeyboardMarkup{}
	row := make([]tgbotapi.InlineKeyboardButton, 0)

	i := 0
	for i < len(taskStatuses)+1 {
		membersRow := make([]tgbotapi.InlineKeyboardButton, 0)
		firstRowMemberBtn := tgbotapi.NewInlineKeyboardButtonData(strconv.Itoa(i+1), strconv.Itoa(i+1))
		membersRow = append(membersRow, firstRowMemberBtn)
		i++

		if i != len(taskStatuses)+1 {
			secondRowMemberBtn := tgbotapi.NewInlineKeyboardButtonData(strconv.Itoa(i+1), strconv.Itoa(i+1))
			membersRow = append(membersRow, secondRowMemberBtn)
			i++
		}

		keyboard.InlineKeyboard = append(keyboard.InlineKeyboard, membersRow)
	}

	cancelBtn := tgbotapi.NewInlineKeyboardButtonData("Отмена", "cancel_btn")
	row = append(row, cancelBtn)
	keyboard.InlineKeyboard = append(keyboard.InlineKeyboard, row)

	msg.ReplyMarkup = keyboard

	botUtil.Bot.Send(msg)

	cancelText := "Отмена создания статуса задач"

	resultText := ""
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
			resultText = cancelText
			return resultText, false
		default:
			if command == "" {
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
	count := len(taskStatuses) - (page-1)*10
	if count > 10 {
		count = 10
	}
	msg := menu.MakeTaskStatusesRemovingMenu(botUtil.Message, taskStatuses, page, count)
	botUtil.Bot.Send(msg)

	statusIndex := -1

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
			text := "Отмена удаления статуса задач"
			msg = tgbotapi.NewMessage(botUtil.Message.Chat.ID, text)
			return msg
		case "prev_page":
			page--
		case "next_page":
			page++
		default:
			text, index, status := IsTaskStatusIndex(command, len(taskStatuses), page)
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
			text, _ = controller.RemoveTaskStatus(project, taskStatus)
			goto LOOP
		case "no_btn":
			text = "Отмена удаления статуса задач"
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

func IsTaskStatusIndex(command string, count int, page int) (string, int, bool) {
	id, err := strconv.Atoi(command)
	if err != nil {
		log.Println("error in casting command: ", err)
		text := "Вы ввели не номер статуса задач в списке, а '" + command + "'"
		return text, -1, false
	}
	if id > count || id < 1 {
		log.Println(fmt.Sprintf("id not in range 1-%d", count))
		text := fmt.Sprintf("Номер статуса задач должен быть в интервале от 1 до %d", count)
		return text, -1, false
	}

	id = (page-1)*10 + id

	return "", id - 1, true
}
