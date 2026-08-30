package helpers

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"image"
	"image/png"
)

func EncodeImage(image image.Image) (string, error) {
	var buf bytes.Buffer
	err := png.Encode(&buf, image)
	if err != nil {
		return "", fmt.Errorf("PNG encode error: %w", err)
	}

	resultBase64 := base64.StdEncoding.EncodeToString(buf.Bytes())
	return "data:image/png;base64," + resultBase64, nil
}
