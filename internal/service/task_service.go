package service

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"media-processing-platform/internal/domain"
	"media-processing-platform/internal/dto"
	"media-processing-platform/internal/repository"
)

type TaskService struct {
	taskRepository repository.Task
	taskProcessor  *TaskProcessor
}

func NewTaskService(taskRepo repository.Task, taskProcessor *TaskProcessor) *TaskService {
	return &TaskService{
		taskRepository: taskRepo,
		taskProcessor: taskProcessor,
	}
}

func (taskService *TaskService) Create(ctx context.Context, taskDTO *dto.Task) (uuid.UUID, error) {
	id, err := uuid.NewRandom()
	if err != nil {
		return uuid.Nil, err
	}

	task := &domain.Task{
		ID: id,

	}

	if err := taskService.taskRepository.CreateTask(ctx, task); err != nil {
		return uuid.Nil, err
	}

	if err := taskService.taskProcessor.Process(id); err != nil {
		return uuid.Nil, err
	}

	return id, nil
}

func (taskService *TaskService) GetResult(ctx context.Context, id uuid.UUID) (string, error) {
	result, err := taskService.taskRepository.GetTaskResultByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return "", domain.ErrTaskNotFound
		}

		return "", err
	}

	return result, nil
}

func (taskService *TaskService) GetStatus(ctx context.Context, id uuid.UUID) (domain.TaskStatus, error) {
	status, err := taskService.taskRepository.GetTaskStatusByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return "", domain.ErrTaskNotFound
		}

		return "", err
	}

	return status, nil
}
