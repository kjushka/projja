package graph

import (
	"log"
	"math"
	"projja-exec/model"
	"sort"
	"time"
)

const (
	timeCoefficient float32 = 1
	skillsCoefficient float32 = 1
)

type Project struct {
	Id    int64
	Graph *Graph
}

type Graph struct {
	Users        []*model.User
	Tasks        []*model.Task
	Skills       []string
	UserToTask   map[int][]int
	UserToSkills map[int][]string
	TaskToSkill  map[int][]string
}

type rating struct {
	UsersRating map[int]*userRating
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
			Users:        make([]*model.User, 1),
			Tasks:        make([]*model.Task, 0),
			Skills:       make([]string, 0),
			UserToTask:   make(map[int][]int),
			UserToSkills: make(map[int][]string, 1),
			TaskToSkill:  make(map[int][]string, 1),
		},
	}
	project.Graph.Users[0] = newProject.Owner
	project.Graph.Skills = append(project.Graph.Skills, newProject.Owner.Skills...)
	project.Graph.UserToSkills[0] = newProject.Owner.Skills

	return project
}

func (g *Graph) AddExecutor(executor *model.User) {
	g.Users = append(g.Users, executor)

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

	g.UserToSkills[len(g.Users)-1] = executor.Skills
}

func (g *Graph) AddTaskWithExecutor(userIndex int, task *model.Task) {
	taskIndex := len(g.Tasks)
	g.Tasks = append(g.Tasks, task)
	g.UserToTask[userIndex] = append(g.UserToTask[userIndex], taskIndex)
}

func (g *Graph) CalculateNewTaskExecutor(task *model.Task) (*model.User, int) {
	ratio := g.calculateRatingBySkills(task.Skills)
	g.calculateRatingByTime(task.Deadline, ratio)
	executor, index := g.selectExecutorByRating(ratio)
	go g.checkCorrectWork()
	
	return executor, index
}

func (g *Graph) calculateRatingBySkills(taskSkills []string) *rating {
	ratingStruct := &rating{UsersRating: make(map[int]*userRating, len(g.Users))}

	for i, user := range g.Users {
		ratio := g.checkSkillsSimilarity(g.UserToSkills[i], taskSkills)
		userRate := &userRating{
			User:         user,
			SkillsRating: ratio,
			TimeRating:   0,
		}
		ratingStruct.UsersRating[i] = userRate
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
	for i, _ := range g.Users {
		rate := g.checkTime(i, deadline)
		ratio.UsersRating[i].TimeRating = rate
	}
}

func (g *Graph) checkTime(userIndex int, deadline time.Time) float32 {
	tasksDeadlines := make([]int, len(g.UserToTask[userIndex]))
	for i, index := range g.UserToTask[userIndex] {
		daysTo := int(math.Ceil(time.Until(g.Tasks[index].Deadline).Hours()))
		if daysTo > 0 {
			tasksDeadlines[i] = daysTo
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

func (g *Graph) selectExecutorByRating(ratio *rating) (*model.User, int) {
	var bestUser *model.User = nil
	bestUserIndex := 0
	var bestRating float32 = -1
	for index, userRatio := range ratio.UsersRating {
		userRate := userRatio.TimeRating * timeCoefficient + userRatio.SkillsRating * skillsCoefficient
		if userRate > bestRating {
			bestUser = userRatio.User
			bestRating = userRate
			bestUserIndex = index
		}
	}

	return bestUser, bestUserIndex
}

func (g *Graph) checkCorrectWork() {
	log.Println("i check the correction of algorithm")
}
