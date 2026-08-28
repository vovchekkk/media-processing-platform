package service

import (
	"context"
	"errors"
	"fmt"
	pb "media-processing-platform/pkg/proto"
	"media-processing-platform/server/internal/infrastructure/rabbitmq"

	"github.com/google/uuid"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/structpb"
	"gorm.io/gorm"

	"media-processing-platform/server/internal/domain"
	"media-processing-platform/server/internal/dto"
	"media-processing-platform/server/internal/repository"
)

type TaskService struct {
	taskRepository repository.Task
	producer       *rabbitmq.Producer
}

func NewTaskService(taskRepo repository.Task, producer *rabbitmq.Producer) *TaskService {
	return &TaskService{
		taskRepository: taskRepo,
		producer:       producer,
	}
}

func (taskService *TaskService) Create(ctx context.Context, taskDTO *dto.Task, userID uuid.UUID) (uuid.UUID, error) {
	id, err := uuid.NewRandom()
	if err != nil {
		return uuid.Nil, err
	}

	task := &domain.Task{
		ID:     id,
		UserID: userID,
		Filter: taskDTO.Filter,
		Image:  taskDTO.Image,
	}

	if err := taskService.taskRepository.CreateTask(ctx, task); err != nil {
		return uuid.Nil, err
	}

	protoParams, err := structpb.NewStruct(taskDTO.Filter.Parameters)
	if err != nil {
		return uuid.Nil, err
	}

	protoMsg := &pb.TaskMessage{
		Id:    task.ID.String(),
		Image: task.Image,
		Filter: &pb.Filter{
			Name:           taskDTO.Filter.Name,
			ParametersJson: protoParams,
		},
	}

	body, err := proto.Marshal(protoMsg)
	if err != nil {
		return uuid.Nil, err
	}

	if prodErr := taskService.producer.Publish(ctx, body); prodErr != nil {
		if updateErr := taskService.taskRepository.UpdateTaskStatusAndResult(ctx, task.ID, userID, "failed", "failed to publish to queue"); updateErr != nil {
			return uuid.Nil, fmt.Errorf("publication error: %w (suppressed DB error: %v)", prodErr, updateErr)
		}

		return uuid.Nil, prodErr
	}

	return id, nil
}

func (taskService *TaskService) GetResult(ctx context.Context, id uuid.UUID, userID uuid.UUID) (string, error) {
	result, err := taskService.taskRepository.GetTaskResult(ctx, id, userID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return "", domain.ErrTaskNotFound
		}

		return "", err
	}

	return result, nil
}

func (taskService *TaskService) GetStatus(ctx context.Context, id uuid.UUID, userID uuid.UUID) (domain.TaskStatus, error) {
	status, err := taskService.taskRepository.GetTaskStatus(ctx, id, userID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return "", domain.ErrTaskNotFound
		}

		return "", err
	}

	return status, nil
}
