package main

import (
	"edudv-auto/Alert"
	"edudv-auto/Scrape"
	"fmt"
)

func main() {
	SourceCode := Scrape.GetSourceCode()
	courses := Scrape.ParseHTML(SourceCode)

	// Print the parsed courses
	for _, course := range courses {
		fmt.Printf("Course:\n Hours: %s\n Name: %s\n Teacher: %s\n Link: %s\n ZoomLink: %s\n DVLLink: %s\n\n",
			course.Hours, course.Name, course.Teacher, course.Link, course.ZoomLink, course.DVLLink)
	}

	Alert.GetAttendance(courses)
}
