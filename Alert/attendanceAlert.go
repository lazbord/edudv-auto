package Alert

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"

	model "edudv-auto/Model"

	"golang.org/x/net/html"
)

func GetAttendance(courses []model.Course) {
	for _, course := range courses {
		if course.Link == "" {
			continue
		}

		url := "https://www.leonard-de-vinci.net" + course.Link

		req, err := http.NewRequest("GET", url, nil)
		if err != nil {
			log.Fatal(err)
		}

		// Essential headers only
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

		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			log.Fatal(err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			fmt.Printf("Failed to fetch attendance for course %s, status code: %d\n", course.Name, resp.StatusCode)
			continue
		}

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			log.Fatal(err)
		}

		checkAlertPresence(course, string(body))
	}
}

func checkAlertPresence(course model.Course, body string) {
	doc, err := html.Parse(strings.NewReader(body))
	if err != nil {
		log.Fatal(err)
	}

	var alertMessage string

	// Traverse the HTML nodes
	var traverse func(*html.Node)
	traverse = func(n *html.Node) {
		if n.Type == html.ElementNode && n.Data == "div" {
			if hasClass(n, "alert alert-warning") {
				alertMessage = getNodeText(n)
			}
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			traverse(c)
		}
	}

	traverse(doc)

	if alertMessage != "" {
		fmt.Printf("PRESENCE IS UP %s:\n", course.Name)
	} else {
		fmt.Printf("No presence is up %s.\n", course.Name)
	}
}

// Helper function to check if an HTML node has a specific class
func hasClass(n *html.Node, class string) bool {
	for _, attr := range n.Attr {
		if attr.Key == "class" && strings.Contains(attr.Val, class) {
			return true
		}
	}
	return false
}

// Helper function to get text content inside an HTML node
func getNodeText(n *html.Node) string {
	var text strings.Builder
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		if c.Type == html.TextNode {
			text.WriteString(strings.TrimSpace(c.Data))
		}
	}
	return text.String()
}
