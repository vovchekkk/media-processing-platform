package service

import (
	"context"
	"log/slog"
	"time"

	"github.com/google/uuid"

	"media-processing-platform/internal/config"
	"media-processing-platform/internal/domain"
	"media-processing-platform/internal/repository"
)

type TaskProcessor struct {
	processingDuration time.Duration
	repository         repository.Task
	log                *slog.Logger
}

func NewTaskProcessor(cfg config.TaskProcessorConfig, repository repository.Task, log *slog.Logger) *TaskProcessor {
	return &TaskProcessor{
		processingDuration: cfg.ProcessingDuration,
		repository:         repository,
		log:                log,
	}
}

func (taskProccessor *TaskProcessor) startTaskProcessing(ctx context.Context, id uuid.UUID) {
	if err := taskProccessor.repository.UpdateTaskStatus(
		ctx,
		id,
		domain.StatusInProgress,
	); err != nil {
		taskProccessor.log.Error("failed to set task status to in progress",
			"task_id", id,
			"error", err,
		)
		return
	}

	taskProccessor.log.Info("task processing started", "task_id", id)
}

func (taskProccessor *TaskProcessor) doVeryExpensiveTask(ctx context.Context, id uuid.UUID) {
	time.Sleep(taskProccessor.processingDuration)

	if err := taskProccessor.repository.SetTaskResult(
		ctx,
		id,
		"Good job! Task completed successfully.",
	); err != nil {
		taskProccessor.log.Error("failed to set task result",
			"task_id", id,
			"error", err,
		)
		return
	}
}

func (taskProccessor *TaskProcessor) finishTaskProcessing(ctx context.Context, id uuid.UUID) {
	if err := taskProccessor.repository.UpdateTaskStatus(
		ctx,
		id,
		domain.StatusReady,
	); err != nil {
		taskProccessor.log.Error("failed to set task status to ready",
			"task_id", id,
			"error", err,
		)
		return
	}

	taskProccessor.log.Info("task processing completed", "task_id", id)
}

func (taskProccessor *TaskProcessor) Process(id uuid.UUID) error {
	go func() {
		ctx := context.Background()

		taskProccessor.startTaskProcessing(ctx, id)

		taskProccessor.doVeryExpensiveTask(ctx, id)

		taskProccessor.finishTaskProcessing(ctx, id)
	}()

	return nil
}
