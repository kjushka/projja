package menu

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"math"
	"projja_telegram/command/util"
	"projja_telegram/model"
	"strings"
)

func MakeMembersMenu(
	message *util.MessageData,
	project *model.Project,
	members []*model.User,
	page int,
	count int,
) tgbotapi.MessageConfig {
	msg := tgbotapi.MessageConfig{}
	start := (page - 1) * 4
	end := start + count
	textStrings := make([]string, len(members[start:end]))
	for i, member := range members[start:end] {
		textStrings[i] = fmt.Sprintf("%d. %s aka %s", i+1, member.Name, member.Username)
	}
	text := fmt.Sprintf(
		"Участники проекта '%s':\n%s\n",
		project.Name,
		strings.Join(textStrings, "\n"),
	)
	msg = tgbotapi.NewMessage(message.Chat.ID, text)

	rows := make([][]tgbotapi.KeyboardButton, 0)

	if len(members) != 0 {
		pagesCount := int(math.Ceil(float64(len(members)) / 4.0))
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

	row1 := make([]tgbotapi.KeyboardButton, 0)
	row2 := make([]tgbotapi.KeyboardButton, 0)
	addBtn := tgbotapi.NewKeyboardButton("Добавить участника")
	row1 = append(row1, addBtn)
	if len(members) != 1 {
		removeBtn := tgbotapi.NewKeyboardButton("Удалить участника")
		row1 = append(row1, removeBtn)
	}
	projectMenuBtn := tgbotapi.NewKeyboardButton("Назад")
	row2 = append(row2, projectMenuBtn)
	rows = append(rows, row1, row2)

	keyboard := tgbotapi.NewReplyKeyboard(rows...)
	msg.ReplyMarkup = keyboard

	return msg
}

func MakeMembersRemovingMenu(
	message *util.MessageData,
	project *model.Project,
	members []*model.User,
	page int,
	count int,
) tgbotapi.MessageConfig {
	msg := tgbotapi.MessageConfig{}
	start := (page - 1) * 4
	end := start + count
	textStrings := make([]string, len(members[start:end]))
	for i, member := range members[start:end] {
		textStrings[i] = fmt.Sprintf("%d. %s aka %s", i+1, member.Name, member.Username)
	}
	text := fmt.Sprintf(
		"Выберите участника проекта '%s' для удаления:\n%s\n",
		project.Name,
		strings.Join(textStrings, "\n"),
	)
	msg = tgbotapi.NewMessage(message.Chat.ID, text)

	rows := make([][]tgbotapi.KeyboardButton, 0)

	i := start
	for i < end {
		membersRow := make([]tgbotapi.KeyboardButton, 0)
		firstRowMemberBtn := tgbotapi.NewKeyboardButton(members[i].Username)
		membersRow = append(membersRow, firstRowMemberBtn)
		i++

		if i != end {
			secondRowMemberBtn := tgbotapi.NewKeyboardButton(members[i].Username)
			membersRow = append(membersRow, secondRowMemberBtn)
			i++
		}

		rows = append(rows, membersRow)
	}

	pagesCount := int(math.Ceil(float64(len(members)) / 4.0))
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
