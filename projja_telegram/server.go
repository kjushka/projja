package main

import (
	"github.com/go-telegram-bot-api/telegram-bot-api"
	"log"
	"os"
	rootv "projja_telegram/command/root/view"
	"projja_telegram/command/util"
)

//1854133506:AAFi2RLmybsgjAuNQtB207xsXaRqiIaipm8

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

	botUtil := &util.BotUtil{
		Message: nil,
		Bot:     bot,
		Updates: updates,
	}

	rootv.ListenRootCommands(botUtil)
}
