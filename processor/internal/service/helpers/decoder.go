package helpers

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"image"
	pb "media-processing-platform/pkg/proto"
	"media-processing-platform/processor/internal/domain"
	"strings"

	"github.com/disintegration/imaging"
	"github.com/google/uuid"
)

func DecodeTask(taskMsg *pb.TaskMessage) (*domain.Task, error) {
	taskID, err := uuid.Parse(taskMsg.GetId())
	if err != nil {
		return nil, fmt.Errorf("failed to parse task id: %w", err)
	}

	filter := taskMsg.GetFilter()
	if filter == nil {
		return nil, fmt.Errorf("filter is nil")
	}

	parametersJSON := filter.GetParameters()
	if parametersJSON == nil {
		return nil, fmt.Errorf("filter parameters are nil")
	}

	img, err := DecodeImage(taskMsg.GetImage())
	if err != nil {
		return nil, fmt.Errorf("failed to decode task image: %w", err)
	}

	return &domain.Task{
		ID:    taskID,
		Image: img,
		Filter: &domain.ImageFilter{
			Name:       filter.GetName(),
			Parameters: parametersJSON.AsMap(),
		},
	}, nil
}

func DecodeImage(base64Str string) (image.Image, error) {
	idx := strings.Index(base64Str, ",")
	if idx != -1 {
		base64Str = base64Str[idx+1:]
	}

	imgBytes, err := base64.StdEncoding.DecodeString(base64Str)
	if err != nil {
		return nil, fmt.Errorf("base64 decode error: %w", err)
	}

	srcImg, err := imaging.Decode(bytes.NewReader(imgBytes))
	if err != nil {
		return nil, fmt.Errorf("image parse error: %w", err)
	}

	return srcImg, nil
}
