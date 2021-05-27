package view

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"log"
	"math"
	"projja_telegram/command/current_project/controller"
	"projja_telegram/command/current_project/menu"
	"projja_telegram/command/util"
	"projja_telegram/model"
	"strings"
	"time"
)

func ManageProjectTasks(botUtil *util.BotUtil, project *model.Project) {
	page := 1
	tasks, msg, status := ShowTasksMenu(botUtil, project, page)
	botUtil.Bot.Send(msg)
	if !status {
		return
	}

	for update := range botUtil.Updates {
		mes := update.Message
		command := ""

		if mes.Text != "" {
			command = mes.Text
		}

		switch command {
		case "Создать новую задачу":
			msg = CreateTask(botUtil, project)
			botUtil.Bot.Send(msg)
		case "Предыдущая страница":
			page--
		case "Следующая страница":
			page++
		case "Назад":
			return
		default:
			text, index, status := IsTaskName(tasks, command)
			msg = tgbotapi.NewMessage(botUtil.Message.Chat.ID, text)
			botUtil.Bot.Send(msg)
			if status {
				ManageTask(botUtil, project, tasks[index])
			}
		}

		tasks, msg, status = ShowTasksMenu(botUtil, project, page)
		botUtil.Bot.Send(msg)
		if !status {
			return
		}
	}
}

func ShowTasksMenu(botUtil *util.BotUtil, project *model.Project, page int) ([]*model.Task, tgbotapi.MessageConfig, bool) {
	tasks, status := controller.GetProjectTasks(project)
	if !status {
		errorText := "Не удалось получить список задач\n" +
			"Попробуйте позже"
		msg := tgbotapi.NewMessage(botUtil.Message.Chat.ID, errorText)
		return nil, msg, false
	}

	count := len(tasks) - (page-1)*4
	if count > 4 {
		count = 4
	}
	msg := menu.MakeProjectTasksMenu(botUtil.Message, project, tasks, page, count)
	return tasks, msg, true
}

func CreateTask(botUtil *util.BotUtil, project *model.Project) tgbotapi.MessageConfig {
	cancelText := "Отмена создания задачи"
	text, cancelStatus := util.ListenForText(botUtil,
		"Введите название описание новой задачи",
		cancelText,
	)
	if !cancelStatus {
		msg := tgbotapi.NewMessage(botUtil.Message.Chat.ID, text)
		return msg
	}

	newTaskDescription := text

	newTaskDeadline, cancelStatus := listenForDeadline(botUtil)

	if !cancelStatus {
		msg := tgbotapi.NewMessage(botUtil.Message.Chat.ID, cancelText)
		return msg
	}

	newTaskPriority, cancelStatus := listenForPriority(botUtil)

	if !cancelStatus {
		msg := tgbotapi.NewMessage(botUtil.Message.Chat.ID, cancelText)
		return msg
	}

	newTaskSkills, cancelStatus := listenForTaskSkills(botUtil)

	if !cancelStatus {
		msg := tgbotapi.NewMessage(botUtil.Message.Chat.ID, cancelText)
		return msg
	}

	task := &model.Task{
		Description: newTaskDescription,
		Deadline:    newTaskDeadline.Format("2006-01-02"),
		Priority:    newTaskPriority,
		Skills:      newTaskSkills,
	}

	acceptingString := fmt.Sprintf("Вы хотите создать:\n"+
		"Описание: %s\n"+
		"Deadline: %s\n"+
		"Приоритет: %s\n"+
		"Навыки: %s\n",
		task.Description,
		task.Deadline,
		task.Priority,
		strings.Join(task.Skills, ", "),
	)
	msg := util.GetAcceptingMessage(botUtil.Message, acceptingString)

	botUtil.Bot.Send(msg)

	var executor *model.User

	for update := range botUtil.Updates {
		mes := update.Message
		command := ""

		if mes.Text != "" {
			command = mes.Text
		}

		switch command {
		case "Да":
			text = "Вычисляем подходящего исполнителя..."
			msg = tgbotapi.NewMessage(botUtil.Message.Chat.ID, text)
			botUtil.Bot.Send(msg)

			user, err := controller.CalculateExecutor(project, task)
			if err != nil {
				errorText := "При получении исполнителя произошла ошибка\nПопробуйте позже ещё раз"
				msg = tgbotapi.NewMessage(botUtil.Message.Chat.ID, errorText)
				return msg
			}
			executor = user

			goto LOOP
		case "Нет":
			text = cancelText
			msg = tgbotapi.NewMessage(botUtil.Message.Chat.ID, cancelText)
			return msg
		default:
			text = "Пожалуйста, выберите один из вариантов"
			msg = tgbotapi.NewMessage(botUtil.Message.Chat.ID, text)
			botUtil.Bot.Send(msg)

			msg = util.GetAcceptingMessage(botUtil.Message, acceptingString)
			botUtil.Bot.Send(msg)
		}
	}

LOOP:
	text = fmt.Sprintf("Предлагаемый исполнитель:\n"+
		"Имя: %s\nUsername: %s\nНавыки: %s\n",
		executor.Name,
		executor.Username,
		strings.Join(executor.Skills, ", "),
	)
	msg = tgbotapi.NewMessage(botUtil.Message.Chat.ID, text)
	botUtil.Bot.Send(msg)
	acceptingString = "Вы согласны с выбором?"
	msg = util.GetAcceptingMessage(botUtil.Message, acceptingString)
	botUtil.Bot.Send(msg)

	for update := range botUtil.Updates {
		mes := update.Message
		command := ""

		if mes.Text != "" {
			command = mes.Text
		}

		switch command {
		case "Да":
			task.Executor = executor
			createText, status := controller.CreateTask(project, task)
			text = createText

			if status {
				log.Println(botUtil.Message.Chat.ID, executor.Id)
				msg := tgbotapi.NewMessage(executor.ChatId, fmt.Sprintf("Вы получили новую задачу:\n%s", task.Description))
				botUtil.Bot.Send(msg)
			}

			goto BREAK
		case "Нет":
			executor, text = listenForExecutor(botUtil, project, cancelText)
			if executor == nil {
				return tgbotapi.NewMessage(botUtil.Message.Chat.ID, text)
			}
			task.Executor = executor
			text, _ = controller.CreateTask(project, task)
			goto BREAK
		default:
			text = "Пожалуйста, выберите один из вариантов"
			msg = tgbotapi.NewMessage(botUtil.Message.Chat.ID, text)
			botUtil.Bot.Send(msg)

			msg = util.GetAcceptingMessage(botUtil.Message, acceptingString)
			botUtil.Bot.Send(msg)
		}
	}

BREAK:
	msg = tgbotapi.NewMessage(botUtil.Message.Chat.ID, text)
	return msg
}

