package main

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/davecgh/go-spew/spew"
	"github.com/gocolly/colly"
)

func toNum(s string) int {
	i, _ := strconv.Atoi(s)
	return i
}

type Deps struct {
	user   string
	repo   string
	avatar string
	url    string
}

func scrape(c *colly.Collector, url string) int {

	// ----------------------------------------------
	var count int

	c.OnHTML("a.btn-link.selected", func(e *colly.HTMLElement) {
		lines := strings.Split(e.Text, "\n")
		for i, line := range lines {
			if strings.Contains(line, "Repositories") {
				trimmed := strings.TrimSpace(lines[i-1])
				count = toNum(trimmed)
			}
		}
	})

	// ----------------------------------------------
	dependents := []Deps{}

	c.OnHTML("div.Box-row", func(e *colly.HTMLElement) {
		avatar := e.ChildAttr("img", "src")
		user := e.ChildAttrs("a", "href")[0][1:]
		repo := e.ChildText("a.text-bold")
		url := fmt.Sprintf("https://github.com/%s/%s", user, repo)
		dependents = append(dependents, Deps{
			user:   user,
			repo:   repo,
			avatar: avatar,
			url:    url,
		})
	})

	// ----------------------------------------------
	c.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting", r.URL)
	})

	c.Visit(url)

	fmt.Println("count:", count)
	spew.Dump(dependents)
	return count
}

func main() {
	c := colly.NewCollector()
	url := "https://github.com/aquasecurity/trivy/network/dependents"
	scrape(c, url)
}
