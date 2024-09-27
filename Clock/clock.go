package Clock

import (
	alert "edudv-auto/Alert"
	discord "edudv-auto/Discord"
	model "edudv-auto/Model"
	"strings"
	"time"
)

// MainClock checks the schedule of courses and waits for the start time to check presence.
func MainClock(courses []model.Course) {
	currentTime := time.Now()

	for _, course := range courses {
		// Split multiple time ranges, if provided
		timeRanges := strings.Split(course.Hours, ",")

		for _, rangeStr := range timeRanges {
			// Parse the time range (start-end)
			times := strings.Split(strings.TrimSpace(rangeStr), "-")
			startTime, _ := time.Parse("15:04", times[0])
			endTime, _ := time.Parse("15:04", times[1])

			// Set course start and end times to today
			startTimeToday := time.Date(currentTime.Year(), currentTime.Month(), currentTime.Day(), startTime.Hour(), startTime.Minute(), 0, 0, time.Local)
			endTimeToday := time.Date(currentTime.Year(), currentTime.Month(), currentTime.Day(), endTime.Hour(), endTime.Minute(), 0, 0, time.Local)

			// Check if the current time is before the start time
			if currentTime.Before(startTimeToday) {
				discord.SendLoggerMessage("Waiting for course " + course.Name + " to start at " + startTimeToday.Format("15:04"))
				time.Sleep(time.Until(startTimeToday)) // Sleep until the course start time
				currentTime = time.Now()               // Update current time after sleeping
			}

			// Start checking for presence every 30 seconds until the end time
			CheckPresence(course, endTimeToday)
		}
	}
}

// CheckPresence checks the presence status every 30 seconds until it's available or time runs out.
func CheckPresence(course model.Course, endTime time.Time) {
	for {
		// Call alert.CheckPresence to see if presence is available
		isPresenceUp := alert.CheckPresence(course)

		if isPresenceUp {
			// Send a message to Discord if presence is up
			discord.SendDiscordMessage("Appel !")
			discord.SendDiscordMessage("- " + course.Name)
			discord.SendDiscordMessage("- " + course.Teacher)
			discord.SendDiscordMessage("- " + course.Hours)
			discord.SendDiscordMessage("- " + course.ZoomLink)
			discord.SendDiscordMessage("- " + course.DVLLink)
			discord.SendLoggerMessage("Presence is up for  " + course.Name + ". Stopping presence check ")
			break // Exit the loop once presence is up
		}

		// Check if current time has passed the course end time
		if time.Now().After(endTime) {
			discord.SendLoggerMessage("Course " + course.Name + " has ended. Stopping presence check ")
			break
		}

		// Wait for 30 seconds before checking again
		time.Sleep(30 * time.Second)
	}
}