func listenForDeadline(botUtil *util.BotUtil) (time.Time, bool) {
	mesText := "Введите дату дедлайна в формате YYYY-MM-DD"
	msg := tgbotapi.NewMessage(botUtil.Message.Chat.ID, mesText)

	row := make([]tgbotapi.KeyboardButton, 0)

	cancelBtn := tgbotapi.NewKeyboardButton("Отмена")
	row = append(row, cancelBtn)

	keyboard := tgbotapi.NewReplyKeyboard(row)
	msg.ReplyMarkup = keyboard

	botUtil.Bot.Send(msg)

	result := time.Time{}
	ready := false
	for update := range botUtil.Updates {
		mes := update.Message
		command := ""

		if mes.Text != "" {
			command = mes.Text
		}

		switch command {
		case "Отмена":
			return time.Now(), false
		default:
			if command == "" {
				botUtil.Bot.Send(msg)
				continue
			}

			resultDeadline, err := time.Parse("2006-01-02", command)
			if err != nil {
				log.Println("incorrect time format")
				text := "Вы ввели дату в неверном формате\nПопробуйте ещё раз"
				errorMsg := tgbotapi.NewMessage(botUtil.Message.Chat.ID, text)
				botUtil.Bot.Send(errorMsg)
				botUtil.Bot.Send(msg)
				continue
			}
			if time.Until(resultDeadline) < 0 {
				text := "Вы ввели дату, которая уже прошла\nПопробуйте ещё раз"
				errorMsg := tgbotapi.NewMessage(botUtil.Message.Chat.ID, text)
				botUtil.Bot.Send(errorMsg)
				botUtil.Bot.Send(msg)
				continue
			}

			result = resultDeadline
			ready = true
		}

		if ready {
			break
		}

		botUtil.Bot.Send(msg)
	}

	return result, true
}

func listenForTaskSkills(botUtil *util.BotUtil) ([]string, bool) {
	mesText := "Перечислите через пробел навыки, которые нужны для выполнения задачи\n" +
		"Пример:\nfrontend js angular"
	msg := tgbotapi.NewMessage(botUtil.Message.Chat.ID, mesText)

	row := make([]tgbotapi.KeyboardButton, 0)

	cancelBtn := tgbotapi.NewKeyboardButton("Отмена")
	row = append(row, cancelBtn)

	keyboard := tgbotapi.NewReplyKeyboard(row)
	msg.ReplyMarkup = keyboard

	botUtil.Bot.Send(msg)

	var result []string
	ready := false
	for update := range botUtil.Updates {
		mes := update.Message
		command := ""

		if mes.Text != "" {
			command = mes.Text
		}

		switch command {
		case "Отмена":
			return nil, false
		default:
			if command == "" {
				botUtil.Bot.Send(msg)
				continue
			}

			result = strings.Split(command, " ")
			ready = true
		}

		if ready {
			break
		}
	}

	return result, true
}

