package main

import (
	"edudv-auto/Clock"
	"edudv-auto/Scrape"
)

func main() {
	courses := Scrape.GetCoursesOfTheDay()
	Clock.MainClock(courses)
}
