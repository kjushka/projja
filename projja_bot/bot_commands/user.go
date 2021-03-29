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
	"strings"
)

// Данная функция возвращает ошибку или сообщение об 
// удачной регистрации пользователя
func RegiserUser(from *tgbotapi.User) string {
	// Телега возвращает id типа int
	var userName string = from.UserName
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

	// TODO: Нужно возаращать результат работы удалось/не удалось зарегать юзера
	// 500 ошибка может возвращаться, если ты пытаешься зарегать юзера, который уже есть в бд
	fmt.Print(resp.Status)
	logger.LogCommandResult(resp.Status)

	if(resp.StatusCode >= 500) {
		jsonUser, err := ioutil.ReadAll(resp.Body)
		defer resp.Body.Close()
		logger.ForError(err)

		var duplicateUser bool = strings.HasPrefix(string(jsonUser), string(jsonUser))
		if (duplicateUser) {
			return "Такой пользователь уже зарегистрирован!"
		}

		return "Неизвестная ошибка"
	} else if (resp.StatusCode < 300) {
		fmt.Print(userName)
		fmt.Print(from.UserName)
		return fmt.Sprintf("Пользователь %s был успешно зарегистрирован!", from.UserName);
	}

	logger.LogCommandResult("Non-standard situation during registration.")
	return "Что-то пошло не так..." 
}

func GetUser(userName string) {
	fmt.Println(userName)

	resp, err := http.Get(betypes.GetPathToMySQl("http") + "api/user/get/" + userName);
	fmt.Println(resp.StatusCode)
	logger.ForError(err)

	getUserInfo, err := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()
	logger.ForError(err)
	fmt.Println(string(getUserInfo))

} 