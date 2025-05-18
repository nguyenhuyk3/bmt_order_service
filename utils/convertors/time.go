package convertors

import (
	"fmt"
	"time"
)

func ConvertDateStringToTime(input string) (time.Time, error) {
	parsedTime, err := time.Parse("2006-01-02", input)
	if err != nil {
		return time.Time{}, fmt.Errorf("invalid date format (%s): %w", input, err)
	}

	return parsedTime, nil
}
