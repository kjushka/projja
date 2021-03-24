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
	userToTask  [][]float32
	userToSkill [][]bool
	taskToSkill [][]bool
}

func (g *Graph) AddTask(*model.Task) {

}

func (g *Graph) CaluculateExecutors() {
}

func (g *Graph) calculateExecutorsBySkill() []*model.User {
	return nil
}

func (g *Graph) caluculateExecutorsByTime() []*model.User {
	return nil
}

func (g *Graph) GetExecutorsToTasks() {}

func (p *Project) LoadGraph() {
	graph := &Graph{
		Users:       make([]*model.User, 0),
		Tasks:       make([]*model.Task, 0),
		Skills:      make([]string, 0),
		userToTask:  make([][]float32, 0),
		userToSkill: make([][]bool, 0),
		taskToSkill: make([][]bool, 0),
	}
	p.Graph = graph
}
