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

func (c *controller) CalculateTaskExecutor(params martini.Params, w http.ResponseWriter, r *http.Request) (int, string) {
	id := params["id"]
	intId, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		log.Println(err)
		return 500, err.Error()
	}

	jsonTask, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Println("error in printing body: ", err)
		return 500, err.Error()
	}
	defer r.Body.Close()

	task := &model.Task{}
	err = json.Unmarshal(jsonTask, task)
	if err != nil {
		log.Println("error in unmarshalling: ", err)
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
			defer c.Mutex.Unlock()
			return 500, err.Error()
		}

		c.Projects[intId] = &usingProject{1, project}
	}
	c.Mutex.Unlock()

	executor := project.Graph.CalculateNewTaskExecutor(task)
	task.Executor = executor

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

	w.Header().Set("Content-Type", "application/json")
	return c.makeContentResponse(200, "Task with executor", task)
}

func (c *controller) GetRedisData(params martini.Params) (int, string) {
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

func (c *controller) setSkillsToUserInGraph(skillsData *userSkillsData) error {
	return nil
}

func (c *controller) updateUserInfoInGraph(userInfo *model.User) error {
	return nil
}

func (c *controller) addProject(newProject *model.Project) error {
	err := c.saveNewProject(newProject)
	return err
}

func (c *controller) addMemberInGraph(newMemberData *addingMemberData) error {
	var project *graph.Project
	var err error
	c.Mutex.Lock()
	if using, ok := c.Projects[newMemberData.ProjectId]; ok {
		project = using.Project
		using.UsingCount++
	} else {
		project, err = c.readData(newMemberData.ProjectId)
		if err != nil {
			log.Println(err)
			defer c.Mutex.Unlock()
			return err
		}

		c.Projects[newMemberData.ProjectId] = &usingProject{1, project}
	}
	c.Mutex.Unlock()

	project.Graph.AddExecutor(newMemberData.Member)

	c.Mutex.Lock()
	if using, ok := c.Projects[newMemberData.ProjectId]; ok && using.UsingCount == 1 {
		delete(c.Projects, newMemberData.ProjectId)
	} else if ok && using.UsingCount > 1 {
		using.Project = project
		using.UsingCount--
	}
	c.Mutex.Unlock()

	err = c.writeProjectToRedis(project)
	return err
}

func (c *controller) removeMemberInGraph(removingMember *removingMemberData) error {
	var project *graph.Project
	var err error
	c.Mutex.Lock()
	if using, ok := c.Projects[removingMember.ProjectId]; ok {
		project = using.Project
		using.UsingCount++
	} else {
		project, err = c.readData(removingMember.ProjectId)
		if err != nil {
			log.Println(err)
			defer c.Mutex.Unlock()
			return err
		}

		c.Projects[removingMember.ProjectId] = &usingProject{1, project}
	}
	c.Mutex.Unlock()

	project.Graph.RemoveMember(removingMember.MemberUsername)

	c.Mutex.Lock()
	if using, ok := c.Projects[removingMember.ProjectId]; ok && using.UsingCount == 1 {
		delete(c.Projects, removingMember.ProjectId)
	} else if ok && using.UsingCount > 1 {
		using.Project = project
		using.UsingCount--
	}
	c.Mutex.Unlock()

	err = c.writeProjectToRedis(project)
	return err
}

func (c *controller) createTaskInGraph(taskData *newTaskData) error {
	var project *graph.Project
	var err error
	c.Mutex.Lock()
	if using, ok := c.Projects[taskData.ProjectId]; ok {
		project = using.Project
		using.UsingCount++
	} else {
		project, err = c.readData(taskData.ProjectId)
		if err != nil {
			log.Println(err)
			defer c.Mutex.Unlock()
			return err
		}

		c.Projects[taskData.ProjectId] = &usingProject{1, project}
	}
	c.Mutex.Unlock()

	project.Graph.AddTaskWithExecutor(taskData.Task)

	c.Mutex.Lock()
	if using, ok := c.Projects[taskData.ProjectId]; ok && using.UsingCount == 1 {
		delete(c.Projects, taskData.ProjectId)
	} else if ok && using.UsingCount > 1 {
		using.Project = project
		using.UsingCount--
	}
	c.Mutex.Unlock()

	err = c.writeProjectToRedis(project)
	return err
}