func listenForPriority(botUtil *util.BotUtil) (string, bool) {
	mesText := "Выберите приоритет задачи"
	msg := tgbotapi.NewMessage(botUtil.Message.Chat.ID, mesText)

	priorityRow := make([]tgbotapi.KeyboardButton, 0)
	critical := tgbotapi.NewKeyboardButton("Критический")
	high := tgbotapi.NewKeyboardButton("Высокий")
	medium := tgbotapi.NewKeyboardButton("Средний")
	low := tgbotapi.NewKeyboardButton("Низкий")
	priorityRow = append(priorityRow, critical)
	priorityRow = append(priorityRow, high)
	priorityRow = append(priorityRow, medium)
	priorityRow = append(priorityRow, low)

	row := make([]tgbotapi.KeyboardButton, 0)
	cancelBtn := tgbotapi.NewKeyboardButton("Отмена")
	row = append(row, cancelBtn)

	keyboard := tgbotapi.NewReplyKeyboard(priorityRow, row)

	msg.ReplyMarkup = keyboard

	botUtil.Bot.Send(msg)

	var result string
	ready := false
	for update := range botUtil.Updates {
		mes := update.Message
		command := ""

		if mes.Text != "" {
			command = mes.Text
		}

		switch command {
		case "Отмена":
			return "", false
		case "Критический":
			result = "critical"
			ready = true
		case "Высокий":
			result = "high"
			ready = true
		case "Средний":
			result = "medium"
			ready = true
		case "Низкий":
			result = "low"
			ready = true
		default:
			text := "Пожалуйста, выберите один из вариантов"
			errorMsg := tgbotapi.NewMessage(botUtil.Message.Chat.ID, text)
			botUtil.Bot.Send(errorMsg)
			botUtil.Bot.Send(msg)
		}

		if ready && result != "" {
			break
		}
	}

	var priority string
	switch result {
	case "critical":
		priority = "критический"
	case "high":
		priority = "высокий"
	case "medium":
		priority = "средний"
	case "low":
		priority = "низкий"
	}
	text := fmt.Sprintf("Для задачи был выбран %s приоритет", priority)
	msg = tgbotapi.NewMessage(botUtil.Message.Chat.ID, text)
	botUtil.Bot.Send(msg)

	return result, true
}

func listenForExecutor(botUtil *util.BotUtil, project *model.Project, cancelString string) (*model.User, string) {
	members, status := controller.GetMembers(project)
	if !status {
		errorText := "Не удалось получить список участников\n" +
			"Попробуйте позже"
		return nil, errorText
	}

	page := 1
	count := len(members) - (page-1)*4
	if count > 4 {
		count = 4
	}
	msg := makeExecutorMenu(botUtil.Message, members, page, count)
	botUtil.Bot.Send(msg)

	var memberIndex int
	exit := false
	for update := range botUtil.Updates {
		mes := update.Message
		command := ""

		if mes.Text != "" {
			command = mes.Text
		}

		switch command {
		case "Отмена":
			return nil, cancelString
		case "Предыдущая страница":
			page--
		case "Следующая страница":
			page++
		default:
			text, index, status := IsMemberName(members, command)
			memberIndex = index
			if !status {
				msg := tgbotapi.NewMessage(botUtil.Message.Chat.ID, text)
				botUtil.Bot.Send(msg)
			} else {
				exit = true
			}
		}

		if exit && memberIndex != -1 {
			break
		}

		msg = makeExecutorMenu(botUtil.Message, members, page, count)
		botUtil.Bot.Send(msg)
	}

	return members[memberIndex], ""
}

func makeExecutorMenu(message *util.MessageData, members []*model.User, page int, count int) tgbotapi.MessageConfig {
	msg := tgbotapi.MessageConfig{}
	start := (page - 1) * 4
	end := start + count
	textStrings := make([]string, len(members[start:end]))
	for i, member := range members {
		textStrings[i] = fmt.Sprintf("%d. %s aka %s", i+1, member.Name, member.Username)
	}
	text := fmt.Sprintf(
		"Выберите исполнителя:\n%s\n",
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

	pagesCount := int(math.Ceil(float64(len(members)) / 10.0))
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

	row := make([]tgbotapi.KeyboardButton, 0)
	projectMenuBtn := tgbotapi.NewKeyboardButton("Отмена")
	row = append(row, projectMenuBtn)
	rows = append(rows, row)

	keyboard := tgbotapi.NewReplyKeyboard(rows...)
	msg.ReplyMarkup = keyboard

	return msg
}

func IsTaskName(tasks []*model.Task, command string) (string, int, bool) {
	if command == "" {
		text := "Задачи с таким описанием не существует"
		return text, -1, false
	}

	index := -1
	found := false
	for i, task := range tasks {
		if task.Description == command {
			found = true
			index = i
			break
		}
	}

	if !found {
		text := "Задачи с таким описанием не существует"
		return text, index, found
	}

	return "", index, found
}
