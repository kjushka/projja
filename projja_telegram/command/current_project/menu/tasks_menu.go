package menu

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"math"
	"projja_telegram/command/util"
	"projja_telegram/model"
	"strconv"
	"strings"
)

func MakeProjectTasksMenu(
	message *util.MessageData,
	project *model.Project,
	tasks []*model.Task,
	page int,
	count int,
) tgbotapi.MessageConfig {
	msg := tgbotapi.MessageConfig{}
	if len(tasks) != 0 {
		textStrings := make([]string, len(tasks))
		for i, task := range tasks {
			textStrings[i] = fmt.Sprintf("%d. %s", i+1, task.Description)
		}
		text := fmt.Sprintf(
			"Задачи проекта '%s':\n%s\n",
			project.Name,
			strings.Join(textStrings, "\n"),
		)
		msg = tgbotapi.NewMessage(message.Chat.ID, text)
	} else {
		text := "Вы ещё не создали ни одной задачи"
		msg = tgbotapi.NewMessage(message.Chat.ID, text)
	}

	keyboard := tgbotapi.InlineKeyboardMarkup{}

	if len(tasks) != 0 {
		pagesCount := int(math.Ceil(float64(count) / 10.0))
		prevNextBntRow := make([]tgbotapi.InlineKeyboardButton, 0)
		if page > 1 {
			prevBnt := tgbotapi.NewInlineKeyboardButtonData("Предыдущая страница", "prev_page")
			prevNextBntRow = append(prevNextBntRow, prevBnt)
		}
		if page < pagesCount {
			nextBnt := tgbotapi.NewInlineKeyboardButtonData("Следующая страница", "next_page")
			prevNextBntRow = append(prevNextBntRow, nextBnt)
		}
		keyboard.InlineKeyboard = append(keyboard.InlineKeyboard, prevNextBntRow)

		i := 0
		for i < len(tasks) {
			projectsRow := make([]tgbotapi.InlineKeyboardButton, 0)
			firstRowProjectBtn := tgbotapi.NewInlineKeyboardButtonData(tasks[i].Description, strconv.Itoa(i+1))
			projectsRow = append(projectsRow, firstRowProjectBtn)
			i++

			if i != len(tasks) {
				secondRowProjectBtn := tgbotapi.NewInlineKeyboardButtonData(tasks[i].Description, strconv.Itoa(i+1))
				projectsRow = append(projectsRow, secondRowProjectBtn)
				i++
			}

			keyboard.InlineKeyboard = append(keyboard.InlineKeyboard, projectsRow)
		}
	}

	row := make([]tgbotapi.InlineKeyboardButton, 0)
	createBtn := tgbotapi.NewInlineKeyboardButtonData("Создать новую задачу", "create_task")
	rootBtn := tgbotapi.NewInlineKeyboardButtonData("Назад", "back_btn")
	row = append(row, createBtn)
	row = append(row, rootBtn)
	keyboard.InlineKeyboard = append(keyboard.InlineKeyboard, row)

	msg.ReplyMarkup = keyboard

	return msg
}
