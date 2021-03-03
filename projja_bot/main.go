package main

import (
	"encoding/json"
	"fmt"
	"github.com/go-telegram-bot-api/telegram-bot-api"
	"io/ioutil"
	"log"
	"net/http"
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

func main() {
	// log.Printf("Autorized on account %s", Bot.Self.UserName)
	logger.ForError(BotErr)
	setWebhook(Bot)

	message := func(w http.ResponseWriter, r *http.Request) {
		text, err := ioutil.ReadAll(r.Body)
		logger.ForError(err)

		var botText betypes.BotMessage
		err = json.Unmarshal(text, &botText)
		logger.ForError(err)

		fmt.Println(fmt.Sprintf("%s", text))
		logger.LogFile.Println(fmt.Sprintf("%s", text))

		firstName := botText.Message.From.First_name
		// userName := botText.Message.From.Username
		chatGroup := botText.Message.Chat.Id
		// mText := botText.Message.Text
		msg := tgbotapi.NewMessage(chatGroup, fmt.Sprintf("Привет, %s", firstName))

		Bot.Send(msg)
	}

	http.HandleFunc("/", message)
	fmt.Println("Server is working!")
	log.Fatal(http.ListenAndServeTLS(fmt.Sprintf("%s:%s", betypes.BotInternalAddress, betypes.BotInternalPort),
		betypes.CertPath, betypes.KeyPath, nil))
}
