package filter

import (
	"image"
	"media-processing-platform/processor/internal/filter/shared"
)

type FlipX struct {
}

func (f *FlipX) Name() string {
	return "flip_x"
}

func (f *FlipX) Apply(src image.Image, parameters map[string]any) (image.Image, error) {
	if err := shared.ValidateNoParameters(parameters); err != nil {
		return nil, err
	}

	bounds := src.Bounds()
	width := bounds.Dx()
	height := bounds.Dy()

	dst := image.NewRGBA(bounds)

	for y := range height {
		for x := range width {
			newX := width - 1 - x

			color := src.At(bounds.Min.X+x, bounds.Min.Y+y)

			dst.Set(bounds.Min.X+newX, bounds.Min.Y+y, color)
		}
	}

	return dst, nil
}
