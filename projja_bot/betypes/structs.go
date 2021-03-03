// Этот файл не нужен, т.к. за структуры отвечает библиотека
// go-telegram-bot-api
package betypes

type BotMessage struct {
	Message struct {
		Message_id int
		From       struct {
			Id            int
			First_name    string
			Last_name     string
			Username      string
			Language_code string
		}
		Chat struct {
			Id int64
		}
		// возможно для даты стоит ипользовать другой тип данных
		Date int
		Text string
	}
}

type BotSendMessageID struct {
	Result struct {
		Message_id int
	}
}
