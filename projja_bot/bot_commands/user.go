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
	"strconv"
)

// Данная функция возвращает ошибку или сообщение об 
// удачной регистрации пользователя
func RegiserUser(from *tgbotapi.User) string {
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
		return fmt.Sprintf("Пользователь %s был успешно зарегистрирован!", from.UserName);
	}

	logger.LogCommandResult("Non-standard situation during registration.")
	return "Что-то пошло не так..." 
}

func GetUser(args string) (string, *betypes.User) {
	if(args == "") {
		return "Вы не указали имя пользавателя!", nil  
	}

	// Берем только первый аргумент
	// возможно тут нужна более хорошая валидация
	var userName string = strings.Split(args, " ")[0]

	resp, err := http.Get(betypes.GetPathToMySQl("http") + "api/user/get/" + userName);
	logger.ForError(err)

	getUserInfo, err := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()
	logger.ForError(err)
	logger.LogCommandResult(string(getUserInfo));

	var userAns *betypes.GetUserAnswer
	if err := json.Unmarshal(getUserInfo, &userAns); err != nil {
    logger.ForError(err)
	}

	if userAns.IsEmpty == true {
		return fmt.Sprintf("Пользоватль с именем %s не зарегистрирован!", userName), nil 
	}

	// Если пользователь есть, то нужно его куда-то сохранить
	return fmt.Sprintf("Вы выбрали пользователя %s", userName), userAns.Content
} 

// ПРИМЕЧАНИЕ
// если вы добавляете юзеру навык, который у него уже есть, 
// то все равно получите код 202 (все хорошо), хотя логичней
// было бы сообщить о том, что такой навык уже есть
// UPD Как сказал сан Антон. Алгоритм добавления навыков служит для координального изменения их
// по этому эта команда стирает все навыки и записывает новые 
func SetSkills(args string) string {	
	if (args == "") {
		return "Вы не указали имя пользователя и его навыки!"
	} 
	
	var argsArr[] string = strings.Split(args, " ")
	if (len(argsArr) == 1){
		return "Вы не указали навыки, которые хотите присвоить пользователю!"
	}

	userName := argsArr[0]
	skills := argsArr[1: len(argsArr)]

	userSkills := &betypes.Skills {
		Skills: skills,
	}

	skillsBytes, err := json.Marshal(userSkills)
	logger.LogCommandResult(string(skillsBytes))
	logger.ForError(err)

	resp, err :=  http.Post(betypes.GetPathToMySQl("http") + fmt.Sprintf("api/user/%s/skills", userName), "application/json", bytes.NewBuffer(skillsBytes) )
	logger.ForError(err)
	logger.LogCommandResult(resp.Status)

	if resp.StatusCode == 404 || resp.StatusCode == 500 {
		return fmt.Sprintf("Пользователь с именем %s не зарегистрирован!", userName)
	}

	if resp.StatusCode == 202 {
		return "Навыки были успешно установленн!"
	}

	return "Что-то пошло не так..."
}


// Нужна ли эта функция
// func ChangeName()

func GetAllProjects(args string) string {
	if (args == "") {
		return "Вы не указали имя пользователя, проекты которого хотите просмотреть!"
	} 

	var userName string = strings.Split(args, " ")[0]

	resp, err := http.Get(betypes.GetPathToMySQl("http") + fmt.Sprintf("api/user/%s/owner/all", userName))
	logger.ForError(err)

	fmt.Println(resp.Status)

	gettingProjects, err := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()
	logger.ForError(err)

	
	var projects *betypes.ProjectsList
	if err := json.Unmarshal(gettingProjects, &projects); err != nil {
    logger.ForError(err)
	}

	ans := fmt.Sprintf("User %s have projects:\nId Project\n", userName)
	for i := 0; i < len(projects.Content); i++ {
		ans += strconv.FormatInt(projects.Content[i].Id, 10) + "  " + projects.Content[i].Name + "\n"
	}

	fmt.Println(ans)
	

	return ans
}