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

func extractCount(url string) int {
	c := colly.NewCollector()
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

	c.Visit(url)
	return count
}

func scrape(url string, deps *[]Dep) (bool, string, []Dep) {
	c := colly.NewCollector()

	// ----------------------------------------------
	// next page button
	nextPageExist := false
	var nextUrl string

	c.OnHTML("a.btn.btn-outline.BtnGroup-item", func(e *colly.HTMLElement) {
		if e.Text == "Next" {
			nextPageExist = true
		} else {
			nextPageExist = false
		}
		nextUrl = e.Attr("href")
	})

	// ----------------------------------------------
	// list of deps
	c.OnHTML("div.Box-row", func(e *colly.HTMLElement) {
		avatar := e.ChildAttr("img", "src")
		user := e.ChildAttrs("a", "href")[0][1:]
		repo := e.ChildText("a.text-bold")
		stars := extractStars(e)

		*deps = append(*deps, Dep{user: user,
			repo:   repo,
			stars:  stars,
			avatar: avatar,
			url:    url,
		})
	})

	// ----------------------------------------------
	c.OnRequest(func(r *colly.Request) {
		fmt.Println("🛩️ Visiting", r.URL)
	})

	c.Visit(url)
	return nextPageExist, nextUrl, *deps
}

func main() {
	// url := "https://github.com/aquasecurity/trivy/network/dependents"
	url := "https://github.com/hwchase17/langchain/network/dependents"

	estimatedCount := extractCount(url)

	allDeps := []Dep{}
	hasNextPage := false

	hasNextPage, url, allDeps = scrape(url, &allDeps)

	for hasNextPage {
		hasNextPage, url, allDeps = scrape(url, &allDeps)
	}
	fmt.Println("estimated: ", estimatedCount)
	fmt.Println("found:     ", len(allDeps))

}
