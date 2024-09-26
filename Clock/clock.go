package Clock

import (
	discord "edudv-auto/Discord"
	model "edudv-auto/Model"
	"fmt"
	"strings"
	"time"
)

func GetHours(courses []model.Course) {
	currentTime := time.Now()

	for _, course := range courses {
		// Split the hours string into individual time ranges
		timeRanges := strings.Split(course.Hours, ",")

		for _, rangeStr := range timeRanges {
			// Split the range string into start and end times
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
				fmt.Printf(course.Name)
			}
		}
	}
	discord.SendDiscordMessage("Pas de cours prévu à " + currentTime.Format("15:04"))
}
