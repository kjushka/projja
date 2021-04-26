package view

import (
	"fmt"
	// "projja_bot/betypes"
	// "projja_bot/logger"
	"projja_bot/logger"
	"projja_bot/bot_commands/controller"
	"github.com/go-telegram-bot-api/telegram-bot-api"
	// "net/http"
	// "strings"
)

func ChooseProjjaAction(message *tgbotapi.Message) tgbotapi.MessageConfig {
	text := fmt.Sprintf("%s , что вы хотите сделать?\n" +
	"Чтобы создать проект, воспользуйтес командой:\n" +
	"/create_project название проекта", message.From.UserName)

	msg := tgbotapi.NewMessage(message.Chat.ID, text)

	keyboard := tgbotapi.InlineKeyboardMarkup{}

	var row1 []tgbotapi.InlineKeyboardButton
	var row2 []tgbotapi.InlineKeyboardButton
	setBtn := tgbotapi.NewInlineKeyboardButtonData("Настроить профиль", "register_user")
	dirProjBtn := tgbotapi.NewInlineKeyboardButtonData("Управлять проектами", "project_control")
	dirTaskBtn := tgbotapi.NewInlineKeyboardButtonData("Управлять задачами", "register_user")

	row1 = append(row1, setBtn)
	row2 = append(row2, dirProjBtn)
	row2 = append(row2, dirTaskBtn)
	keyboard.InlineKeyboard = append(keyboard.InlineKeyboard, row2)
	keyboard.InlineKeyboard = append(keyboard.InlineKeyboard, row1)

	msg.ReplyMarkup = keyboard
	return msg
}


func Start(message *tgbotapi.Message) tgbotapi.MessageConfig {
	_, isRegister := controller.GetUser(message.From.UserName)
	var text string

	if isRegister == nil {
		text = fmt.Sprintf("Привет %s, давай зарегистрируемся в системе:)", message.From.UserName)
		msg := tgbotapi.NewMessage(message.Chat.ID, text)
	
		keyboard := tgbotapi.InlineKeyboardMarkup{}

		var row []tgbotapi.InlineKeyboardButton
		btn := tgbotapi.NewInlineKeyboardButtonData("Регистрация", "register_user")
		row = append(row, btn)
		keyboard.InlineKeyboard = append(keyboard.InlineKeyboard, row)

		msg.ReplyMarkup = keyboard
		return msg
	} else {
		return ChooseProjjaAction(message);
	}
	// Проверка указаны ли скилы
}

func Register(message *tgbotapi.Message) (tgbotapi.MessageConfig, tgbotapi.MessageConfig) {
	var ans string = controller.RegiserUser(message.From)

	msg := tgbotapi.NewMessage(message.Chat.ID, ans)
	msg.ReplyToMessageID = message.MessageID
	
	ans = "Давай теперь узнаем, что ты умеешь:) \n" + 
	"Для этого введи команду /set_skills и перечисли через пробел навыки, которыми ты обладаешь \n" +
	"Пример: \n" + 
	"/set_skills frontend js angular"
	
	return msg, tgbotapi.NewMessage(message.Chat.ID, ans)
}

func CreateProject(message *tgbotapi.Message) tgbotapi.MessageConfig {
	logger.LogCommandResult("Create project");
	var ans string = controller.CreateProject(message.From.UserName, message.CommandArguments())

	msg := tgbotapi.NewMessage(message.Chat.ID, ans)
	msg.ReplyToMessageID = message.MessageID
	return msg
}

func GetAllProjects(message *tgbotapi.Message) tgbotapi.MessageConfig {
	logger.LogCommandResult("Get all projects");
	keyboard, countPrjects := controller.GetAllProjects(message.From.UserName)

	if	countPrjects == 0 {
		return tgbotapi.NewMessage(message.Chat.ID, "На данный момент у вас нет открытых проектов:(")	
	}

	msg := tgbotapi.NewMessage(message.Chat.ID, "Выберите проект с которым вы хотите работать:")
	msg.ReplyMarkup = keyboard	
	return msg
}

func SetSkills(message *tgbotapi.Message) tgbotapi.MessageConfig {
	ans := controller.SetSkills(message.From.UserName, message.CommandArguments())
	return tgbotapi.NewMessage(message.Chat.ID, ans)
}