package controller

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/go-martini/martini"
	"io/ioutil"
	"log"
	"net/http"
	"projja_api/model"
	"strconv"
	"time"
)

func (c *Controller) AddAnswer(w http.ResponseWriter, r *http.Request) (int, string) {
	contentType := r.Header.Get("Content-Type")
	if contentType != "application/json" {
		err := fmt.Sprintf("Unsupportable Content-Type header: %s", contentType)
		log.Println(err)
		return 500, err
	}
	jsonAnswer, err := ioutil.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		log.Println("error in reading body", err)
		return 500, err.Error()
	}
	answer := &model.Answer{}
	err = json.Unmarshal(jsonAnswer, answer)
	if err != nil {
		log.Println("error in unmarshalling")
		return 500, err.Error()
	}

	result, err := c.DB.Exec(
		"insert into answer (task, executor, message_id, chat_id, status, sent_at) values (?, ?, ?, ?, ?, ?)",
		answer.Task.Id,
		answer.Executor.Id,
		answer.MessageId,
		answer.ChatId,
		answer.Status,
		answer.SentAt,
	)
	if err != nil {
		log.Println("error in creating answer:", err)
		return 500, err.Error()
	}

	rowsAffected, _ := result.RowsAffected()

	w.Header().Set("Content-Type", "application/json")
	return c.makeContentResponse(201, "Answer created", struct {
		Name    string
		Content interface{}
	}{
		Name:    "Rows affected",
		Content: rowsAffected,
	})
}

func (c *Controller) GetLastAnswer(params martini.Params, w http.ResponseWriter) (int, string) {
	user, err := c.getUserByUsername(params["uname"])
	if err != nil {
		log.Println("error in getting user: ", err)
		return 500, err.Error()
	}
	taskId, err := strconv.ParseInt(params["tid"], 10, 64)
	if err != nil {
		log.Println("error in parsing taskId: ", err)
		return 500, err.Error()
	}

	row := c.DB.QueryRow(
		"select * from answer where task = ? and executor = ? order by sent_at desc limit 1",
		taskId,
		user.Id,
	)
	answerDto := &struct {
		Id        int64
		Task      int64
		Executor  int64
		MessageId int
		ChatId    int64
		Status    string
		SentAt    time.Time
	}{}
	err = row.Scan(
		&answerDto.Id,
		&answerDto.Task,
		&answerDto.Executor,
		&answerDto.MessageId,
		&answerDto.ChatId,
		&answerDto.Status,
		&answerDto.SentAt,
	)
	if err != nil && err != sql.ErrNoRows {
		log.Println("error in getting answer: ", err)
		return 500, err.Error()
	}
	if err == sql.ErrNoRows {
		log.Println("no answers for task added yet")
		return 404, err.Error()
	}

	task, _ := c.getTaskById(taskId)

	answer := &model.Answer{
		Id:        answerDto.Id,
		Task:      task,
		Executor:  user,
		MessageId: answerDto.MessageId,
		ChatId:    answerDto.ChatId,
		Status:    answerDto.Status,
		SentAt:    answerDto.SentAt,
	}

	w.Header().Set("Content-Type", "application/json")
	return c.makeContentResponse(200, "Last answer", answer)
}
