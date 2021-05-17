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

func MakeTaskStatusesMenu(
	message *util.MessageData,
	project *model.Project,
	taskStatuses []*model.TaskStatus,
	page int,
	count int,
) tgbotapi.MessageConfig {
	msg := tgbotapi.MessageConfig{}
	textStrings := make([]string, len(taskStatuses))

	for i, status := range taskStatuses {
		textStrings[i] = fmt.Sprintf("%d. '%s' level %d", i+1, status.Status, status.Level)
	}
	text := fmt.Sprintf(
		"Статусы задач проекта '%s':\n%s\n",
		project.Name,
		strings.Join(textStrings, "\n"),
	)
	msg = tgbotapi.NewMessage(message.Chat.ID, text)

	keyboard := tgbotapi.InlineKeyboardMarkup{}

	row1 := make([]tgbotapi.InlineKeyboardButton, 0)
	row2 := make([]tgbotapi.InlineKeyboardButton, 0)
	addBtn := tgbotapi.NewInlineKeyboardButtonData("Добавить статус", "add_status")
	row1 = append(row1, addBtn)
	if len(taskStatuses) > 1 {
		removeBtn := tgbotapi.NewInlineKeyboardButtonData("Удалить статус", "remove_status")
		row1 = append(row1, removeBtn)
	}
	projectMenuBtn := tgbotapi.NewInlineKeyboardButtonData("Назад", "back_btn")
	row2 = append(row2, projectMenuBtn)
	keyboard.InlineKeyboard = append(keyboard.InlineKeyboard, row1)
	keyboard.InlineKeyboard = append(keyboard.InlineKeyboard, row2)

	if len(taskStatuses) != 0 {
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
	}

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
	if len(taskStatuses) != 0 {
		textStrings := make([]string, len(taskStatuses))
		for i, status := range taskStatuses {
			textStrings[i] = fmt.Sprintf("%d. '%s' level %d", i+1, status.Status, status.Level)
		}
		text := fmt.Sprintf(
			"Выберите статус задач для удаления:\n%s\n",
			strings.Join(textStrings, "\n"),
		)
		msg = tgbotapi.NewMessage(message.Chat.ID, text)
	} else {
		text := "Вы ещё не добавили ни одного участника"
		msg = tgbotapi.NewMessage(message.Chat.ID, text)
	}

	keyboard := tgbotapi.InlineKeyboardMarkup{}

	row1 := make([]tgbotapi.InlineKeyboardButton, 0)
	cancelBtn := tgbotapi.NewInlineKeyboardButtonData("Отмена", "cancel_btn")
	row1 = append(row1, cancelBtn)
	keyboard.InlineKeyboard = append(keyboard.InlineKeyboard, row1)

	if len(taskStatuses) != 0 {
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
		for i < len(taskStatuses) {
			membersRow := make([]tgbotapi.InlineKeyboardButton, 0)
			firstRowMemberBtn := tgbotapi.NewInlineKeyboardButtonData(taskStatuses[i].Status, strconv.Itoa(i+1))
			membersRow = append(membersRow, firstRowMemberBtn)
			i++

			if i != len(taskStatuses) {
				secondRowMemberBtn := tgbotapi.NewInlineKeyboardButtonData(taskStatuses[i].Status, strconv.Itoa(i+1))
				membersRow = append(membersRow, secondRowMemberBtn)
				i++
			}

			keyboard.InlineKeyboard = append(keyboard.InlineKeyboard, membersRow)
		}
	}

	msg.ReplyMarkup = keyboard

	return msg
}
