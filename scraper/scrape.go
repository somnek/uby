package main

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/gocolly/colly"
)

type Dep struct {
	User    string `json:"user"`
	Repo    string `json:"repo"`
	Stars   int    `json:"stars"`
	Avatar  string `json:"avatar"`
	RepoUrl string `json:"repoUrl"`
	Url     string `json:"depUrl"`
}

type Done string

func SomeLongTask() tea.Cmd {
	return func() tea.Msg {
		// ...
		time.Sleep(2 * time.Second)
		return Done("done")
	}
}

func toNum(s string) int {
	i, _ := strconv.Atoi(s)
	return i
}

func extractStars(e *colly.HTMLElement) int {
	parent := e.ChildText("span.color-fg-muted.text-bold.pl-3")
	split := strings.TrimSpace(strings.Split(parent, " ")[0])
	stars := toNum(split)
	return stars
}

func scrapePage(url string) (bool, string) {
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

	err := c.Visit(url)
	if err != nil {
		return false, ""
	}
	return nextPageExist, nextUrl
}

func Scrape(url string) tea.Cmd {
	return func() tea.Msg {
		var nextUrl string
		var pages int
		var shouldContinue = true

		// next pages
		for shouldContinue {
			if pages == 0 {
				shouldContinue, nextUrl = scrapePage(url)
			} else {
				shouldContinue, nextUrl = scrapePage(nextUrl)
			}
			pages += 1
		}

		return Done(fmt.Sprintf("%d total repo...", pages))
	}
}
