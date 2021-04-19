package main

import (
	"fmt"
	"projja_bot/betypes"
	"projja_bot/logger"
	"projja_bot/bot_commands"
	"github.com/go-telegram-bot-api/telegram-bot-api"
	"net/http"
)

var (
	Bot, BotErr = tgbotapi.NewBotAPI(betypes.BotToken)
)

func setWebhook(bot *tgbotapi.BotAPI) {
	webHookInfo := tgbotapi.NewWebhookWithCert(fmt.Sprintf("https://%s:%s/%s", betypes.BotExternalAddress, betypes.BotExternalPort,
		betypes.BotToken), betypes.CertPath)
	_, err := bot.SetWebhook(webHookInfo)
	logger.ForError(err)
}

func checkUpdates(updates <-chan tgbotapi.Update) {
	// fmt.Println("check updates");

	for update := range updates {
		message := update.Message
		var command string

		if update.CallbackQuery != nil {
			command = update.CallbackQuery.Data
			message = update.CallbackQuery.Message
			// Подменяем from bota на from пользователя нажавшего кнопку
			message.From = update.CallbackQuery.From
		} else if message.IsCommand() {
			command = message.Command()
		}

		switch command {
			case "start":
				_, isRegister := bot_commands.GetUser(message.From.UserName)
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
					Bot.Send(msg)
				} else {
					text = fmt.Sprintf("Добрый день %s! Что вы хотите сделать?\n" +
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
					Bot.Send(msg)
				}
				// Проверка указаны ли скилы

			case "register_user":
				logger.LogCommandResult("Register user")
				var ans string = bot_commands.RegiserUser(message.From)

				msg := tgbotapi.NewMessage(message.Chat.ID, ans)
				msg.ReplyToMessageID = message.MessageID
				Bot.Send(msg)
				
				ans = "Давай теперь узнаем, что ты умеешь:) \n" + 
				"Для этого введи команду /set_skills и перечисли через пробел навыки, которыми ты обладаешь \n" +
				"Пример: \n" + 
				"/set_skills frontend js angular"
				msg = tgbotapi.NewMessage(message.Chat.ID, ans)
				Bot.Send(msg)
			
			case "create_project":	
				logger.LogCommandResult("Create project");
				var ans string = bot_commands.CreateProject(message.From.UserName, message.CommandArguments())

				msg := tgbotapi.NewMessage(message.Chat.ID, ans)
				msg.ReplyToMessageID = message.MessageID
				Bot.Send(msg)
			case "change_profile":
				// var text string;
				// text = "Для изменения настроек пользователя вы можете использовать следующие команды\n"
				// +	"/set_skills навык1 навык2 ... навыкN - изменить навыки пользователя\n" 
				// + "/change_name новое имя - " 
				// TODO 

			case "project_control":
				logger.LogCommandResult("Get all projects");
				keyboard, countPrjects := bot_commands.GetAllProjects(message.From.UserName)

				if	countPrjects == 0 {
					msg := tgbotapi.NewMessage(message.Chat.ID, "На данный момент у вас нет открытых проектов:(")
					Bot.Send(msg)
				} else {
					msg := tgbotapi.NewMessage(message.Chat.ID, "Выберите проект с которым вы хотите работать:")
					msg.ReplyMarkup = keyboard
					Bot.Send(msg)
				}

			}
	
	}	
}

func main() {
	logger.ForError(BotErr)
	setWebhook(Bot)
	
	updates := Bot.ListenForWebhook("/")

	fmt.Println("Server is working!")
	go http.ListenAndServeTLS(fmt.Sprintf("%s:%s", betypes.BotInternalAddress, betypes.BotInternalPort),
		betypes.CertPath, betypes.KeyPath, nil)

	checkUpdates(updates)
}
