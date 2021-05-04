package controller

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"projja-exec/graph"
	"projja-exec/model"
	"strconv"

	"github.com/go-redis/redis/v8"
)

func (c *controller) ListenExecStream(ctx context.Context) {
	for {
		xStreamSlice := c.Rds.XRead(ctx, &redis.XReadArgs{Block: 0, Streams: []string{"exec", "$"}})
		xReadResult, err := xStreamSlice.Result()
		if err != nil {
			if err == redis.Nil {
				continue
			} else {
				panic(err)
			}
		}
		rdsMessages := xReadResult[len(xReadResult)-1].Messages
		for _, rdsMessage := range rdsMessages {
			rdsMap := rdsMessage.Values
			if val, ok := rdsMap["skills"]; ok {
				go c.setSkillsToUser(val)
			}
			if val, ok := rdsMap["info"]; ok {
				go c.updateUserInfo(val)
			}
		}
	}
}

func (c *controller) setSkillsToUser(jsonSkills interface{}) {
	if strJsonSkills, ok := jsonSkills.(string); ok {
		skillsData := &userSkillsData{}
		err := json.Unmarshal([]byte(strJsonSkills), skillsData)
		if err != nil {
			log.Println("error in unmarshalling:", err)
			return
		}
		err = c.setSkillsToUserInGraph(skillsData)
		if err != nil {
			log.Println("error in setting skills to user: ", err)
		}
	} else {
		log.Println("error in casting user skills id")
	}
}

func (c *controller) updateUserInfo(jsonUserInfo interface{}) {
	if strJsonUserData, ok := jsonUserInfo.(string); ok {
		userData := &updateUserData{}
		err := json.Unmarshal([]byte(strJsonUserData), userData)
		if err != nil {
			log.Println("error in unmarshalling:", err)
			return
		}
		err = c.updateUserInfoInGraph(userData)
		if err != nil {
			log.Println("error in updating user info: ", err)
		}
	} else {
		log.Println("error in casting user info")
	}
}

func (c *controller) ListenProjectStream(ctx context.Context) {
	for {
		xStreamSlice := c.Rds.XRead(ctx, &redis.XReadArgs{Block: 0, Streams: []string{"project", "$"}})
		xReadResult, err := xStreamSlice.Result()
		if err != nil {
			if err == redis.Nil {
				continue
			} else {
				panic(err)
			}
		}
		rdsMessages := xReadResult[len(xReadResult)-1].Messages
		for _, rdsMessage := range rdsMessages {
			rdsMap := rdsMessage.Values
			if val, ok := rdsMap["new"]; ok {
				go c.createNewProject(val)
			}
			if val, ok := rdsMap["add-member"]; ok {
				go c.addMember(val)
			}
			if val, ok := rdsMap["remove-member"]; ok {
				go c.removeMember(val)
			}
			if val, ok := rdsMap["task"]; ok {
				go c.createTask(val)
			}
		}
	}
}

func (c *controller) createNewProject(jsonProject interface{}) {
	if strJsonProject, ok := jsonProject.(string); ok {
		newProject := &model.Project{}
		err := json.Unmarshal([]byte(strJsonProject), newProject)
		if err != nil {
			log.Println("error in unmarshalling:", err)
			return
		}
		err = c.addProject(newProject)
		if err != nil {
			log.Println("error in creating project: ", err)
		}
	} else {
		log.Println("error in casting project")
	}
}

func (c *controller) addMember(jsonProjectNewMember interface{}) {
	if strJsonProjectNewMember, ok := jsonProjectNewMember.(string); ok {
		newProjectMemberData := &addingMemberData{}
		err := json.Unmarshal([]byte(strJsonProjectNewMember), newProjectMemberData)
		if err != nil {
			log.Println("error in unmarshalling:", err)
			return
		}
		err = c.addMemberInGraph(newProjectMemberData)
		if err != nil {
			log.Println("error in adding member: ", err)
		}
	} else {
		log.Println("error in casting project member data")
	}
}

func (c *controller) removeMember(jsonRemovingMember interface{}) {
	if strJsonRemovingMember, ok := jsonRemovingMember.(string); ok {
		removingMember := &removingMemberData{}
		err := json.Unmarshal([]byte(strJsonRemovingMember), removingMember)
		if err != nil {
			log.Println("error in unmarshalling:", err)
			return
		}
		err = c.removeMemberInGraph(removingMember)
		if err != nil {
			log.Println("error in removing member: ", err)
		}
	} else {
		log.Println("error in casting removing member data")
	}
}

func (c *controller) createTask(jsonTask interface{}) {
	if strJsonTask, ok := jsonTask.(string); ok {
		task := &newTaskData{}
		err := json.Unmarshal([]byte(strJsonTask), task)
		if err != nil {
			log.Println("error in unmarshalling:", err)
			return
		}
		err = c.createTaskInGraph(task)
		if err != nil {
			log.Println("error in creating task: ", err)
		}
	} else {
		log.Println("error in casting task")
	}
}

