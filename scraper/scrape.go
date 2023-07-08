package main

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/log"
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

func formatUrl(url string) string {
	return fmt.Sprintf("https://github.com/%s/network/dependents", url)
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
				cleaned := strings.ReplaceAll(trimmed, ",", "")
				count = toNum(cleaned)
			}
		}
	})

	err := c.Visit(url)
	if err != nil {
		log.Info(err)
	}
	return count
}

func Scrape(url string) tea.Cmd {
	url = formatUrl(url)
	return func() tea.Msg {
		count := extractCount(url)
		return Done(fmt.Sprintf("done %d", count))
	}
}
