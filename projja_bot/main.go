package main

// 	"io/ioutil"
import (
	"encoding/json"
	"fmt"
	"github.com/go-telegram-bot-api/telegram-bot-api"
	"log"
	"net/http"
	"bytes"
	"projja_bot/betypes"
	"projja_bot/logger"
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

func regiserUser(from *tgbotapi.User) {
	user := &betypes.User{
		Name: from.FirstName + " " + from.LastName,
		Username: from.UserName
		TelegramId: from.ID,
	}

	// message := map[string]interface{}{
	// 	"user": "test",
	// }

	messageBytes, err := json.Marshal(user)
	if err != nil {
		log.Info(err)
	}

	resp, err := http.Post("http://localhost:8080/user/regiser", "application/json", bytes.NewBuffer(messageBytes))
	if err != nil {
		log.Info(err)
	}

	var result map[string]interface{}

	json.NewDecoder(resp.Body).Decode(&result)

	log.Println(result)
}

func checkUpdates(updates <-chan tgbotapi.Update) {

	for update := range updates {
		message := update.Message

		if message.IsCommand() {
			command := message.Command()
      //arguments := message.CommandArguments()

			switch command {
				case "register_user":
					fmt.Println("register user")
					
					fmt.Println(message.From.FirstName)
					fmt.Println(message.From.UserName)
					fmt.Println(message.From.LastName)
					fmt.Printf("%t", message.From.ID)

					regiserUser(message.From)

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
