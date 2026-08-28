package filter

import (
	"fmt"
	"image"
	"media-processing-platform/processor/internal/filter/shared"

	"github.com/disintegration/imaging"
)

type Blur struct {
}

func (b *Blur) Name() string {
	return "blur"
}

func (b *Blur) Apply(src image.Image, parameters map[string]any) (image.Image, error) {
	value, ok := parameters["sigma"]
	if !ok {
		return nil, fmt.Errorf("sigma parameter is missing")
	}

	sigma, err := shared.ParseFloat64(value)
	if err != nil {
		return nil, fmt.Errorf("failed to parse sigma parameter: %w", err)
	}

	dstImg := imaging.Blur(src, sigma)
	return dstImg, nil
}
