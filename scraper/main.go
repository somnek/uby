package main

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/gocolly/colly"
)

const (
	maxPerPage = 30
)

func toNum(s string) int {
	i, _ := strconv.Atoi(s)
	return i
}

type Dep struct {
	user   string
	repo   string
	stars  int
	avatar string
	url    string
}

func extractStars(e *colly.HTMLElement) int {
	parent := e.ChildText("span.color-fg-muted.text-bold.pl-3")
	split := strings.TrimSpace(strings.Split(parent, " ")[0])
	stars := toNum(split)
	return stars
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
	var nextPageUrl string

	c.OnHTML("a.btn.btn-outline.BtnGroup-item", func(e *colly.HTMLElement) {
		if e.Text == "Next" {
			nextPageUrl = e.Attr("href")
		}
		fmt.Println("hot", nextPageUrl)
	})

	// ----------------------------------------------
	dependents := []Dep{}

	c.OnHTML("div.Box-row", func(e *colly.HTMLElement) {
		avatar := e.ChildAttr("img", "src")
		user := e.ChildAttrs("a", "href")[0][1:]
		repo := e.ChildText("a.text-bold")
		url := fmt.Sprintf("https://github.com/%s/%s", user, repo)
		stars := extractStars(e)

		dependents = append(dependents, Dep{
			user:   user,
			repo:   repo,
			stars:  stars,
			avatar: avatar,
			url:    url,
		})

		if len(dependents)%maxPerPage == 0 {
			fmt.Println("nextpagelink", nextPageUrl)
			c.Visit(nextPageUrl)
		}
	})

	// ----------------------------------------------
	c.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting", r.URL)
	})

	c.Visit(url)

	fmt.Println("count:", count)
	fmt.Println("len", len(dependents))
	fmt.Println("nextPageUrl", nextPageUrl)
	// spew.Dump(dependents)
	return count
}

func main() {
	c := colly.NewCollector()
	url := "https://github.com/aquasecurity/trivy/network/dependents"
	scrape(c, url)
}
