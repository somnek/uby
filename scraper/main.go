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

type Deps struct {
	name string
}

func scrapeRepoCount(url string) int {
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

	c.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting", r.URL)
	})

	c.Visit(url)
	fmt.Println(count)
	return count
}

func scrapeDependents(repoCount int) []Deps {
	hasNextPage := false
	dependents := []Deps{}

	for hasNextPage {
		dependents = append(dependents, Deps{name: "test"})
		break
	}
	return dependents
}

func main() {
	url := "https://github.com/aquasecurity/trivy/network/dependents"
	repoCount := scrapeRepoCount(url)
	dependents := scrapeDependents(repoCount)
	_ = dependents

}
