package task

import (
	"fmt"
	"todo/entity"
)

type ServiceRepository interface {
	// DoesThisUesrHaveThisCategoryID(userID, categoryID uint) bool
	CreateNewTask(t entity.Task) (entity.Task, error)
	ListUserTasks(userID uint) ([]entity.Task, error)
}

type Service struct {
	repository ServiceRepository
}

func NewService(repo ServiceRepository) Service {
	return Service{
		repository: repo,
	}
}

type CreateRequest struct {
	Title               string
	DueDate             string
	CategoryID          uint
	AuthenticatedUserID uint
}

type CreateResponse struct {
	Task entity.Task
}

func (t Service) Create(req CreateRequest) (CreateResponse, error) {
	// if !t.repository.DoesThisUesrHaveThisCategoryID(req.AuthenticatedUserID, req.CategoryID) {
	// 	return CreateResponse{}, fmt.Errorf("user does not have this category: %d\n", req.CategoryID)
	// }

	createdTask, cErr := t.repository.CreateNewTask(entity.Task{
		Title:      req.Title,
		DueDate:    req.DueDate,
		CategoryID: req.CategoryID,
		IsDone:     false,
		UserID:     req.AuthenticatedUserID,
	})

	if cErr != nil {
		return CreateResponse{}, fmt.Errorf("can't create new task: %v\n", cErr)
	}

	return CreateResponse{Task: createdTask}, nil
}

type ListRequest struct {
	UserID uint
}

type ListResponse struct {
	Tasks []entity.Task
}

func (t Service) List(req ListRequest) (ListResponse, error) {
	tasks, err := t.repository.ListUserTasks(req.UserID)
	if err != nil {

		return ListResponse{}, fmt.Errorf("can't list user tasks: %v\n", err)
	}

	return ListResponse{Tasks: tasks}, nil
}
