package filter

import (
	"fmt"
	"image"
	"media-processing-platform/processor/internal/filter/shared"
)

type ApplyFunc func(src image.Image, parameters map[string]any) (image.Image, error)

type Registry struct {
	filters map[string]ApplyFunc
}

func NewRegistry(filters ...Filter) *Registry {
	registry := &Registry{
		filters: make(map[string]ApplyFunc),
	}

	for _, filter := range filters {
		registry.filters[filter.Name()] = filter.Apply
	}

	return registry
}

func (registry *Registry) GetApplyFunc(filterName string) (ApplyFunc, error) {
	apply, ok := registry.filters[shared.ToLower(filterName)]
	if !ok {
		return nil, fmt.Errorf("unknown filter: %s", filterName)
	}

	return apply, nil
}
