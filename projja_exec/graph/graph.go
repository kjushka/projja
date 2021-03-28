package graph

import "projja-exec/model"

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

func (p *Project) LoadGraph() {
	graph := &Graph{
		Users:       make([]*model.User, 0),
		Tasks:       make([]*model.Task, 0),
		Skills:      make([]string, 0),
		UserToTask:  make(map[int][]float32),
		UserToSkill: make(map[int][]string),
		TaskToSkill: make(map[int][]string),
	}
	p.Graph = graph
}

func (g *Graph) AddTask(*model.Task) {

}

func (g *Graph) caluculateExecutor() {

}
