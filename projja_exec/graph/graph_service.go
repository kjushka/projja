package graph

import "projja-exec/model"

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

func (g *Graph) ChangeUserData(newUserInfo *model.User) {
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
