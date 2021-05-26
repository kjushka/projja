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

func MakeMembersMenu(
	message *util.MessageData,
	project *model.Project,
	members []*model.User,
	page int,
	count int,
) tgbotapi.MessageConfig {
	msg := tgbotapi.MessageConfig{}
	textStrings := make([]string, len(members))
	for i, member := range members {
		textStrings[i] = fmt.Sprintf("%d. %s aka %s", i+1, member.Name, member.Username)
	}
	text := fmt.Sprintf(
		"Участники проекта '%s':\n%s\n",
		project.Name,
		strings.Join(textStrings, "\n"),
	)
	msg = tgbotapi.NewMessage(message.Chat.ID, text)

	keyboard := tgbotapi.InlineKeyboardMarkup{}

	row1 := make([]tgbotapi.InlineKeyboardButton, 0)
	row2 := make([]tgbotapi.InlineKeyboardButton, 0)
	addBtn := tgbotapi.NewInlineKeyboardButtonData("Добавить участника", "add_member")
	row1 = append(row1, addBtn)
	if len(members) != 1 {
		removeBtn := tgbotapi.NewInlineKeyboardButtonData("Удалить участника", "remove_member")
		row1 = append(row1, removeBtn)
	}
	projectMenuBtn := tgbotapi.NewInlineKeyboardButtonData("Назад", "back_btn")
	row2 = append(row2, projectMenuBtn)
	keyboard.InlineKeyboard = append(keyboard.InlineKeyboard, row1)
	keyboard.InlineKeyboard = append(keyboard.InlineKeyboard, row2)

	if len(members) != 0 {
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

func MakeMembersRemovingMenu(
	message *util.MessageData,
	project *model.Project,
	members []*model.User,
	page int,
	count int,
) tgbotapi.MessageConfig {
	msg := tgbotapi.MessageConfig{}
	textStrings := make([]string, len(members))
	for i, member := range members {
		textStrings[i] = fmt.Sprintf("%d. %s aka %s", i+1, member.Name, member.Username)
	}
	text := fmt.Sprintf(
		"Выберите участника проекта '%s' для удаления:\n%s\n",
		project.Name,
		strings.Join(textStrings, "\n"),
	)
	msg = tgbotapi.NewMessage(message.Chat.ID, text)

	keyboard := tgbotapi.InlineKeyboardMarkup{}

	if len(members) != 0 {
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
		for i < len(members) {
			membersRow := make([]tgbotapi.InlineKeyboardButton, 0)
			firstRowMemberBtn := tgbotapi.NewInlineKeyboardButtonData(members[i].Name, strconv.Itoa(i+1))
			membersRow = append(membersRow, firstRowMemberBtn)
			i++

			if i != len(members) {
				secondRowMemberBtn := tgbotapi.NewInlineKeyboardButtonData(members[i].Name, strconv.Itoa(i+1))
				membersRow = append(membersRow, secondRowMemberBtn)
				i++
			}

			keyboard.InlineKeyboard = append(keyboard.InlineKeyboard, membersRow)
		}
	}

	row1 := make([]tgbotapi.InlineKeyboardButton, 0)
	cancelBtn := tgbotapi.NewInlineKeyboardButtonData("Отмена", "cancel_btn")
	row1 = append(row1, cancelBtn)
	keyboard.InlineKeyboard = append(keyboard.InlineKeyboard, row1)

	msg.ReplyMarkup = keyboard

	return msg
}
