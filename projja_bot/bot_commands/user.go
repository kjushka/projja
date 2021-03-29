package bot_commands

// TODO: не уверен, что это хорошее решение,
// т.к. telegram-bot-api приходится подключать
// в каждом отдельном файле, где используются библиотеги этого api
import (
	"net/http"
	"encoding/json"
	"io/ioutil"
	"bytes"
	"github.com/go-telegram-bot-api/telegram-bot-api"
	"projja_bot/betypes"
	"projja_bot/logger"
	"fmt"
)


func RegiserUser(from *tgbotapi.User) {

	// Телега возвращает id типа int
	user := &betypes.User{
		Name: from.FirstName + " " + from.LastName,
		Username: from.UserName,
		TelegramId: from.ID, 
	}

	messageBytes, err := json.Marshal(user)
	logger.LogCommandResult(string(messageBytes))
	logger.ForError(err)

	resp, err := http.Post(betypes.GetPathToMySQl("http") + "api/user/register", "application/json", bytes.NewBuffer(messageBytes))
	logger.ForError(err)

	// TODO: Нужно возаращать результат работы удалос/не удалось зарегать юзера

	// 500 ошибка может возвращаться, если ты пытаешься зарегать юзера, который уже есть в бд
	fmt.Print(resp.Status)
	logger.LogCommandResult(resp.Status)

	if(resp.StatusCode == 500) {
		jsonUser, err := ioutil.ReadAll(resp.Body)
		defer resp.Body.Close()
		logger.ForError(err)

		var duplicateUser bool = strings.HasPrefix(string(jsonUser), string(jsonUser))
		if (duplicateUser) {
			
		}

	}


}

func GetUser(userName string) {
	// Todo Валидация на наличие имени

	resp, err := http.Get(betypes.GetPathToMySQl("http") + "api/user/get/" + userName);
	fmt.Print(resp)
	logger.ForError(err)


} 