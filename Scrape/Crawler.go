package Scrape

import (
	"io"
	"log"
	"net/http"
	"strings"

	model "edudv-auto/Model"

	"golang.org/x/net/html"
)

func GetSourceCode() string {
	client := &http.Client{}
	req, err := http.NewRequest("GET", "https://www.leonard-de-vinci.net/student/presences/", nil)
	if err != nil {
		log.Fatal(err)
	}
	req.Header.Set("accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.7")
	req.Header.Set("accept-language", "fr-FR,fr;q=0.9,en-US;q=0.8,en;q=0.7")
	req.Header.Set("cache-control", "max-age=0")
	req.Header.Set("cookie", "SimpleSAML=1h9ltt70u2v2lstqso09s8h63s; SimpleSAMLAuthToken=_841b5008a9000756a1945515c4634f71b87a0caab6; alvstu=amtdrkbf1mhqm6s7hvsj5k7dcn; uids=lb201460")
	req.Header.Set("priority", "u=0, i")
	req.Header.Set("referer", "https://www.leonard-de-vinci.net/")
	req.Header.Set("sec-ch-ua", `"Google Chrome";v="129", "Not=A?Brand";v="8", "Chromium";v="129"`)
	req.Header.Set("sec-ch-ua-mobile", "?0")
	req.Header.Set("sec-ch-ua-platform", `"Linux"`)
	req.Header.Set("sec-fetch-dest", "document")
	req.Header.Set("sec-fetch-mode", "navigate")
	req.Header.Set("sec-fetch-site", "same-origin")
	req.Header.Set("sec-fetch-user", "?1")
	req.Header.Set("upgrade-insecure-requests", "1")
	req.Header.Set("user-agent", "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/129.0.0.0 Safari/537.36")
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()
	bodyText, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	return (string(bodyText))
}

func ParseHTML(SourceCode string) []model.Course {
	doc, err := html.Parse(strings.NewReader(SourceCode))
	if err != nil {
		log.Fatal(err)
	}

	var courses []model.Course

	// Function to traverse the HTML nodes
	var traverse func(*html.Node)
	traverse = func(n *html.Node) {
		if n.Type == html.ElementNode && n.Data == "tr" {
			var course model.Course
			tdIndex := 0

			// Traverse <td> inside the <tr>
			for td := n.FirstChild; td != nil; td = td.NextSibling {
				if td.Type == html.ElementNode && td.Data == "td" {
					text := getText(td)

					// Assign fields based on the order of <td> elements
					switch tdIndex {
					case 0:
						cleanTimeRange := strings.ReplaceAll(text, " ", "")
						course.Hours = cleanTimeRange
					case 1:
						course.Name = text
					case 2:
						course.Teacher = text
					case 3:
						course.Link = getHref(td)
					case 4:
						course.ZoomLink = getHref(td)
					case 5:
						course.DVLLink = getHref(td)
					}
					tdIndex++
				}
			}
			// Add the course to the list if it's valid
			if course.Hours != "" {
				courses = append(courses, course)
			}
		}

		// Traverse the children of the current node
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			traverse(c)
		}
	}

	// Start traversing from the root
	traverse(doc)
	return courses
}

// Helper to get text inside a <td>
func getText(n *html.Node) string {
	var content strings.Builder
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		if c.Type == html.TextNode {
			content.WriteString(strings.TrimSpace(c.Data))
		}
	}
	return content.String()
}

// Helper to get href inside an <a> tag
func getHref(n *html.Node) string {
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		if c.Type == html.ElementNode && c.Data == "a" {
			for _, attr := range c.Attr {
				if attr.Key == "href" {
					return attr.Val
				}
			}
		}
	}
	return ""
}
