package main

import (
	"github.com/go-telegram-bot-api/telegram-bot-api"
	"log"
	"os"
	rootv "projja_telegram/command/root/view"
	"projja_telegram/command/util"
)

var UsersChan map[int]*util.BotUtil

func main() {
	botToken := os.Getenv("BOT_TOKEN")
	bot, err := tgbotapi.NewBotAPI(botToken)
	if err != nil {
		log.Panic(err)
	}

	bot.Debug = true

	log.Printf("Authorized on account %s\n", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates, err := bot.GetUpdatesChan(u)

	UsersChan = make(map[int]*util.BotUtil)

	for update := range updates {
		var from *tgbotapi.User
		var chat *tgbotapi.Chat
		if update.Message != nil {
			from = update.Message.From
			chat = update.Message.Chat

			log.Println(update.Message.Text)
		} else if update.CallbackQuery != nil {
			from = update.CallbackQuery.From
			chat = update.CallbackQuery.Message.Chat

			log.Println(update.CallbackQuery)
		}

		if _, ok := UsersChan[from.ID]; !ok {
			UsersChan[from.ID] = &util.BotUtil{
				Message: &util.MessageData{
					From: from,
					Chat: chat,
				},
				Bot:     bot,
				Updates: make(chan tgbotapi.Update),
			}

			go rootv.ListenRootCommands(UsersChan[from.ID])
		}

		UsersChan[from.ID].Updates <- update
	}
}
