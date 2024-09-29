package Clock

import (
	alert "edudv-auto/Alert"
	discord "edudv-auto/Discord"
	model "edudv-auto/Model"
	"edudv-auto/Scrape"
	"fmt"
	"math/rand"
	"strings"
	"time"
)

// MainClock checks the schedule of courses and waits for the start time to check presence.
func ClockForTHeDay(courses []model.Course) {
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
				discord.SendLoggerMessage(" :zzz: Waiting for course " + course.Name + " to start at " + startTimeToday.Format("15:04"))
				time.Sleep(time.Until(startTimeToday)) // Sleep until the course start time
				currentTime = time.Now()               // Update current time after sleeping
			}

			// Start checking for presence every 30 seconds until the end time
			CheckPresence(course, endTimeToday)
		}

		discord.SendLoggerMessage(":white_check_mark: No more course for the end of the day.")
	}
}

// CheckPresence checks the presence status every 30 seconds until it's available or time runs out.
func CheckPresence(course model.Course, endTime time.Time) {
	discord.SendLoggerMessage(" :mag: Started checking presence for **" + course.Name + "**")
	for {
		// Call alert.CheckPresence to see if presence is available
		isPresenceUp := alert.CheckPresence(course)

		if isPresenceUp {
			// Send a message to Discord if presence is up
			discord.SendDiscordMessage(":warning: **Appel** @everyone :warning:")
			discord.SendDiscordMessage("- " + course.Name)
			discord.SendDiscordMessage("- " + course.Teacher)
			discord.SendDiscordMessage("- " + course.Hours)
			discord.SendDiscordMessage("- " + course.ZoomLink)
			discord.SendDiscordMessage("- " + course.DVLLink)
			discord.SendLoggerMessage(" :white_check_mark: Presence is up for  **" + course.Name + "**. Stopping presence check ")
			break // Exit the loop once presence is up
		}

		// Check if current time has passed the course end time
		if time.Now().After(endTime) {
			discord.SendLoggerMessage(" :no_entry: Course **" + course.Name + "** has ended. Stopping presence check ")
			break
		}

		randomSleep := time.Duration(20+rand.Intn(40)) * time.Second

		time.Sleep(randomSleep)
	}
}

func MainClock() {
	for {
		now := time.Now()

		next := time.Date(now.Year(), now.Month(), now.Day(), 7, 0, 0, 0, now.Location())

		if now.After(next) {
			next = next.Add(24 * time.Hour)
		}

		durationUntilNext := time.Until(next)

		waitingTimeMessage := fmt.Sprintf(":alarm_clock: Waiting until next morning, waiting time: %dh%dm", int(durationUntilNext.Hours()), int(durationUntilNext.Minutes())%60)

		discord.SendLoggerMessage(waitingTimeMessage)

		time.Sleep(durationUntilNext)

		courses := Scrape.GetCoursesOfTheDay()
		ClockForTHeDay(courses)
	}
}
