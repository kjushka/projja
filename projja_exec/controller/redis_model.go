package controller

import (
	"projja_exec/model"
	"time"
)

type userSkillsData struct {
	UserId      int64
	Skills      []string
	ProjectsIds []int64
}

type updateUserData struct {
	NewUserInfo *model.User
	ProjectsIds []int64
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

type changeExecutorData struct {
	TaskId    int64
	OldUserId int64
	NewUserId int64
	ProjectId int64
}

type changeDescriptionData struct {
	TaskId      int64
	Description string
	ProjectId   int64
}

type closeTaskData struct {
	TaskId     int64
	ExecutorId int64
	ProjectId  int64
}

type changeDeadlineData struct {
	TaskId    int64
	Deadline  time.Time
	ProjectId int64
}
