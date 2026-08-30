package shared

import (
	"fmt"
	"strconv"
	"strings"
)

func ParseFloat64(value any) (float64, error) {
	var sigma float64

	if value == nil {
		return 0, fmt.Errorf("sigma parameter is nil")
	}

	switch v := value.(type) {
	case float64:
		sigma = v
	case int:
		sigma = float64(v)
	case int64:
		sigma = float64(v)
	case string:
		if strings.TrimSpace(v) == "" {
			return 0, fmt.Errorf("sigma parameter is empty")
		} else {
			parsed, err := strconv.ParseFloat(v, 64)
			if err != nil {
				return 0, fmt.Errorf("failed to parse string sigma '%s' to float64: %w", v, err)
			}
			sigma = parsed
		}
	default:
		return 0, fmt.Errorf("invalid type for sigma parameter: expected float or string, got %T", value)
	}

	if sigma <= 0 {
		return 0, fmt.Errorf("sigma must be greater than 0, got %f", sigma)
	}

	return sigma, nil
}

func ValidateNoParameters(parameters map[string]any) error {
	if len(parameters) > 0 {
		return fmt.Errorf("filter does not accept parameters")
	}

	return nil
}

func ToLower(s string) string {
	return strings.ToLower(s)
}
