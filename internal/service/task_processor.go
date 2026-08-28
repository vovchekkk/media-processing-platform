package service

import (
	"context"
	"log/slog"
	"time"

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

func (taskProccessor *TaskProcessor) startTaskProcessing(ctx context.Context, task *domain.Task) {
	if err := taskProccessor.repository.UpdateTaskStatusByIDAndUserID(
		ctx,
		task.ID,
		task.UserID,
		domain.StatusInProgress,
	); err != nil {
		taskProccessor.log.Error("failed to set task status to in progress",
			"task_id", task.ID,
			"error", err,
		)
		return
	}

	taskProccessor.log.Info("task processing started", "task_id", task.ID)
}

func (taskProccessor *TaskProcessor) doVeryExpensiveTask(ctx context.Context, task *domain.Task) {
	time.Sleep(taskProccessor.processingDuration)

	if err := taskProccessor.repository.SetTaskResultByIDAndUserID(
		ctx,
		task.ID,
		task.UserID,
		"Good job! Task completed successfully.",
	); err != nil {
		taskProccessor.log.Error("failed to set task result",
			"task_id", task.ID,
			"error", err,
		)
		return
	}
}

func (taskProccessor *TaskProcessor) finishTaskProcessing(ctx context.Context, task *domain.Task) {
	if err := taskProccessor.repository.UpdateTaskStatusByIDAndUserID(
		ctx,
		task.ID,
		task.UserID,
		domain.StatusReady,
	); err != nil {
		taskProccessor.log.Error("failed to set task status to ready",
			"task_id", task.ID,
			"error", err,
		)
		return
	}

	taskProccessor.log.Info("task processing completed", "task_id", task.ID)
}

func (taskProccessor *TaskProcessor) Process(task *domain.Task) error {
	go func() {
		ctx := context.Background()

		taskProccessor.startTaskProcessing(ctx, task)

		taskProccessor.doVeryExpensiveTask(ctx, task)

		taskProccessor.finishTaskProcessing(ctx, task)
	}()

	return nil
}
