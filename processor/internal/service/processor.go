package service

import (
	"context"
	"fmt"
	"log/slog"
	pb "media-processing-platform/pkg/proto"
	"media-processing-platform/processor/internal/domain"
	"media-processing-platform/processor/internal/metrics"
	"media-processing-platform/processor/internal/repository"
	"media-processing-platform/processor/internal/service/helpers"
	"time"

	"google.golang.org/protobuf/proto"
)

type ProcessorService struct {
	taskRepo     repository.Task
	imageService *ImageService
	metrics      *metrics.Metrics
	log          *slog.Logger
}

func NewProcessorService(taskRepo repository.Task, imageService *ImageService, metrics *metrics.Metrics, log *slog.Logger) *ProcessorService {
	return &ProcessorService{
		taskRepo:     taskRepo,
		imageService: imageService,
		metrics:      metrics,
		log:          log,
	}
}

func (processor *ProcessorService) Process(ctx context.Context, body []byte) (err error) {
	var taskMsg pb.TaskMessage
	if err := proto.Unmarshal(body, &taskMsg); err != nil {
		return fmt.Errorf("failed to unmarshal proto: %w", err)
	}

	task, decodeErr := helpers.DecodeTask(&taskMsg)
	if decodeErr != nil {
		return fmt.Errorf("failed to decode task: %w", decodeErr)
	}

	start := time.Now()
	filterName := task.Filter.Name

	defer func() {
		processor.metrics.TaskDuration.WithLabelValues(filterName).Observe(time.Since(start).Seconds())

		if err != nil {
			processor.metrics.TasksTotal.WithLabelValues(filterName, "failed").Inc()

			if updateErr := processor.taskRepo.UpdateTaskStatusAndResult(
				context.Background(),
				task.ID,
				domain.StatusFailed,
				err.Error(),
			); updateErr != nil {
				processor.log.Error(
					"failed to set task status to failed",
					"task_id", task.ID,
					"error", updateErr,
				)
			}
		} else {
			processor.metrics.TasksTotal.WithLabelValues(filterName, "success").Inc()
		}
	}()

	image, err := processor.imageService.Process(task)
	if err != nil {
		processor.log.Error("task processing failed", "task_id", task.ID, "error", err)
		return err
	}

	resultData, err := helpers.EncodeImage(image)
	if err != nil {
		return fmt.Errorf("failed to encode image: %w", err)
	}

	err = processor.taskRepo.UpdateTaskStatusAndResult(ctx, task.ID, domain.StatusReady, resultData)
	if err != nil {
		processor.log.Error("failed to save result", "task_id", task.ID, "error", err)
		return err
	}

	processor.log.Info("task processing completed and result saved", "task_id", task.ID, "result", resultData)
	return nil
}
