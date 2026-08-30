package filter

import "image"

type Filter interface {
	Name() string
	Apply(src image.Image, parameters map[string]any) (image.Image, error)
}
