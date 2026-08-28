package filter

import (
	"image"
	"image/color"
	"media-processing-platform/processor/internal/filter/shared"
)

type Negative struct {
}

func (n *Negative) Name() string {
	return "negative"
}

func (n *Negative) Apply(src image.Image, parameters map[string]any) (image.Image, error) {
	if err := shared.ValidateNoParameters(parameters); err != nil {
		return nil, err
	}

	bounds := src.Bounds()
	width := bounds.Dx()
	height := bounds.Dy()

	dst := image.NewRGBA(bounds)

	for y := range height {
		for x := range width {
			r, g, b, a := src.At(x, y).RGBA()

			r8 := uint8(r >> 8)
			g8 := uint8(g >> 8)
			b8 := uint8(b >> 8)
			a8 := uint8(a >> 8)

			invertedColor := color.RGBA{
				R: 255 - r8,
				G: 255 - g8,
				B: 255 - b8,
				A: a8,
			}

			dst.SetRGBA(x, y, invertedColor)
		}
	}

	return dst, nil
}
