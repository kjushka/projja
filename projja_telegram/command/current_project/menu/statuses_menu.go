package menu

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"math"
	"projja_telegram/command/util"
	"projja_telegram/model"
	"strings"
)

func MakeTaskStatusesMenu(
	message *util.MessageData,
	project *model.Project,
	taskStatuses []*model.TaskStatus,
	page int,
	count int,
) tgbotapi.MessageConfig {
	msg := tgbotapi.MessageConfig{}
	start := (page - 1) * 4
	end := start + count
	textStrings := make([]string, len(taskStatuses[start:end]))
	for i, status := range taskStatuses[start:end] {
		textStrings[i] = fmt.Sprintf("%d. '%s' level %d", i+1, status.Status, status.Level)
	}
	text := fmt.Sprintf(
		"Статусы задач проекта '%s':\n%s\n",
		project.Name,
		strings.Join(textStrings, "\n"),
	)
	msg = tgbotapi.NewMessage(message.Chat.ID, text)

	rows := make([][]tgbotapi.KeyboardButton, 0)

	row1 := make([]tgbotapi.KeyboardButton, 0)
	row2 := make([]tgbotapi.KeyboardButton, 0)
	addBtn := tgbotapi.NewKeyboardButton("Добавить статус")
	row1 = append(row1, addBtn)
	if len(taskStatuses) > 1 {
		removeBtn := tgbotapi.NewKeyboardButton("Удалить статус")
		row1 = append(row1, removeBtn)
	}
	projectMenuBtn := tgbotapi.NewKeyboardButton("Назад")
	row2 = append(row2, projectMenuBtn)
	rows = append(rows, row1, row2)

	pagesCount := int(math.Ceil(float64(len(taskStatuses)) / 4.0))
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

	keyboard := tgbotapi.NewReplyKeyboard(rows...)
	msg.ReplyMarkup = keyboard

	return msg
}

func MakeTaskStatusesRemovingMenu(
	message *util.MessageData,
	taskStatuses []*model.TaskStatus,
	page int,
	count int,
) tgbotapi.MessageConfig {
	msg := tgbotapi.MessageConfig{}
	start := (page - 1) * 4
	end := start + count
	textStrings := make([]string, len(taskStatuses[start:end]))
	for i, status := range taskStatuses[start:end] {
		textStrings[i] = fmt.Sprintf("%d. '%s' level %d", i+1, status.Status, status.Level)
	}
	text := fmt.Sprintf(
		"Выберите статус задач для удаления:\n%s\n",
		strings.Join(textStrings, "\n"),
	)
	msg = tgbotapi.NewMessage(message.Chat.ID, text)

	rows := make([][]tgbotapi.KeyboardButton, 0)

	i := start
	for i < end {
		taskStatusesRow := make([]tgbotapi.KeyboardButton, 0)
		firstRowMemberBtn := tgbotapi.NewKeyboardButton(taskStatuses[i].Status)
		taskStatusesRow = append(taskStatusesRow, firstRowMemberBtn)
		i++

		if i != end {
			secondRowMemberBtn := tgbotapi.NewKeyboardButton(taskStatuses[i].Status)
			taskStatusesRow = append(taskStatusesRow, secondRowMemberBtn)
			i++
		}

		rows = append(rows, taskStatusesRow)
	}

	pagesCount := int(math.Ceil(float64(len(taskStatuses)) / 4.0))
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

	row1 := make([]tgbotapi.KeyboardButton, 0)
	cancelBtn := tgbotapi.NewKeyboardButton("Отмена")
	row1 = append(row1, cancelBtn)
	rows = append(rows, row1)

	keyboard := tgbotapi.NewReplyKeyboard(rows...)
	msg.ReplyMarkup = keyboard

	return msg
}
