package menu

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"log"
	"math"
	"projja_telegram/command/util"
	"projja_telegram/model"
	"strings"
)

func MakeProjectAnswersMenu(message *util.MessageData, answers []*model.Answer, page int, count int) tgbotapi.MessageConfig {
	msg := tgbotapi.MessageConfig{}
	start := (page - 1) * 4
	end := start + count
	if len(answers) != 0 {
		textStrings := make([]string, len(answers[start:end]))

		log.Println(start, count, end)
		for i, answer := range answers[start:end] {
			textStrings[i] = fmt.Sprintf(
				"%d. Ответ на задачу '%s', отправлен %s",
				i+1,
				answer.Task.Description,
				answer.SentAt.Format("2006-01-02"),
			)
		}
		text := fmt.Sprintf("Новые ответы на задачи:\n%s\n",
			strings.Join(textStrings, "\n"),
		)
		msg = tgbotapi.NewMessage(message.Chat.ID, text)
	} else {
		text := "Вы ещё не получили новых ответов"
		msg = tgbotapi.NewMessage(message.Chat.ID, text)
	}

	rows := make([][]tgbotapi.KeyboardButton, 0)

	if len(answers) != 0 {
		i := start
		for i < end {
			tasksRow := make([]tgbotapi.KeyboardButton, 0)
			firstRowTaskBtn := tgbotapi.NewKeyboardButton(answers[i].Task.Description)
			tasksRow = append(tasksRow, firstRowTaskBtn)
			i++

			if i != end {
				secondRowTaskBtn := tgbotapi.NewKeyboardButton(answers[i].Task.Description)
				tasksRow = append(tasksRow, secondRowTaskBtn)
				i++
			}

			rows = append(rows, tasksRow)
		}

		pagesCount := int(math.Ceil(float64(len(answers)) / 4.0))
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
