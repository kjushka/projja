package graph

import (
	"projja-exec/model"
)

type Project struct {
	Id    int64
	Graph *Graph
}

type Graph struct {
	Users       []*model.User
	Tasks       []*model.Task
	Skills      []string
	UserToTask  map[int][]float32
	UserToSkill map[int][]string
	TaskToSkill map[int][]string
}

func MakeNewProject(newProject *model.Project) *Project {
	project := &Project{
		Id: newProject.Id,
		Graph: &Graph{
			Users:       make([]*model.User, 1),
			Tasks:       make([]*model.Task, 0),
			Skills:      make([]string, 0),
			UserToTask:  make(map[int][]float32),
			UserToSkill: make(map[int][]string, 1),
			TaskToSkill: make(map[int][]string, 1),
		},
	}
	project.Graph.Users[0] = newProject.Owner
	project.Graph.Skills = append(project.Graph.Skills, newProject.Owner.Skills...)
	project.Graph.UserToSkill[0] = newProject.Owner.Skills

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

	g.UserToSkill[len(g.Users)-1] = executor.Skills
}

func (g *Graph) AddTask(*model.Task) {

}

func (g *Graph) caluculateExecutor() {

}
