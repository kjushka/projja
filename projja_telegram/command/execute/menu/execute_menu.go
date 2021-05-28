package menu

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"math"
	"projja_telegram/command/util"
	"projja_telegram/model"
	"strings"
)

func MakeExecutedTasksMenu(message *util.MessageData, tasks []*model.Task, page int, count int) tgbotapi.MessageConfig {
	msg := tgbotapi.MessageConfig{}
	start := (page - 1) * 4
	end := start + count
	if len(tasks) != 0 {
		textStrings := make([]string, len(tasks[start:end]))
		for i, task := range tasks[start:end] {
			var priority string
			switch task.Priority {
			case "critical":
				priority = "критический"
			case "high":
				priority = "высокий"
			case "medium":
				priority = "средний"
			case "low":
				priority = "низкий"
			}
			textStrings[i] = fmt.Sprintf("%d. %s до %s, приоритет: %s", i+1, task.Description, task.Deadline, priority)
		}
		text := fmt.Sprintf(
			"Ваши задачи:\n%s\n",
			strings.Join(textStrings, "\n"),
		)
		msg = tgbotapi.NewMessage(message.Chat.ID, text)
	} else {
		text := "Вы ещё не получили ни одной задачи"
		msg = tgbotapi.NewMessage(message.Chat.ID, text)
	}

	rows := make([][]tgbotapi.KeyboardButton, 0)
	if len(tasks) != 0 {
		i := start
		for i < end {
			tasksRow := make([]tgbotapi.KeyboardButton, 0)
			firstRowTaskBtn := tgbotapi.NewKeyboardButton(tasks[i].Description)
			tasksRow = append(tasksRow, firstRowTaskBtn)
			i++

			if i != end {
				secondRowTaskBtn := tgbotapi.NewKeyboardButton(tasks[i].Description)
				tasksRow = append(tasksRow, secondRowTaskBtn)
				i++
			}

			rows = append(rows, tasksRow)
		}

		pagesCount := int(math.Ceil(float64(len(tasks)) / 4.0))
		prevNextBntRow := make([]tgbotapi.KeyboardButton, 0)
		if page > 1 {
			prevBnt := tgbotapi.NewKeyboardButton("Предыдущая страница")
			prevNextBntRow = append(prevNextBntRow, prevBnt)
		}
		if page < pagesCount {
			nextBnt := tgbotapi.NewKeyboardButton("Следующая страница")
			prevNextBntRow = append(prevNextBntRow, nextBnt)
		}
		rows = append(rows, prevNextBntRow)
	}

	row := make([]tgbotapi.KeyboardButton, 0)
	rootBtn := tgbotapi.NewKeyboardButton("Назад")
	row = append(row, rootBtn)
	rows = append(rows, row)

	keyboard := tgbotapi.NewReplyKeyboard(rows...)
	msg.ReplyMarkup = keyboard

	return msg
}
