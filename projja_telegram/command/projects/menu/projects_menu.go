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
			"Для работы с проектом введите его номер в списке или создайте новый проект", strings.Join(textStrings, "\n"))
		msg = tgbotapi.NewMessage(message.Chat.ID, text)
	} else {
		text := "Вы ещё не создали ни одного проекта"
		msg = tgbotapi.NewMessage(message.Chat.ID, text)
	}

	keyboard := tgbotapi.InlineKeyboardMarkup{}

	row := make([]tgbotapi.InlineKeyboardButton, 0)
	createBtn := tgbotapi.NewInlineKeyboardButtonData("Создать новый проект", "create_project")
	rootBtn := tgbotapi.NewInlineKeyboardButtonData("Главное меню", "root")
	row = append(row, createBtn)
	row = append(row, rootBtn)
	keyboard.InlineKeyboard = append(keyboard.InlineKeyboard, row)

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

		upperProjectRow := make([]tgbotapi.InlineKeyboardButton, 0)
		upperRowCount := len(projects)
		if len(projects) > 5 {
			upperRowCount = 5
		}
		for i := 0; i < upperRowCount; i++ {
			projectIndexBtn := tgbotapi.NewInlineKeyboardButtonData(strconv.Itoa(i+1), strconv.Itoa(i+1))
			upperProjectRow = append(upperProjectRow, projectIndexBtn)
		}
		keyboard.InlineKeyboard = append(keyboard.InlineKeyboard, upperProjectRow)

		if len(projects) > 5 {
			lowerProjectRow := make([]tgbotapi.InlineKeyboardButton, 0)
			lowerRowCount := len(projects)
			for i := 5; i < lowerRowCount; i++ {
				projectIndexBtn := tgbotapi.NewInlineKeyboardButtonData(strconv.Itoa(i+1), strconv.Itoa(i+1))
				lowerProjectRow = append(lowerProjectRow, projectIndexBtn)
			}
			keyboard.InlineKeyboard = append(keyboard.InlineKeyboard, lowerProjectRow)
		}
	}

	msg.ReplyMarkup = keyboard

	return msg, projects, true
}