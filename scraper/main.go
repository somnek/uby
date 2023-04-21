package main

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/gocolly/colly"
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
	nextPageExist := true
	page := 1
	maxPerPage := 30
	depSizeOnPage := 0 // reset to 0 on each page

	// ----------------------------------------------
	// count
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
	// next page button
	var nextPageUrl string

	c.OnHTML("a.btn.btn-outline.BtnGroup-item", func(e *colly.HTMLElement) {
		if e.Text == "Next" {
			nextPageUrl = e.Attr("href")
			nextPageExist = true
		} else {
			nextPageExist = false
		}
	})

	// ----------------------------------------------
	// list of deps
	dependents := []Dep{}

	c.OnHTML("div.Box-row", func(e *colly.HTMLElement) {
		avatar := e.ChildAttr("img", "src")
		user := e.ChildAttrs("a", "href")[0][1:]
		repo := e.ChildText("a.text-bold")
		stars := extractStars(e)

		dependents = append(dependents, Dep{
			user:   user,
			repo:   repo,
			stars:  stars,
			avatar: avatar,
			url:    url,
		})
		depSizeOnPage += 1

		if depSizeOnPage == maxPerPage {
			url = nextPageUrl
			page += 1
			depSizeOnPage = 0
		} else {
			return
		}
	})

	// ----------------------------------------------
	c.OnRequest(func(r *colly.Request) {
		fmt.Println("🛩️ Visiting", r.URL)
	})

	for nextPageExist {
		c.Visit(url)
	}

	fmt.Println("count:", count)
	fmt.Println("len", len(dependents))
	// spew.Dump(dependents)
	return count
}

func main() {
	c := colly.NewCollector()
	// url := "https://github.com/aquasecurity/trivy/network/dependents"
	url := "https://github.com/charmbracelet/bubbletea/network/dependents?package_id=UGFja2FnZS0yMjc1ODk0MDQy"
	scrape(c, url)
}