func (c *controller) ListenTaskStream(ctx context.Context) {
	for {
		xStreamSlice := c.Rds.XRead(ctx, &redis.XReadArgs{Block: 0, Streams: []string{"task", "$"}})
		xReadResult, err := xStreamSlice.Result()
		if err != nil {
			if err == redis.Nil {
				continue
			} else {
				panic(err)
			}
		}
		rdsMessages := xReadResult[len(xReadResult)-1].Messages
		for _, rdsMessage := range rdsMessages {
			rdsMap := rdsMessage.Values
			if val, ok := rdsMap["executor"]; ok {
				go c.changeTaskExecutor(val)
			}
			if val, ok := rdsMap["description"]; ok {
				go c.changeTaskDescription(val)
			}
			if val, ok := rdsMap["close"]; ok {
				go c.closeTask(val)
			}
			if val, ok := rdsMap["deadline"]; ok {
				go c.changeTaskDeadline(val)
			}
		}
	}
}

func (c *controller) changeTaskExecutor(jsonChangeExecutorData interface{}) {
	if strJsonChangeExecutorData, ok := jsonChangeExecutorData.(string); ok {
		changeExecutor := &changeExecutorData{}
		err := json.Unmarshal([]byte(strJsonChangeExecutorData), changeExecutor)
		if err != nil {
			log.Println("error in unmarshalling:", err)
			return
		}
		err = c.changeTaskExecutorInGraph(changeExecutor)
		if err != nil {
			log.Println("error in changing task executor: ", err)
		}
	} else {
		log.Println("error in casting change executor data")
	}
}

func (c *controller) changeTaskDescription(jsonChangeDescriptionData interface{}) {
	if strJsonChangeDescriptionData, ok := jsonChangeDescriptionData.(string); ok {
		changeDescription := &changeDescriptionData{}
		err := json.Unmarshal([]byte(strJsonChangeDescriptionData), changeDescription)
		if err != nil {
			log.Println("error in unmarshalling:", err)
			return
		}
		err = c.changeTaskDescriptionInGraph(changeDescription)
		if err != nil {
			log.Println("error in change description: ", err)
		}
	} else {
		log.Println("error in casting change description")
	}
}

func (c *controller) closeTask(jsonCloseData interface{}) {
	if strJsonCloseData, ok := jsonCloseData.(string); ok {
		closeData := &closeTaskData{}
		err := json.Unmarshal([]byte(strJsonCloseData), closeData)
		if err != nil {
			log.Println("error in unmarshalling:", err)
			return
		}
		err = c.closeTaskInGraph(closeData)
		if err != nil {
			log.Println("error in closing task: ", err)
		}
	} else {
		log.Println("error in casting close task")
	}
}

func (c *controller) changeTaskDeadline(jsonChangeDeadlineData interface{}) {
	if strJsonChangeDeadlineData, ok := jsonChangeDeadlineData.(string); ok {
		changeDeadline := &changeDeadlineData{}
		err := json.Unmarshal([]byte(strJsonChangeDeadlineData), changeDeadline)
		if err != nil {
			log.Println("error in unmarshalling:", err)
			return
		}
		err = c.changeTaskDeadlineInGraph(changeDeadline)
		if err != nil {
			log.Println("error in change deadline: ", err)
		}
	} else {
		log.Println("error in casting change deadline")
	}
}

func (c *controller) saveNewProject(newProject *model.Project) error {
	project := graph.MakeNewProject(newProject)

	err := c.writeProjectToRedis(project)
	if err != nil {
		return err
	}
	return nil
}

func (c *controller) writeProjectToRedis(project *graph.Project) error {
	ctx := context.Background()

	byteGraph, err := json.Marshal(project.Graph)
	if err != nil {
		return err
	}

	status := c.Rds.Set(ctx, strconv.FormatInt(project.Id, 10), string(byteGraph), 0)
	if status.Err() != nil {
		return status.Err()
	}
	return nil
}

func (c *controller) readData(id int64) (*graph.Project, error) {
	val, err := c.Rds.Get(context.Background(), strconv.FormatInt(id, 10)).Result()
	switch {
	case err == redis.Nil:
		log.Println("key does not exist")
		return nil, err
	case err != nil:
		log.Println("Get failed", err)
		return nil, err
	case val == "":
		err = errors.New("value is empty")
		log.Println(err)
		return nil, err
	}

	log.Println(val)
	g := &graph.Graph{}
	err = json.Unmarshal([]byte(val), g)
	if err != nil {
		return nil, err
	}
	return &graph.Project{
		Id:    id,
		Graph: g,
	}, nil
}
