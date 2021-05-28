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
		"select * from answer where task = ? and executor = ? and status in (?, ?) order by sent_at desc limit 1",
		taskId,
		user.Id,
		"not checked",
		"declined",
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

func (c *Controller) GetProjectAnswers(params martini.Params, w http.ResponseWriter) (int, string) {
	projectId, err := strconv.ParseInt(params["pid"], 10, 64)
	if err != nil {
		log.Println("error in parse project id: ", err.Error())
		return 500, err.Error()
	}

	rows, err := c.DB.Query(
		"select max(a.id), a.task, t.description, cast(t.deadline as char), t.priority, "+
			"t.status, t.level, max(a.message_id), max(a.chat_id), max(a.status), max(cast(a.sent_at as char)) "+
			"from answer a "+
			"inner join (select t.id, t.description, t.deadline, t.priority, "+
			"ts.status, ts.status_level level from task t "+
			"left join task_status ts on ts.id = t.status "+
			"where t.project = ? and t.is_closed <> true and t.is_closed = 0) t on t.id = a.task "+
			"where a.status = ? "+
			"group by a.task "+
			"order by max(a.id) asc",
		projectId,
		"not checked",
	)
	if err != nil {
		log.Println("error in getting answers: ", err)
		return 500, err.Error()
	}
	answers := make([]*model.Answer, 0)
	for rows.Next() {
		answer := &model.Answer{
			Id: 0,
			Task: &model.Task{
				Id:          0,
				Description: "",
				Project:     nil,
				Deadline:    "",
				Priority:    "",
				Status: &model.TaskStatus{
					Status: "",
					Level:  0,
				},
				IsClosed: false,
				Executor: nil,
				Skills:   nil,
			},
			Executor:  nil,
			MessageId: 0,
			ChatId:    0,
			Status:    "",
			SentAt:    time.Time{},
		}
		var date string

		err = rows.Scan(
			&answer.Id,
			&answer.Task.Id,
			&answer.Task.Description,
			&answer.Task.Deadline,
			&answer.Task.Priority,
			&answer.Task.Status.Status,
			&answer.Task.Status.Level,
			&answer.MessageId,
			&answer.ChatId,
			&answer.Status,
			&date,
		)
		if err != nil {
			log.Println("error in scanning answers: ", err)
			continue
		}

		answer.SentAt, err = time.Parse("2006-01-02 15:04:05", date)
		if err != nil {
			log.Println("error in casting date: ", err)
			continue
		}

		answers = append(answers, answer)
	}

	w.Header().Set("Content-Type", "application/json")

	return c.makeContentResponse(200, "Project answers", answers)
}

func (c *Controller) AcceptAnswer(params martini.Params, w http.ResponseWriter) (int, string) {
	answerId, err := strconv.ParseInt(params["aid"], 10, 64)
	if err != nil {
		log.Println("error in parsing answer id: ", err)
		return 500, err.Error()
	}

	result, err := c.DB.Exec(
		"update answer set status = ? where id = ?",
		"accepted",
		answerId,
	)
	if err != nil {
		log.Println("error in accepting answer: ", err.Error())
		return 500, err.Error()
	}
	rowsAffected, _ := result.RowsAffected()

	w.Header().Set("Content-Type", "application/json")
	return c.makeContentResponse(200, "Answer accepted", struct {
		Name    string
		Content interface{}
	}{
		Name:    "Rows affected",
		Content: rowsAffected,
	})
}

func (c *Controller) DeclineAnswer(params martini.Params, w http.ResponseWriter) (int, string) {
	answerId, err := strconv.ParseInt(params["aid"], 10, 64)
	if err != nil {
		log.Println("error in parsing answer id: ", err)
		return 500, err.Error()
	}

	result, err := c.DB.Exec(
		"update answer set status = ? where id = ?",
		"declined",
		answerId,
	)
	if err != nil {
		log.Println("error in declining answer: ", err.Error())
		return 500, err.Error()
	}
	rowsAffected, _ := result.RowsAffected()

	w.Header().Set("Content-Type", "application/json")
	return c.makeContentResponse(200, "Answer declined", struct {
		Name    string
		Content interface{}
	}{
		Name:    "Rows affected",
		Content: rowsAffected,
	})
}
