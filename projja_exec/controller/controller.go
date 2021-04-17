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
	"sync"

	"github.com/go-martini/martini"
	"github.com/go-redis/redis/v8"
)

type controller struct {
	Rds      *redis.Client
	Projects map[int64]*usingProject
	Mutex    *sync.Mutex
}

type usingProject struct {
	UsingCount int
	Project    *graph.Project
}

func NewController(options *redis.Options) *controller {
	return &controller{
		Rds:      redis.NewClient(options),
		Projects: make(map[int64]*usingProject),
		Mutex:    &sync.Mutex{},
	}
}

func (c *controller) CheckContentType(w http.ResponseWriter, r *http.Request) {
	contentType := r.Header.Get("Content-Type")
	if r.Method == "POST" && contentType != "application/json" {
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
	id := params["id"]
	intId, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		log.Println(err)
		return 500, err.Error()
	}

	jsonExec, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Println("error during reading body:", err)
		return 500, err.Error()
	}
	exec := &model.User{}
	err = json.Unmarshal(jsonExec, exec)
	if err != nil {
		log.Println("error during unmarshalling:", err)
		return 500, err.Error()
	}

	var project *graph.Project
	c.Mutex.Lock()
	if using, ok := c.Projects[intId]; ok {
		project = using.Project
		using.UsingCount++
	} else {
		project, err = c.readData(intId)
		if err != nil {
			log.Println(err)
			return 500, err.Error()
		}

		c.Projects[intId] = &usingProject{1, project}
	}
	c.Mutex.Unlock()

	project.Graph.AddExecutor(exec)

	c.Mutex.Lock()
	if using, ok := c.Projects[intId]; ok && using.UsingCount == 1 {
		delete(c.Projects, intId)
	} else if ok && using.UsingCount > 1 {
		using.Project = project
		using.UsingCount--
	}
	c.Mutex.Unlock()

	err = c.writeProjectToRedis(project)
	if err != nil {
		log.Println(err)
		return 500, err.Error()
	}

	return 200, "Executor added"
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
	project, err := c.readData(intId)
	if err != nil {
		log.Println(err.Error())
		return 500, err.Error()
	}
	return c.makeContentResponse(200, "project", project)
}
