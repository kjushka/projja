package graph

import (
	"projja-exec/model"
	"time"
)

func (g *Graph) SetSkillsToUser(userId int64, skills []string) {
	executor := g.Users[userId]
	newSkills := make([]string, 0)
	for _, s := range skills {
		isNotExist := true

		for _, gs := range g.Skills {
			if s == gs {
				isNotExist = false
				break
			}
		}

		if isNotExist {
			newSkills = append(newSkills, s)
		}
	}

	g.Skills = append(g.Skills, newSkills...)
	executor.Skills = skills
	g.Users[userId] = executor
}

func (g *Graph) UpdateUserInfo(newUserInfo *model.User) {
	g.Users[newUserInfo.Id] = newUserInfo
}

func (g *Graph) AddExecutor(executor *model.User) {
	g.Users[executor.Id] = executor

	newSkills := make([]string, 0)
	for _, s := range executor.Skills {
		isNotExist := true

		for _, gs := range g.Skills {
			if s == gs {
				isNotExist = false
				break
			}
		}

		if isNotExist {
			newSkills = append(newSkills, s)
		}
	}

	g.Skills = append(g.Skills, newSkills...)
}

func (g *Graph) RemoveMember(memberUsername string) {
	for userId, user := range g.Users {
		if user.Username == memberUsername {
			delete(g.Users, userId)
			tasksIds := g.UserToTask[userId]

			for _, taskId := range tasksIds {
				delete(g.Tasks, taskId)
				delete(g.TaskToSkill, taskId)
			}

			delete(g.UserToTask, userId)
			break
		}
	}
}

func (g *Graph) AddTaskWithExecutor(task *model.Task) {
	g.Tasks[task.Id] = task
	g.UserToTask[task.Executor.Id] = append(g.UserToTask[task.Executor.Id], task.Id)
}

func (g *Graph) ChangeTaskExecutor(oldUserId int64, newUserId int64, taskId int64) {
	for index, task := range g.UserToTask[oldUserId] {
		if task == taskId {
			g.UserToTask[oldUserId][index] = g.UserToTask[oldUserId][len(g.UserToTask[oldUserId])-1]
			g.UserToTask[oldUserId][len(g.UserToTask[oldUserId])-1] = 0
			g.UserToTask[oldUserId] = g.UserToTask[oldUserId][:len(g.UserToTask[oldUserId])-1]
			break
		}
	}
	g.UserToTask[newUserId] = append(g.UserToTask[newUserId], taskId)
}

func (g *Graph) ChangeTaskDescription(taskId int64, description string) {
	g.Tasks[taskId].Description = description
}

func (g *Graph) CloseTask(taskId int64, executorId int64) {
	delete(g.Tasks, taskId)

	for index, task := range g.UserToTask[executorId] {
		if task == taskId {
			g.UserToTask[executorId][index] = g.UserToTask[executorId][len(g.UserToTask[executorId])-1]
			g.UserToTask[executorId][len(g.UserToTask[executorId])-1] = 0
			g.UserToTask[executorId] = g.UserToTask[executorId][:len(g.UserToTask[executorId])-1]
			break
		}
	}
}

func (g *Graph) ChangeTaskDeadline(taskId int64, deadline time.Time) {
	g.Tasks[taskId].Deadline = deadline
}
