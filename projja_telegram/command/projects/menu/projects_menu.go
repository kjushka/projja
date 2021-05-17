package menu

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"math"
	"projja_telegram/command/projects/controller"
	"projja_telegram/command/util"
	"projja_telegram/model"
	"strconv"
	"strings"
)

func MakeProjectsMenu(message *util.MessageData, page int, count int) (tgbotapi.MessageConfig, []*model.Project, bool) {
	projects, status := controller.GetProjects(message.From, page, count)

	if !status {
		errorText := "Не удалось получить список проектов\n" +
			"Попробуйте позже"
		msg := tgbotapi.NewMessage(message.Chat.ID, errorText)
		return msg, nil, false
	}

	msg := tgbotapi.MessageConfig{}
	if len(projects) != 0 {
		textStrings := make([]string, len(projects))
		for i, project := range projects {
			textStrings[i] = fmt.Sprintf("%d. '%s' статус: %s", i+1, project.Name, project.Status)
		}
		text := fmt.Sprintf("Ваши проекты:\n%s\n"+
			"Выберите проект для работы", strings.Join(textStrings, "\n"))
		msg = tgbotapi.NewMessage(message.Chat.ID, text)
	} else {
		text := "Вы ещё не создали ни одного проекта"
		msg = tgbotapi.NewMessage(message.Chat.ID, text)
	}

	keyboard := tgbotapi.InlineKeyboardMarkup{}

	if len(projects) != 0 {
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
		for i < len(projects) {
			projectsRow := make([]tgbotapi.InlineKeyboardButton, 0)
			firstRowProjectBtn := tgbotapi.NewInlineKeyboardButtonData(projects[i].Name, strconv.Itoa(i+1))
			projectsRow = append(projectsRow, firstRowProjectBtn)
			i++

			if i != len(projects) {
				secondRowProjectBtn := tgbotapi.NewInlineKeyboardButtonData(projects[i].Name, strconv.Itoa(i+1))
				projectsRow = append(projectsRow, secondRowProjectBtn)
				i++
			}

			keyboard.InlineKeyboard = append(keyboard.InlineKeyboard, projectsRow)
		}
	}

	row := make([]tgbotapi.InlineKeyboardButton, 0)
	createBtn := tgbotapi.NewInlineKeyboardButtonData("Создать новый проект", "create_project")
	rootBtn := tgbotapi.NewInlineKeyboardButtonData("Назад", "back_btn")
	row = append(row, createBtn)
	row = append(row, rootBtn)
	keyboard.InlineKeyboard = append(keyboard.InlineKeyboard, row)

	msg.ReplyMarkup = keyboard

	return msg, projects, true
}
