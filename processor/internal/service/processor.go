package service

import (
	"context"
	"fmt"
	"log/slog"
	pb "media-processing-platform/pkg/proto"
	"media-processing-platform/processor/internal/domain"
	"media-processing-platform/processor/internal/repository"

	"github.com/google/uuid"
	"google.golang.org/protobuf/proto"
)

type ProcessorService struct {
	taskRepo repository.Task
	log      *slog.Logger
}

func NewProcessorService(taskRepo repository.Task, log *slog.Logger) *ProcessorService {
	return &ProcessorService{
		taskRepo: taskRepo,
		log:      log,
	}
}

func (processor *ProcessorService) Process(ctx context.Context, body []byte) error {
	var taskMsg pb.TaskMessage
	if err := proto.Unmarshal(body, &taskMsg); err != nil {
		return fmt.Errorf("failed to unmarshal proto: %w", err)
	}

	taskID, err := uuid.Parse(taskMsg.GetId())
	if err != nil {
		return fmt.Errorf("failed to parse task id: %w", err)
	}

	resultData, execErr := processor.executeProcessing(&taskMsg)
	if execErr != nil {
		processor.log.Error("task processing failed", "task_id", taskID, "error", execErr)
		if updateErr := processor.taskRepo.UpdateTaskStatusAndResult(ctx, taskID, domain.StatusFailed, execErr.Error()); updateErr != nil {
			processor.log.Error("failed to update task status", "task_id", taskID, "error", updateErr)
			return fmt.Errorf("execution error: %w (suppressed DB error: %v)", execErr, updateErr)
		}

		return execErr
	}

	err = processor.taskRepo.UpdateTaskStatusAndResult(ctx, taskID, domain.StatusReady, resultData)
	if err != nil {
		processor.log.Error("failed to save completed status to DB", "task_id", taskID, "error", err)
		return err
	}

	processor.log.Info("task processing completed and result saved", "task_id", taskID, "result", resultData)
	return nil
}

func (processor *ProcessorService) executeProcessing(msg *pb.TaskMessage) (string, error) {
	return "data:image/png;base64,...", nil
}
