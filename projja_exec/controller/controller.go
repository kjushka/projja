package controller

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"projja-exec/graph"
	"projja-exec/model"
	"strconv"

	"github.com/go-martini/martini"
	"github.com/go-redis/redis/v8"
)

type controller struct {
	Rds      *redis.Client
	Projects map[int64]*usingProject
}

type usingProject struct {
	UsingCount int
	Project    *graph.Project
}

func NewController(options *redis.Options) *controller {
	return &controller{
		Rds:      redis.NewClient(options),
		Projects: make(map[int64]*usingProject),
	}
}

func (c *controller) CheckContentType(w http.ResponseWriter, r *http.Request) {
	contentType := r.Header.Get("Content-Type")
	if contentType != "application/json" {
		err := fmt.Sprintf("Unsupportable Content-Type header: %s", contentType)
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err))
		return
	}
}

func (c *controller) AddProject(w http.ResponseWriter, r *http.Request) (int, string) {
	jsonProject, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Println("error in reading body: ", err)
		return 500, err.Error()
	}
	defer r.Body.Close()

	newProject := &model.Project{}
	err = json.Unmarshal(jsonProject, newProject)
	if err != nil {
		log.Println("error in unmarshalling new project: ", err)
		return 500, err.Error()
	}

	if newProject.Status == "closed" {
		err := "saving closed project declined"
		log.Println(err)
		return 500, err
	}

	err = c.saveNewProject(newProject)
	if err != nil {
		log.Println("error in saving new project: ", err)
		return 500, err.Error()
	}

	return c.makeContentResponse(200, "Project saved", newProject)
}

func (c *controller) AddExecutorToProject(params martini.Params, w http.ResponseWriter, r *http.Request) (int, string) {

	return 200, ""
}

func (c *controller) AddTaskToProject(params martini.Params, w http.ResponseWriter, r *http.Request) (int, string) {

	return 200, ""
}

func (c *controller) GetRedisData(params martini.Params, w http.ResponseWriter) (int, string) {
	id := params["id"]
	intId, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		log.Println(err)
		return 500, err.Error()
	}
	project, err := c.ReadData(intId)
	if err != nil {
		log.Println(err.Error())
		return 500, err.Error()
	}
	return c.makeContentResponse(200, "project", project)
}
