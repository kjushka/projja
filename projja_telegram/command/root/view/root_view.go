package view

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"log"
	"projja_telegram/command/projects"
	rootc "projja_telegram/command/root/controller"
	"projja_telegram/command/util"
	"strings"
)

func ListenRootCommands(bot *tgbotapi.BotAPI, updates tgbotapi.UpdatesChannel) {
	for update := range updates {
		message := update.Message
		var command string

		if update.CallbackQuery != nil {
			response := strings.Split(update.CallbackQuery.Data, " ")
			command = response[0]

			message = update.CallbackQuery.Message
			message.From = update.CallbackQuery.From
		} else if message.IsCommand() {
			command = message.Command()
		} else if message.Text != "" {
			command = message.Text
		}

		log.Println(command)
		//log.Printf("[%s] %s\n", update.Message.From.UserName, update.Message.Text)
		messageData := util.MessageToMessageData(message)

		switch command {
		case "start":
			msg := Start(messageData)
			bot.Send(msg)
		case "register":
			Register(messageData, bot, updates)
		case "set_skills":
			ChangeSkills(messageData, bot, updates)
		case "project_management":
			projects.SelectProject(messageData, bot, updates)
		default:
			SendUnknownMessage(messageData, command, bot)
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
		return GetRootMenu(message)
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

func Register(message *util.MessageData, bot *tgbotapi.BotAPI, updates tgbotapi.UpdatesChannel) {
	status, text := rootc.RegisterUser(message.From)

	if !status {
		msg := getRegisterMessage(message, text)
		bot.Send(msg)

		return
	}

	msg := tgbotapi.NewMessage(message.Chat.ID, text)
	bot.Send(msg)

	SetSkills(message, bot, updates)

	msg = GetRootMenu(message)
	bot.Send(msg)
}

func SetSkills(message *util.MessageData, bot *tgbotapi.BotAPI, updates tgbotapi.UpdatesChannel) {
	text := "Давайте теперь узнаем, что вы умеете\n" +
		"Для этого перечислите через пробел навыки, которыми вы обладаете\n" +
		"Пример: \n" +
		"frontend js angular"
	msg := tgbotapi.NewMessage(message.Chat.ID, text)
	bot.Send(msg)

	status := false
	var skills []string
	for !status {
		skills, status = ListenForSkills(updates)
	}

	status = rootc.SetSkills(message.From.UserName, skills)
	if !status {
		text = "Во время регистрации навыков произошла ошибка\n" +
			"Попробуйте ввести навыки ещё раз"
		msg = tgbotapi.NewMessage(message.Chat.ID, text)
		bot.Send(msg)
		SetSkills(message, bot, updates)
	}

	text = fmt.Sprintf("%s, поздравляем, ваши навыки были успешно установлены!", message.From.UserName)
	msg = tgbotapi.NewMessage(message.Chat.ID, text)
	bot.Send(msg)
}

func ListenForSkills(updates tgbotapi.UpdatesChannel) ([]string, bool) {
	for update := range updates {
		message := update.Message
		if message != nil {
			skillsString := message.Text
			if skillsString == "" {
				return nil, false
			}
			skills := strings.Split(skillsString, " ")
			return skills, true
		}
	}
	return nil, false
}

func ChangeSkills(message *util.MessageData, bot *tgbotapi.BotAPI, updates tgbotapi.UpdatesChannel) {
	defer func(message *util.MessageData, bot *tgbotapi.BotAPI) {
		msg := GetRootMenu(message)
		bot.Send(msg)
	}(message, bot)

	user := rootc.GetUser(message.From.UserName)

	var text string

	if user == nil {
		text = "К сожалению, я не смог получить ваши данные\n" +
			"Давайте попробуем в другой раз"
		msg := tgbotapi.NewMessage(message.Chat.ID, text)
		bot.Send(msg)

		return
	}

	SetSkills(message, bot, updates)

	text = fmt.Sprintf("Сейчас ваш профиль выглядит так:\n"+
		"Имя: %s\n"+
		"Username: %s\n"+
		"Навыки: %s",
		user.Name,
		user.Username,
		strings.Join(user.Skills, ", "),
	)
	msg := tgbotapi.NewMessage(message.Chat.ID, text)
	bot.Send(msg)
}

func SendUnknownMessage(message *util.MessageData, command string, bot *tgbotapi.BotAPI) {
	text := fmt.Sprintf("Я не знаю команды '%s'", command)
	msg := tgbotapi.NewMessage(message.Chat.ID, text)
	bot.Send(msg)

	msg = GetRootMenu(message)
	bot.Send(msg)
}

func GetRootMenu(message *util.MessageData) tgbotapi.MessageConfig {
	text := fmt.Sprintf("%s, что вы хотите сделать?\n", message.From.UserName)
	msg := tgbotapi.NewMessage(message.Chat.ID, text)

	keyboard := tgbotapi.InlineKeyboardMarkup{}

	var row []tgbotapi.InlineKeyboardButton
	skillsBtn := tgbotapi.NewInlineKeyboardButtonData("Изменить навыки", "set_skills")
	projectsManageBtn := tgbotapi.NewInlineKeyboardButtonData("Управлять проектами", "project_management")

	row = append(row, skillsBtn)
	row = append(row, projectsManageBtn)
	keyboard.InlineKeyboard = append(keyboard.InlineKeyboard, row)

	msg.ReplyMarkup = keyboard
	return msg
}
