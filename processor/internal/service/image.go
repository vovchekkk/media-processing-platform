package service

import (
	"image"
	"media-processing-platform/processor/internal/domain"
	"media-processing-platform/processor/internal/filter"
)

type ImageService struct {
	registry *filter.Registry
}

func NewImageService(registry *filter.Registry) *ImageService {
	return &ImageService{
		registry: registry,
	}
}

func (imageService *ImageService) Process(task *domain.Task) (image.Image, error) {
	apply, err := imageService.registry.GetApplyFunc(task.Filter.Name)
	if err != nil {
		return nil, err
	}

	return apply(task.Image, task.Filter.Parameters)
}
