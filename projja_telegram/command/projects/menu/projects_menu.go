package menu

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"math"
	"projja_telegram/command/projects/controller"
	"projja_telegram/command/util"
	"projja_telegram/model"
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

	rows := make([][]tgbotapi.KeyboardButton, 0)

	if len(projects) != 0 {
		i := 0
		for i < len(projects) {
			projectsRow := make([]tgbotapi.KeyboardButton, 0)
			firstRowProjectBtn := tgbotapi.NewKeyboardButton(projects[i].Name)
			projectsRow = append(projectsRow, firstRowProjectBtn)
			i++

			if i != len(projects) {
				secondRowProjectBtn := tgbotapi.NewKeyboardButton(projects[i].Name)
				projectsRow = append(projectsRow, secondRowProjectBtn)
				i++
			}

			rows = append(rows, projectsRow)
		}

		pagesCount := int(math.Ceil(float64(count) / 4.0))
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
	createBtn := tgbotapi.NewKeyboardButton("Создать новый проект")
	rootBtn := tgbotapi.NewKeyboardButton("Назад")
	row = append(row, createBtn)
	row = append(row, rootBtn)
	rows = append(rows, row)

	keyboard := tgbotapi.NewReplyKeyboard(rows...)
	msg.ReplyMarkup = keyboard

	return msg, projects, true
}
