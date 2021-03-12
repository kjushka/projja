package main

// 	"io/ioutil"
import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"projja_bot/betypes"
	"projja_bot/logger"
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

func regiserUser(from *tgbotapi.User) {

	// TODO: нужно переделать id на int
	user := &betypes.User{
		Name: from.FirstName + " " + from.LastName,
		Username: from.UserName,
		TelegramId: from.ID, 
	}


	messageBytes, err := json.Marshal(user)
	fmt.Println(string(messageBytes))
	logger.ForError(err)

	resp, err := http.Post("http://localhost:8080/api/user/register", "application/json", bytes.NewBuffer(messageBytes))
	logger.ForError(err)

	// TODO: Можно добавить обработку ошибок
	
	// 500 ошибка может возвращаться, если ты пытаешься зарегать юзера, который уже есть в бд
	fmt.Println("Status *********")
	fmt.Println(resp.Status)

	jsonUser, err := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()
	logger.ForError(err)
	fmt.Println(string(jsonUser));

	// newUser := &betypes.User2{}
	// err = json.Unmarshal(jsonUser, newUser)
	// logger.ForError(err)

}


func checkUpdates(updates <-chan tgbotapi.Update) {
	// fmt.Println("check updates");

	for update := range updates {
		message := update.Message
		// fmt.Println("update");

		if message.IsCommand() {
			command := message.Command()
      //arguments := message.CommandArguments()

			switch command {
				case "register_user":
					fmt.Println("register user")
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
