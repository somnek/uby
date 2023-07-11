package main

import (
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

type ScrapeTick struct {
	nextUrl string
}
type InitScrapeTick string

func scrapePage(url string) string {
	c := colly.NewCollector()
	var nextUrl string

	// next page button
	c.OnHTML("a.btn.btn-outline.BtnGroup-item", func(e *colly.HTMLElement) {
		if e.Text == "Next" {
			nextUrl = e.Attr("href")
		} else {
			nextUrl = ""
		}
	})

	err := c.Visit(url)
	if err != nil {
		return ""
	}
	return nextUrl
}

func InitScrape() tea.Cmd {
	return func() tea.Msg {
		return InitScrapeTick("Tick")
	}
}

func Scrape(url string) tea.Cmd {
	return func() tea.Msg {
		nextUrl := scrapePage(url)

		return ScrapeTick{
			nextUrl: nextUrl,
		}
	}
}
