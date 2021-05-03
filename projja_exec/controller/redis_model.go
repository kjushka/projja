package controller

import "projja-exec/model"

type userSkillsData struct {
	UserId int64
	Skills []string
}

type addingMemberData struct {
	ProjectId int64
	Member    *model.User
}

type removingMemberData struct {
	ProjectId      int64
	MemberUsername string
}

type newTaskData struct {
	ProjectId int64
	Task      *model.Task
}
