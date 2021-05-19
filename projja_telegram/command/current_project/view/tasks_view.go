package view

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"projja_telegram/command/current_project/controller"
	"projja_telegram/command/current_project/menu"
	"projja_telegram/command/util"
	"projja_telegram/model"
	"strings"
)

func ManageProjectTasks(botUtil *util.BotUtil, project *model.Project) {
	page := 1
	_, msg, status := ShowTasksMenu(botUtil, project, page)
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

		_, msg, status = ShowTasksMenu(botUtil, project, page)
		botUtil.Bot.Send(msg)
		if !status {
			return
		}
	}
}

func ShowTasksMenu(botUtil *util.BotUtil, project *model.Project, page int) ([]*model.Task, tgbotapi.MessageConfig, bool) {
	tasks, status := controller.GetProjectTasks(project)
	if !status {
		errorText := "Не удалось получить список задач\n" +
			"Попробуйте позже"
		msg := tgbotapi.NewMessage(botUtil.Message.Chat.ID, errorText)
		return nil, msg, false
	}

	count := len(tasks) - (page-1)*10
	if count > 10 {
		count = 10
	}
	msg := menu.MakeProjectTasksMenu(botUtil.Message, project, tasks, page, count)
	return tasks, msg, true
}
