package controller

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"projja_telegram/config"
	"projja_telegram/model"
	"time"
)

func CalculateExecutor(project *model.Project, task *model.Task) (*model.User, error) {
	jsonTask, err := json.Marshal(task)
	if err != nil {
		log.Println("error in marshalling task: ", err)
		return nil, err
	}

	resp, err := http.Post(config.GetExecAddr()+fmt.Sprintf("/project/%d/calc/task", project.Id),
		"application/json",
		bytes.NewBuffer(jsonTask),
	)
	if err != nil {
		log.Println("error in request for calculating task executor: ", err)
		return nil, err
	}

	respData := &struct {
		Description string
		Content     *model.Task
	}{}
	jsonBody, err := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()
	if err != nil {
		log.Println("error in reading response body: ", err)
		return nil, err
	}

	err = json.Unmarshal(jsonBody, respData)
	if err != nil {
		log.Println("error in unmarshalling user: ", err)
		return nil, err
	}

	return respData.Content.Executor, nil
}

func CreateTask(project *model.Project, task *model.Task) (string, bool) {
	errorText := "Во время создания задачи произошла ошибка\nПопробуйте позже ещё раз"

	jsonTask, err := json.Marshal(task)
	if err != nil {
		log.Println("error in marshalling task: ", err)
		return errorText, false
	}

	resp, err := http.Post(config.GetAPIAddr()+fmt.Sprintf("/project/%d/create/task", project.Id),
		"application/json",
		bytes.NewBuffer(jsonTask),
	)
	if err != nil {
		log.Println("error in request for creating task: ", err)
		return errorText, false
	}

	if resp.StatusCode == http.StatusInternalServerError {
		log.Println("error in request for creating task")
		return errorText, false
	}

	return "Задача успешно создана", true
}

func ChangeTaskDescription(task *model.Task, newDescription string) (string, bool) {
	descriptionStruct := &struct {
		Description string
	}{newDescription}

	errorText := "Во время смены описания задачи произошла ошибка\nПопробуйте позже ещё раз"

	jsonDescriptionStruct, err := json.Marshal(descriptionStruct)
	if err != nil {
		log.Println("error in marshalling description: ", err)
		return errorText, false
	}

	resp, err := http.Post(config.GetAPIAddr()+fmt.Sprintf("/task/%d/change/description", task.Id),
		"application/json",
		bytes.NewBuffer(jsonDescriptionStruct),
	)
	if err != nil {
		log.Println("error in request for changing task description: ", err)
		return errorText, false
	}

	if resp.StatusCode == http.StatusInternalServerError {
		log.Println("error in request for changing task description")
		return errorText, false
	}

	task.Description = newDescription
	return "Описание задачи успешно изменено", true
}

func ChangeTaskDeadline(task *model.Task, newDeadline time.Time) (string, bool) {
	deadlineStruct := &struct {
		Deadline string
	}{newDeadline.Format("2006-01-02")}

	errorText := "Во время изменения дедлайна задачи произошла ошибка\nПопробуйте позже ещё раз"

	jsonDeadlineStruct, err := json.Marshal(deadlineStruct)
	if err != nil {
		log.Println("error in marshalling deadline: ", err)
		return errorText, false
	}

	resp, err := http.Post(config.GetAPIAddr()+fmt.Sprintf("/task/%d/change/deadline", task.Id),
		"application/json",
		bytes.NewBuffer(jsonDeadlineStruct),
	)
	if err != nil {
		log.Println("error in request for changing task deadline: ", err)
		return errorText, false
	}

	if resp.StatusCode == http.StatusInternalServerError {
		log.Println("error in request for changing task deadline")
		return errorText, false
	}

	task.Deadline = newDeadline.Format("2006-01-02")
	return "Дедлайн задачи успешно изменен", true
}

func ChangeTaskPriority(task *model.Task, newPriority string) (string, bool) {
	priorityStruct := &struct {
		Priority string
	}{newPriority}

	errorText := "Во время изменения приоритета задачи произошла ошибка\nПопробуйте позже ещё раз"

	jsonPriorityStruct, err := json.Marshal(priorityStruct)
	if err != nil {
		log.Println("error in marshalling priority: ", err)
		return errorText, false
	}

	resp, err := http.Post(config.GetAPIAddr()+fmt.Sprintf("/task/%d/change/priority", task.Id),
		"application/json",
		bytes.NewBuffer(jsonPriorityStruct),
	)
	if err != nil {
		log.Println("error in request for changing task priority: ", err)
		return errorText, false
	}

	if resp.StatusCode == http.StatusInternalServerError {
		log.Println("error in request for changing task priority")
		return errorText, false
	}

	task.Priority = newPriority
	return "Приоритет задачи успешно изменен", true
}

func ChangeTaskExecutor(task *model.Task, newExecutor *model.User) (string, bool) {
	errorText := "Во время изменения дедлайна задачи произошла ошибка\nПопробуйте позже ещё раз"

	jsonExecutor, err := json.Marshal(newExecutor)
	if err != nil {
		log.Println("error in marshalling deadline: ", err)
		return errorText, false
	}

	resp, err := http.Post(config.GetAPIAddr()+fmt.Sprintf("/task/%d/change/executor", task.Id),
		"application/json",
		bytes.NewBuffer(jsonExecutor),
	)
	if err != nil {
		log.Println("error in request for changing task executor: ", err)
		return errorText, false
	}

	if resp.StatusCode == http.StatusInternalServerError {
		log.Println("error in request for changing task executor")
		return errorText, false
	}

	task.Executor = newExecutor
	return "Исполнитель задачи успешно изменен", true
}

func CloseTask(task *model.Task) (string, bool) {
	url := config.GetAPIAddr() + fmt.Sprintf("/task/%d/close", task.Id)

	errorText := "Во время закрытия задачи произошла ошибка\nПопробуйте позже ещё раз"

	resp, err := http.Get(url)
	if err != nil {
		log.Println("error in request for closing task: ", err)
		return errorText, false
	}

	if resp.StatusCode == http.StatusInternalServerError {
		log.Println("error in request for closing task")
		return errorText, false
	}

	task.IsClosed = true
	return "Задача успешно закрыта", true
}
