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
	var err error = nil
	for _, projectId := range skillsData.ProjectsIds {
		err = c.setSkillsToUserInProject(projectId, skillsData.UserId, skillsData.Skills)
		if err != nil {
			return err
		}
	}

	return err
}

func (c *controller) setSkillsToUserInProject(projectId int64, userId int64, skills []string) error {
	project, err := c.getProject(projectId)
	if err != nil {
		return err
	}

	project.Graph.SetSkillsToUser(userId, skills)

	err = c.closeProjectWork(project, projectId)

	return err
}

func (c *controller) updateUserInfoInGraph(userInfo *updateUserData) error {
	var err error = nil
	for _, projectId := range userInfo.ProjectsIds {
		err = c.changeUserData(projectId, userInfo.NewUserInfo)
		if err != nil {
			return err
		}
	}

	return err
}

func (c *controller) changeUserData(projectId int64, newUserInfo *model.User) error {
	project, err := c.getProject(projectId)
	if err != nil {
		return err
	}

	project.Graph.ChangeUserData(newUserInfo)

	err = c.closeProjectWork(project, projectId)

	return err
}

func (c *controller) addProject(newProject *model.Project) error {
	err := c.saveNewProject(newProject)
	return err
}

func (c *controller) addMemberInGraph(newMemberData *addingMemberData) error {
	project, err := c.getProject(newMemberData.ProjectId)
	if err != nil {
		return err
	}

	project.Graph.AddExecutor(newMemberData.Member)

	err = c.closeProjectWork(project, newMemberData.ProjectId)

	return err
}

func (c *controller) removeMemberInGraph(removingMember *removingMemberData) error {
	project, err := c.getProject(removingMember.ProjectId)
	if err != nil {
		return err
	}

	project.Graph.RemoveMember(removingMember.MemberUsername)

	err = c.closeProjectWork(project, removingMember.ProjectId)

	return err
}

func (c *controller) createTaskInGraph(taskData *newTaskData) error {
	project, err := c.getProject(taskData.ProjectId)
	if err != nil {
		return err
	}

	project.Graph.AddTaskWithExecutor(taskData.Task)

	err = c.closeProjectWork(project, taskData.ProjectId)

	return err
}

func (c *controller) getProject(projectId int64) (*graph.Project, error) {
	var project *graph.Project
	var err error
	c.Mutex.Lock()
	if using, ok := c.Projects[projectId]; ok {
		project = using.Project
		using.UsingCount++
	} else {
		project, err = c.readData(projectId)
		if err != nil {
			log.Println(err)
			defer c.Mutex.Unlock()
			return nil, err
		}

		c.Projects[projectId] = &usingProject{1, project}
	}
	c.Mutex.Unlock()

	return project, err
}

func (c *controller) closeProjectWork(project *graph.Project, projectId int64) error {
	c.Mutex.Lock()
	if using, ok := c.Projects[projectId]; ok && using.UsingCount == 1 {
		delete(c.Projects, projectId)
	} else if ok && using.UsingCount > 1 {
		using.Project = project
		using.UsingCount--
	}
	c.Mutex.Unlock()

	err := c.writeProjectToRedis(project)
	return err
}
