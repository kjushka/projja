package view

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"projja_telegram/command/current_project/view"
	"projja_telegram/command/execute/controller"
	"projja_telegram/command/execute/menu"
	"projja_telegram/command/util"
	"projja_telegram/model"
	"strings"
)

func ExecuteTasks(botUtil *util.BotUtil) {
	page := 1
	tasks, msg, status := ShowExecutedTasksMenu(botUtil, page)
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
			text, index, status := view.IsTaskId(command, len(tasks), page)
			msg = tgbotapi.NewMessage(botUtil.Message.Chat.ID, text)
			botUtil.Bot.Send(msg)
			if status {
				ManageExecutorAnswers(botUtil, tasks[index])
			}
		}

		tasks, msg, status = ShowExecutedTasksMenu(botUtil, page)
		botUtil.Bot.Send(msg)
		if !status {
			return
		}
	}
}

func ShowExecutedTasksMenu(botUtil *util.BotUtil, page int) ([]*model.Task, tgbotapi.MessageConfig, bool) {
	tasks, status := controller.GetExecutedTasks(botUtil.Message.From)
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
	msg := menu.MakeExecutedTasksMenu(botUtil.Message, tasks, page, count)
	return tasks, msg, true
}

func ManageExecutorAnswers(botUtil *util.BotUtil, task *model.Task) {
	answerMenu, status := MakeAddAnswerMenu(botUtil, task)
	botUtil.Bot.Send(answerMenu)
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

		if mes != nil {
			if mes.Photo != nil {
				command = "photo"
			}
			if mes.Document != nil {
				command = "document"
			}
			if mes.Audio != nil {
				command = "audio"
			}
			if mes.Voice != nil {
				command = "voice"
			}
			if mes.VideoNote != nil {
				command = "video_note"
			}
			if mes.Video != nil {
				command = "video"
			}
		}

		switch command {
		case "back_btn":
			return
		case "photo":
			text := "its photo"
			msg := tgbotapi.NewMessage(botUtil.Message.Chat.ID, text)
			botUtil.Bot.Send(msg)
		case "document":
			text := "its document"
			msg := tgbotapi.NewMessage(botUtil.Message.Chat.ID, text)
			botUtil.Bot.Send(msg)
		case "audio":
			text := "its audio"
			msg := tgbotapi.NewMessage(botUtil.Message.Chat.ID, text)
			botUtil.Bot.Send(msg)
		case "video":
			text := "its video"
			msg := tgbotapi.NewMessage(botUtil.Message.Chat.ID, text)
			botUtil.Bot.Send(msg)
		case "video_note":
			text := "its video_note"
			msg := tgbotapi.NewMessage(botUtil.Message.Chat.ID, text)
			botUtil.Bot.Send(msg)
		case "voice":
			text := "its voice"
			msg := tgbotapi.NewMessage(botUtil.Message.Chat.ID, text)
			botUtil.Bot.Send(msg)
		default:
			msg := util.GetUnknownMessage(botUtil, command)
			botUtil.Bot.Send(msg)
		}

		answerMenu, status = MakeAddAnswerMenu(botUtil, task)
		botUtil.Bot.Send(answerMenu)
		if !status {
			return
		}
	}
}

func MakeAddAnswerMenu(botUtil *util.BotUtil, task *model.Task) (tgbotapi.MessageConfig, bool) {
	answer, status := controller.GetLastAnswer(botUtil.Message.From, task)
	if !status {
		msg := tgbotapi.NewMessage(
			botUtil.Message.Chat.ID,
			"При получении последнего решения произошла ошибка\n"+
				"Попробуйте позже ещё раз",
		)
		return msg, false
	}

	var answerAsString string
	if answer == nil {
		answerAsString = "Вы ещё не отправляли решение для этой задачи"
	} else {
		answerAsString = fmt.Sprintf("Ваш последний ответ: \n%s - %s", answer.Answer, answer.Status)
	}
	answerMenu := tgbotapi.NewMessage(botUtil.Message.Chat.ID, answerAsString)
	keyboard := tgbotapi.InlineKeyboardMarkup{}
	row := make([]tgbotapi.InlineKeyboardButton, 0)
	addAnswerBtn := tgbotapi.NewInlineKeyboardButtonData("Добавить ответ", "add_answer")
	rootBtn := tgbotapi.NewInlineKeyboardButtonData("Назад", "back_btn")
	row = append(row, addAnswerBtn)
	row = append(row, rootBtn)
	keyboard.InlineKeyboard = append(keyboard.InlineKeyboard, row)
	answerMenu.ReplyMarkup = keyboard

	return answerMenu, true
}

func AddFile() {

}
