package graph

import (
	"log"
	"math"
	"projja-exec/model"
	"sort"
	"time"
)

const (
	timeCoefficient   float32 = 1
	skillsCoefficient float32 = 1
)

type Project struct {
	Id    int64
	Graph *Graph
}

type Graph struct {
	Users       map[int64]*model.User
	Tasks       map[int64]*model.Task
	Skills      []string
	UserToTask  map[int64][]int64
	TaskToSkill map[int64][]string
}

type rating struct {
	UsersRating map[int64]*userRating
}

type userRating struct {
	User         *model.User
	TimeRating   float32
	SkillsRating float32
}

func MakeNewProject(newProject *model.Project) *Project {
	project := &Project{
		Id: newProject.Id,
		Graph: &Graph{
			Users:       make(map[int64]*model.User, 1),
			Tasks:       make(map[int64]*model.Task, 0),
			Skills:      make([]string, 0),
			UserToTask:  make(map[int64][]int64),
			TaskToSkill: make(map[int64][]string, 1),
		},
	}
	project.Graph.Users[newProject.Owner.Id] = newProject.Owner
	project.Graph.Skills = append(project.Graph.Skills, newProject.Owner.Skills...)

	return project
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

func (g *Graph) CalculateNewTaskExecutor(task *model.Task) *model.User {
	ratio := g.calculateRatingBySkills(task.Skills)
	g.calculateRatingByTime(task.Deadline, ratio)
	executor := g.selectExecutorByRating(ratio)
	go g.checkCorrectWork()

	return executor
}

func (g *Graph) calculateRatingBySkills(taskSkills []string) *rating {
	ratingStruct := &rating{UsersRating: make(map[int64]*userRating, len(g.Users))}

	for userId, user := range g.Users {
		ratio := g.checkSkillsSimilarity(user.Skills, taskSkills)
		userRate := &userRating{
			User:         user,
			SkillsRating: ratio,
			TimeRating:   0,
		}
		ratingStruct.UsersRating[userId] = userRate
	}

	return ratingStruct
}

func (g *Graph) checkSkillsSimilarity(userSkills []string, taskSkills []string) float32 {
	count := len(taskSkills)
	var ratio float32 = 0
	for _, skill := range userSkills {
		contains := false
		for _, taskSkill := range taskSkills {
			if skill == taskSkill {
				contains = true
				break
			}
		}
		if contains {
			ratio += float32(1 / count)
		}
	}

	return ratio
}

func (g *Graph) calculateRatingByTime(deadline time.Time, ratio *rating) {
	for userId, _ := range g.Users {
		rate := g.checkTime(userId, deadline)
		ratio.UsersRating[userId].TimeRating = rate
	}
}

func (g *Graph) checkTime(userId int64, deadline time.Time) float32 {
	tasksDeadlines := make([]int, 0)
	for _, taskId := range g.UserToTask[userId] {
		daysTo := int(math.Ceil(time.Until(g.Tasks[taskId].Deadline).Hours()))
		if daysTo > 0 {
			tasksDeadlines = append(tasksDeadlines, daysTo)
		}
	}

	sort.Ints(tasksDeadlines)
	ratio := float32(0)
	daysToTaskDeadline := int(math.Ceil(time.Until(deadline).Hours()))
	prev := 0
	count := len(tasksDeadlines) + 1

	for _, days := range tasksDeadlines {
		if days > daysToTaskDeadline {
			break
		}
		intervalDateRate := float32((days - prev) / daysToTaskDeadline)
		intervalTaskRate := float32(1 / count)
		ratio += intervalTaskRate * intervalDateRate

		prev = days
		count--
	}

	return ratio
}

func (g *Graph) selectExecutorByRating(ratio *rating) *model.User {
	var bestUser *model.User = nil
	var bestRating float32 = -1
	for _, userRatio := range ratio.UsersRating {
		userRate := userRatio.TimeRating*timeCoefficient + userRatio.SkillsRating*skillsCoefficient
		if userRate > bestRating {
			bestUser = userRatio.User
			bestRating = userRate
		}
	}

	return bestUser
}

func (g *Graph) checkCorrectWork() {
	log.Println("i check the correction of algorithm")
}
