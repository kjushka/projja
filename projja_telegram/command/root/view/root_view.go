package view

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	view2 "projja_telegram/command/execute/view"
	"projja_telegram/command/projects/view"
	rootc "projja_telegram/command/root/controller"
	"projja_telegram/command/root/menu"
	"projja_telegram/command/util"
	"strings"
)

func ListenRootCommands(botUtil *util.BotUtil) {
	for update := range botUtil.Updates {
		message := update.Message
		command := ""

		if message != nil && message.IsCommand() {
			command = message.Command()
		} else if message != nil && message.Text != "" {
			command = message.Text
		}

		switch command {
		case "start":
			msg := Start(botUtil)
			botUtil.Bot.Send(msg)
		case "Изменить навыки":
			ChangeSkills(botUtil)
		case "Обновить данные профиля":
			UpdateData(botUtil)
		case "Управлять проектами":
			view.SelectProject(botUtil)
		case "Ваши задачи":
			view2.ExecuteTasks(botUtil)
		default:
			SendUnknownMessage(botUtil)
		}

		msg := menu.GetRootMenu(botUtil.Message)
		botUtil.Bot.Send(msg)
	}
}

func Start(botUtil *util.BotUtil) tgbotapi.MessageConfig {
	user := rootc.GetUser(botUtil.Message.From.UserName)

	if user != nil {
		return menu.GetRootMenu(botUtil.Message)
	}

	text := fmt.Sprintf("Привет %s, давайте зарегистрируемся в системе", botUtil.Message.From.UserName)
	msg := getRegisterMessage(botUtil.Message, text)
	botUtil.Bot.Send(msg)

	ready := false
	for update := range botUtil.Updates {
		message := update.Message
		command := ""

		if message.Text != "" {
			command = message.Text
		}

		switch command {
		case "Регистрация":
			regMsg, status := Register(botUtil)
			if !status {
				botUtil.Bot.Send(regMsg)
				continue
			}
			msg = regMsg
			ready = true
		default:
			SendUnknownMessage(botUtil)
		}

		if ready {
			break
		}
	}

	return msg
}

func getRegisterMessage(message *util.MessageData, text string) tgbotapi.MessageConfig {
	msg := tgbotapi.NewMessage(message.Chat.ID, text)

	row := make([]tgbotapi.KeyboardButton, 0)
	btn := tgbotapi.NewKeyboardButton("Регистрация")
	row = append(row, btn)
	keyboard := tgbotapi.NewReplyKeyboard(row)

	msg.ReplyMarkup = keyboard

	return msg
}

func Register(botUtil *util.BotUtil) (tgbotapi.MessageConfig, bool) {
	status, text := rootc.RegisterUser(botUtil.Message)

	if !status {
		msg := getRegisterMessage(botUtil.Message, text)
		return msg, false
	}

	msg := tgbotapi.NewMessage(botUtil.Message.Chat.ID, text)
	botUtil.Bot.Send(msg)

	msg = SetSkills(botUtil, true)
	botUtil.Bot.Send(msg)

	msg = getUserData(botUtil.Message)
	return msg, true
}

func ChangeSkills(botUtil *util.BotUtil) {
	msg := SetSkills(botUtil, false)
	botUtil.Bot.Send(msg)

	msg = getUserData(botUtil.Message)
	botUtil.Bot.Send(msg)
}

func SetSkills(botUtil *util.BotUtil, isFirst bool) tgbotapi.MessageConfig {
	text := "Давайте теперь узнаем, что вы умеете\n" +
		"Для этого перечислите через пробел навыки, которыми вы обладаете\n" +
		"Пример: \n" +
		"frontend js angular"
	msg := tgbotapi.NewMessage(botUtil.Message.Chat.ID, text)

	if !isFirst {
		row := make([]tgbotapi.KeyboardButton, 0)
		cancelBtn := tgbotapi.NewKeyboardButton("Отмена")
		row = append(row, cancelBtn)

		keyboard := tgbotapi.NewReplyKeyboard(row)
		msg.ReplyMarkup = keyboard
	}

	var st string
	var skills []string
	for st != "success" {
		botUtil.Bot.Send(msg)
		skills, st = ListenForSkills(botUtil.Updates)
		if st == "cancel" {
			text = "Отмена обновления навыков"
			msg = tgbotapi.NewMessage(botUtil.Message.Chat.ID, text)
			return msg
		}
		if st == "error" {
			text := "Вы ввели некорректные данные"
			msg := tgbotapi.NewMessage(botUtil.Message.Chat.ID, text)
			botUtil.Bot.Send(msg)
		}
	}

	status := rootc.SetSkills(botUtil.Message.From.UserName, skills)
	if !status {
		text = "Во время регистрации навыков произошла ошибка\n" +
			"Попробуйте ввести навыки ещё раз"
		msg = tgbotapi.NewMessage(botUtil.Message.Chat.ID, text)
		botUtil.Bot.Send(msg)
		return SetSkills(botUtil, isFirst)
	}

	text = fmt.Sprintf("%s, поздравляем, ваши навыки были успешно установлены!", botUtil.Message.From.UserName)
	msg = tgbotapi.NewMessage(botUtil.Message.Chat.ID, text)
	return msg
}

func ListenForSkills(updates tgbotapi.UpdatesChannel) ([]string, string) {
	for update := range updates {
		mes := update.Message
		command := ""
		if mes.Text != "" {
			command = mes.Text
		}

		switch command {
		case "Отмена":
			return nil, "cancel"
		default:
			if command == "" {
				return nil, "error"
			}
			skills := strings.Split(command, " ")
			return skills, "success"
		}
	}
	return nil, "error"
}

func UpdateData(botUtil *util.BotUtil) {
	_, text := rootc.UpdateUserData(botUtil.Message)

	msg := tgbotapi.NewMessage(botUtil.Message.Chat.ID, text)
	botUtil.Bot.Send(msg)

	msg = getUserData(botUtil.Message)
	botUtil.Bot.Send(msg)
}

func getUserData(message *util.MessageData) tgbotapi.MessageConfig {
	user := rootc.GetUser(message.From.UserName)

	var text string

	if user == nil {
		text = "К сожалению, я не смог получить ваши данные\n" +
			"Давайте попробуем в другой раз"
		msg := tgbotapi.NewMessage(message.Chat.ID, text)

		return msg
	}

	text = fmt.Sprintf("Сейчас ваш профиль выглядит так:\n"+
		"Имя: %s\n"+
		"Username: %s\n"+
		"Навыки: %s",
		user.Name,
		user.Username,
		strings.Join(user.Skills, ", "),
	)
	msg := tgbotapi.NewMessage(message.Chat.ID, text)

	return msg
}

func SendUnknownMessage(botUtil *util.BotUtil) {
	text := "Пожалуйста, выберите один из вариантов"
	msg := tgbotapi.NewMessage(botUtil.Message.Chat.ID, text)
	botUtil.Bot.Send(msg)
}
