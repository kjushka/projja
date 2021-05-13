package util

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"strings"
)

func GetUnknownMessage(botUtil *BotUtil, command string) tgbotapi.MessageConfig {
	text := fmt.Sprintf("Я не знаю команды '%s'", command)
	msg := tgbotapi.NewMessage(botUtil.Message.Chat.ID, text)
	return msg
}

func ListenForText(botUtil *BotUtil, mesText string, cancelText string) (string, bool) {
	msg := tgbotapi.NewMessage(botUtil.Message.Chat.ID, mesText)

	keyboard := tgbotapi.InlineKeyboardMarkup{}
	row := make([]tgbotapi.InlineKeyboardButton, 0)
	cancelBtn := tgbotapi.NewInlineKeyboardButtonData("Отмена", "cancel_btn")
	row = append(row, cancelBtn)
	keyboard.InlineKeyboard = append(keyboard.InlineKeyboard, row)

	msg.ReplyMarkup = keyboard

	botUtil.Bot.Send(msg)

	resultText := ""
	for update := range botUtil.Updates {
		mes := update.Message
		var command string

		if update.CallbackQuery != nil {
			response := strings.Split(update.CallbackQuery.Data, " ")
			command = response[0]
		} else if mes.IsCommand() {
			command = mes.Command()
		} else if mes.Text != "" {
			command = mes.Text
		}

		switch command {
		case "cancel_btn":
			resultText = cancelText
			return resultText, false
		default:
			if command == "" {
				continue
			}
			resultText = command
		}

		if resultText != "" {
			break
		}
	}

	return resultText, true
}
