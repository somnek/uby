package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/charmbracelet/log"
	"github.com/gocolly/colly"
)

func toNum(s string) int {
	i, _ := strconv.Atoi(s)
	return i
}

type Dep struct {
	User   string `json:"user"`
	Repo   string `json:"repo"`
	Stars  int    `json:"stars"`
	Avatar string `json:"avatar"`
	Url    string `json:"url"`
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
				cleaned := strings.Replace(trimmed, ",", "", -1)
				count = toNum(cleaned)
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

		*deps = append(*deps, Dep{User: user,
			Repo:   repo,
			Stars:  stars,
			Avatar: avatar,
			Url:    url,
		})
	})

	// ----------------------------------------------
	c.OnRequest(func(r *colly.Request) {
		// fmt.Println("🛩️ Visiting", r.URL)
		log.Info("count: ", len(*deps), "🛩️ Visiting", r.URL)
	})

	c.Visit(url)
	return nextPageExist, nextUrl, *deps
}

func writeJson(deps []Dep) {
	jsonData, err := json.MarshalIndent(deps, "", "  ")
	if err != nil {
		log.Fatal(err)
	}

	err = os.WriteFile("deps.json", jsonData, 0644)
	if err != nil {
		log.Fatal(err)
	}
	log.Info("JSON data written to file: deps.json")
}

func main() {
	url := "https://github.com/aquasecurity/trivy/network/dependents"
	// url := "https://github.com/hwchase17/langchain/network/dependents"

	estimatedCount := extractCount(url)

	allDeps := []Dep{}
	hasNextPage := false

	hasNextPage, url, allDeps = scrape(url, &allDeps)

	for hasNextPage {
		hasNextPage, url, allDeps = scrape(url, &allDeps)
	}
	fmt.Println("estimated: ", estimatedCount)
	fmt.Println("found:     ", len(allDeps))

	writeJson(allDeps)
}
