package clock

import (
	model "edudv-auto/Model"
	"strings"
	"time"
)

func getHours(courses []model.Course) bool {
	currentTime := time.Now()
	currentHour := currentTime.Format("15:04") // Current time in "HH:MM"

	for _, course := range courses {
		for _, rangeStr := range course.Hours {
			// Split the range string
			times := strings.Split(strings.TrimSpace(rangeStr), " - ")
			if len(times) != 2 {
				continue // Skip malformed entries
			}

			startTime, err1 := time.Parse("15:04", times[0])
			endTime, err2 := time.Parse("15:04", times[1])

			if err1 != nil || err2 != nil {
				continue // Skip if parsing fails
			}

			// Check if the current time is within the range
			if currentTime.After(startTime) && currentTime.Before(endTime) {
				return true // Current hour is within the specified ranges
			}
		}
	}
	return false
}
