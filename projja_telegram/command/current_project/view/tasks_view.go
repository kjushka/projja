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
	"strconv"
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
		var command string

		if update.CallbackQuery != nil {
			response := strings.Split(update.CallbackQuery.Data, " ")
			command = response[0]
		} else if mes.IsCommand() {
			command = mes.Command()
		} else if mes.Text != "" {
			command = mes.Text
		}

		switch command {
		case "create_task":
			msg = CreateTask(botUtil, project)
			botUtil.Bot.Send(msg)
		case "prev_page":
			page--
		case "next_page":
			page++
		case "back_btn":
			return
		default:
			text, index, status := IsTaskId(command, len(tasks), page)
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

	count := len(tasks) - (page-1)*10
	if count > 10 {
		count = 10
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
		var command string

		if update.CallbackQuery != nil {
			response := strings.Split(update.CallbackQuery.Data, " ")
			command = response[0]
		} else if mes.IsCommand() {
			command = mes.Command()
		} else if mes.Text != "" {
			command = mes.Text
		}

		switch command {
		case "yes_btn":
			text = "Вычисляем подходящего исполнителя..."
			msg = tgbotapi.NewMessage(botUtil.Message.Chat.ID, text)
			botUtil.Bot.Send(msg)

			user, err := controller.CalculateExecutor(project, task)
			if err != nil {
				errorText := "При получении исполнителя произошла ошибка\nПопробуйте позже"
				msg = tgbotapi.NewMessage(botUtil.Message.Chat.ID, errorText)
				return msg
			}
			executor = user

			goto LOOP
		case "no_btn":
			text = cancelText
			msg = tgbotapi.NewMessage(botUtil.Message.Chat.ID, cancelText)
			return msg
		default:
			text = "Неизвестная команда"
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
		var command string

		if update.CallbackQuery != nil {
			response := strings.Split(update.CallbackQuery.Data, " ")
			command = response[0]
		} else if mes.IsCommand() {
			command = mes.Command()
		} else if mes.Text != "" {
			command = mes.Text
		}

		switch command {
		case "yes_btn":
			task.Executor = executor
			text, _ = controller.CreateTask(project, task)
			goto BREAK
		case "no_btn":
			executor, text = listenForExecutor(botUtil, project, cancelText)
			if executor == nil {
				return tgbotapi.NewMessage(botUtil.Message.Chat.ID, text)
			}
			task.Executor = executor
			text, _ = controller.CreateTask(project, task)
			goto BREAK
		default:
			text = "Неизвестная команда"
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

	keyboard := tgbotapi.InlineKeyboardMarkup{}
	row := make([]tgbotapi.InlineKeyboardButton, 0)

	cancelBtn := tgbotapi.NewInlineKeyboardButtonData("Отмена", "cancel_btn")
	row = append(row, cancelBtn)
	keyboard.InlineKeyboard = append(keyboard.InlineKeyboard, row)

	msg.ReplyMarkup = keyboard

	botUtil.Bot.Send(msg)

	result := time.Time{}
	ready := false
	for update := range botUtil.Updates {
		mes := update.Message
		var command string

		if update.CallbackQuery != nil {
			response := strings.Split(update.CallbackQuery.Data, " ")
			command = response[0]
		} else if mes.IsCommand() {
			command = mes.Command()
		} else if mes.Text != "" {
			command = mes.Text
		}

		switch command {
		case "cancel_btn":
			return time.Now(), false
		default:
			if command == "" {
				continue
			}

			resultDeadline, err := time.Parse("2006-01-02", command)
			if err != nil {
				log.Println("incorrect time format")
				text := "Вы ввели дату в неверном формате\nПопробуйте ещё раз"
				msg = tgbotapi.NewMessage(botUtil.Message.Chat.ID, text)
				botUtil.Bot.Send(msg)
				continue
			}
			if time.Until(resultDeadline) < 0 {
				text := "Вы ввели дату, которая уже прошла\nПопробуйте ещё раз"
				msg = tgbotapi.NewMessage(botUtil.Message.Chat.ID, text)
				botUtil.Bot.Send(msg)
				continue
			}

			result = resultDeadline
			ready = true
		}

		if ready {
			break
		}
	}

	return result, true
}

func listenForTaskSkills(botUtil *util.BotUtil) ([]string, bool) {
	mesText := "Перечислите через пробел навыки, которые нужны для выполнения задачи\n" +
		"Пример:\nfrontend js angular"
	msg := tgbotapi.NewMessage(botUtil.Message.Chat.ID, mesText)

	keyboard := tgbotapi.InlineKeyboardMarkup{}
	row := make([]tgbotapi.InlineKeyboardButton, 0)

	cancelBtn := tgbotapi.NewInlineKeyboardButtonData("Отмена", "cancel_btn")
	row = append(row, cancelBtn)
	keyboard.InlineKeyboard = append(keyboard.InlineKeyboard, row)

	msg.ReplyMarkup = keyboard

	botUtil.Bot.Send(msg)

	var result []string
	ready := false
	for update := range botUtil.Updates {
		mes := update.Message
		var command string

		if update.CallbackQuery != nil {
			response := strings.Split(update.CallbackQuery.Data, " ")
			command = response[0]
		} else if mes.IsCommand() {
			command = mes.Command()
		} else if mes.Text != "" {
			command = mes.Text
		}

		switch command {
		case "cancel_btn":
			return nil, false
		default:
			if command == "" {
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

	keyboard := tgbotapi.InlineKeyboardMarkup{}

	priorityRow := make([]tgbotapi.InlineKeyboardButton, 0)
	critical := tgbotapi.NewInlineKeyboardButtonData("Критический", "critical")
	high := tgbotapi.NewInlineKeyboardButtonData("Высокий", "high")
	medium := tgbotapi.NewInlineKeyboardButtonData("Средний", "medium")
	low := tgbotapi.NewInlineKeyboardButtonData("Низкий", "low")
	priorityRow = append(priorityRow, critical)
	priorityRow = append(priorityRow, high)
	priorityRow = append(priorityRow, medium)
	priorityRow = append(priorityRow, low)

	row := make([]tgbotapi.InlineKeyboardButton, 0)
	cancelBtn := tgbotapi.NewInlineKeyboardButtonData("Отмена", "cancel_btn")
	row = append(row, cancelBtn)

	keyboard.InlineKeyboard = append(keyboard.InlineKeyboard, priorityRow)
	keyboard.InlineKeyboard = append(keyboard.InlineKeyboard, row)

	msg.ReplyMarkup = keyboard

	botUtil.Bot.Send(msg)

	var result string
	ready := false
	for update := range botUtil.Updates {
		mes := update.Message
		var command string

		if update.CallbackQuery != nil {
			response := strings.Split(update.CallbackQuery.Data, " ")
			command = response[0]
		} else if mes.IsCommand() {
			command = mes.Command()
		} else if mes.Text != "" {
			command = mes.Text
		}

		switch command {
		case "cancel_btn":
			return "", false
		case "critical":
			fallthrough
		case "high":
			fallthrough
		case "medium":
			fallthrough
		case "low":
			result = command
			ready = true
		default:
			text := "Неизвестная команда\nВыберите статус задачи"
			msg = tgbotapi.NewMessage(botUtil.Message.Chat.ID, text)
			botUtil.Bot.Send(msg)
		}

		if ready {
			break
		}
	}

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
	count := len(members) - (page-1)*10
	if count > 10 {
		count = 10
	}
	msg := makeExecutorMenu(botUtil.Message, members, page, count)
	botUtil.Bot.Send(msg)

	var memberIndex int
	exit := false
	for update := range botUtil.Updates {
		mes := update.Message
		var command string

		if update.CallbackQuery != nil {
			response := strings.Split(update.CallbackQuery.Data, " ")
			command = response[0]
		} else if mes.IsCommand() {
			command = mes.Command()
		} else if mes.Text != "" {
			command = mes.Text
		}

		switch command {
		case "cancel_btn":
			return nil, cancelString
		case "prev_page":
			page--
		case "next_page":
			page++
		default:
			text, index, status := IsMemberId(command, len(members), page)
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
	textStrings := make([]string, len(members))
	for i, member := range members {
		textStrings[i] = fmt.Sprintf("%d. %s aka %s", i+1, member.Name, member.Username)
	}
	text := fmt.Sprintf(
		"Выберите исполнителя:\n%s\n",
		strings.Join(textStrings, "\n"),
	)
	msg = tgbotapi.NewMessage(message.Chat.ID, text)

	keyboard := tgbotapi.InlineKeyboardMarkup{}

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

	row := make([]tgbotapi.InlineKeyboardButton, 0)
	projectMenuBtn := tgbotapi.NewInlineKeyboardButtonData("Отмена", "cancel_btn")
	row = append(row, projectMenuBtn)
	keyboard.InlineKeyboard = append(keyboard.InlineKeyboard, row)

	msg.ReplyMarkup = keyboard

	return msg
}

func IsTaskId(command string, count int, page int) (string, int, bool) {
	id, err := strconv.Atoi(command)
	if err != nil {
		log.Println("error in casting command: ", err)
		text := "Вы ввели не номер задачи в списке, а '" + command + "'"
		return text, -1, false
	}
	if id > count || id < 1 {
		log.Println(fmt.Sprintf("id not in range 1-%d", count))
		text := fmt.Sprintf("Номер задачи должен быть в интервале от 1 до %d", count)
		return text, -1, false
	}

	id = (page-1)*10 + id

	return "", id - 1, true
}
