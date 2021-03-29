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
		// fmt.Println("update");

		if message.IsCommand() {
			command := message.Command()
      //arguments := message.CommandArguments()
			fmt.Println(message.Text)

			switch command {
				case "register_user":
					logger.LogCommandResult("Register user")
					bot_commands.RegiserUser(message.From)
				case "get_user":
					logger.LogCommandResult("Get user")
				//	bot_commands.GetUser()

				default:
					fmt.Println("other command")
			}
			
		} else {
			fmt.Println("it is'n a command")
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
