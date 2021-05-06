package main

import (
	"fmt"
	"net/http"
	"projja_bot/betypes"
	"projja_bot/bot_commands/view"
	"projja_bot/logger"
	"strings"

	"github.com/go-telegram-bot-api/telegram-bot-api"
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
	for update := range updates {
		message := update.Message
		var command string
		var args[] string

		if update.CallbackQuery != nil { 
			// Если произошло нажатие на inlain кнопку, то отделяем команду от её аргументов
			response := strings.Split(update.CallbackQuery.Data, " ")
			command = response[0]
			args = response[1: len(response)]

			// Подменяем from bota на from пользователя нажавшего кнопку
			message = update.CallbackQuery.Message
			message.From = update.CallbackQuery.From
		} else if message.IsCommand() {
			command = message.Command()
		}

		fmt.Println(command)

		switch command {
			case "start":
				msg := view.Start(message)
				Bot.Send(msg)
			case "register_user":
				msg1, msg2 := view.Register(message)
				Bot.Send(msg1)
				Bot.Send(msg2)
			case "set_skills":
				msg := view.SetSkills(message)
				Bot.Send(msg)
				msg = view.ChooseProjjaAction(message)
				Bot.Send(msg)
			case "create_project":	
				msg := view.CreateProject(message);
				Bot.Send(msg)
				msg = view.ChooseProjjaAction(message)
				Bot.Send(msg)
			case "change_profile":
				// var text string;
				// text = "Для изменения настроек пользователя вы можете использовать следующие команды\n"
				// +	"/set_skills навык1 навык2 ... навыкN - изменить навыки пользователя\n" 
				// + "/change_name новое имя - " 
				// TODO 

			case "project_control":
				// Тут обрабатывается логика, выбора проекта с которым хочет работать пользователь
				msg := view.GetAllProjects(message);
				Bot.Send(msg)				
			case "select_project":
				// Тут находится логика, которую можно выполнить после выбора проекта
				projectId := args[0]
				selectedProject := args[1]
				

				msg := view.SelectProject(message, projectId, selectedProject)
				Bot.Send(msg)
				msg = view.ChosePrjectAction(message);
				Bot.Send(msg)
			case "members_management":
				msg := view.MembersManagment(message)
				Bot.Send(msg)	
			case "add_member":
				msg := view.AddMemberToProject(message)	
				Bot.Send(msg)	
			case "add_member_yes":
				msg := view.AddMemberYes(message)
				Bot.Send(msg)	
			case "add_member_no":
				msg := view.AddMemberNo(message)
				Bot.Send(msg)	
			case "get_members":
				msg := view.GetProjectMembers(message)
				Bot.Send(msg)	
			case "remove_member":
				msg := view.RemoveMemberFromProject(message)
				Bot.Send(msg)		
			}
	
	}	
}

func main() {
	logger.ForError(BotErr)
	setWebhook(Bot)
	
	updates := Bot.ListenForWebhook("/")

	fmt.Println("Server is working!")
	go http.ListenAndServeTLS(fmt.Sprintf("%s:%s", betypes.BotInternalAddress, betypes.BotInternalPort), betypes.CertPath, betypes.KeyPath, nil)

	checkUpdates(updates)
}
