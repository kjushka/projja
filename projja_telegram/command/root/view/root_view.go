package view

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"log"
	"projja_telegram/command/projects/view"
	rootc "projja_telegram/command/root/controller"
	"projja_telegram/command/root/menu"
	"projja_telegram/command/util"
	"strings"
)

func ListenRootCommands(botUtil *util.BotUtil) {
	for update := range botUtil.Updates {
		message := update.Message
		var command string

		if update.CallbackQuery != nil {
			response := strings.Split(update.CallbackQuery.Data, " ")
			command = response[0]
		} else if message.IsCommand() {
			command = message.Command()
		} else if message.Text != "" {
			command = message.Text
		}

		log.Println(command)
		botUtil.Message = util.MessageToMessageData(message)

		switch command {
		case "start":
			msg := Start(botUtil.Message)
			botUtil.Bot.Send(msg)
		case "register":
			Register(botUtil)
		case "set_skills":
			ChangeSkills(botUtil)
		case "update_data":
			UpdateData(botUtil)
		case "project_management":
			view.SelectProject(botUtil)
		default:
			SendUnknownMessage(botUtil, command)
		}
	}
}

func Start(message *util.MessageData) tgbotapi.MessageConfig {
	isRegister := rootc.GetUser(message.From.UserName)
	var text string

	if isRegister == nil {
		text = fmt.Sprintf("Привет %s, давайте зарегистрируемся в системе", message.From.UserName)
		return getRegisterMessage(message, text)
	} else {
		return menu.GetRootMenu(message)
	}
}

func getRegisterMessage(message *util.MessageData, text string) tgbotapi.MessageConfig {
	msg := tgbotapi.NewMessage(message.Chat.ID, text)

	keyboard := tgbotapi.InlineKeyboardMarkup{}

	var row []tgbotapi.InlineKeyboardButton
	btn := tgbotapi.NewInlineKeyboardButtonData("Регистрация", "register")
	row = append(row, btn)
	keyboard.InlineKeyboard = append(keyboard.InlineKeyboard, row)

	msg.ReplyMarkup = keyboard

	return msg
}

func Register(botUtil *util.BotUtil) {
	status, text := rootc.RegisterUser(botUtil.Message.From)

	if !status {
		msg := getRegisterMessage(botUtil.Message, text)
		botUtil.Bot.Send(msg)

		return
	}

	msg := tgbotapi.NewMessage(botUtil.Message.Chat.ID, text)
	botUtil.Bot.Send(msg)

	msg = SetSkills(botUtil, true)
	botUtil.Bot.Send(msg)

	msg = menu.GetRootMenu(botUtil.Message)
	botUtil.Bot.Send(msg)
}

func ChangeSkills(botUtil *util.BotUtil) {
	defer func(message *util.MessageData, bot *tgbotapi.BotAPI) {
		msg := menu.GetRootMenu(message)
		bot.Send(msg)
	}(botUtil.Message, botUtil.Bot)

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
		keyboard := tgbotapi.InlineKeyboardMarkup{}
		row := make([]tgbotapi.InlineKeyboardButton, 0)
		cancelBtn := tgbotapi.NewInlineKeyboardButtonData("Отмена", "cancel_btn")

		row = append(row, cancelBtn)
		keyboard.InlineKeyboard = append(keyboard.InlineKeyboard, row)
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
	_, text := rootc.UpdateUserData(botUtil.Message.From)

	msg := tgbotapi.NewMessage(botUtil.Message.Chat.ID, text)
	botUtil.Bot.Send(msg)

	msg = getUserData(botUtil.Message)
	botUtil.Bot.Send(msg)

	msg = menu.GetRootMenu(botUtil.Message)
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

func SendUnknownMessage(botUtil *util.BotUtil, command string) {
	text := fmt.Sprintf("Я не знаю команды '%s'", command)
	msg := tgbotapi.NewMessage(botUtil.Message.Chat.ID, text)
	botUtil.Bot.Send(msg)

	msg = menu.GetRootMenu(botUtil.Message)
	botUtil.Bot.Send(msg)
}
