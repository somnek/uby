package main

import (
	"fmt"

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
	Page    int    `json:"page"`
}

type PageTick struct {
	nextUrl string
	deps    []Dep
}
type InitScrapeTick string

func scrapePage(url string, deps []Dep) (string, []Dep) {
	c := colly.NewCollector()
	var nextUrl string

	// pagination
	c.OnHTML("a.btn.btn-outline.BtnGroup-item", func(e *colly.HTMLElement) {
		if e.Text == "Next" {
			nextUrl = e.Attr("href")
		} else {
			nextUrl = ""
		}
	})

	// deps
	c.OnHTML("div.Box-row", func(e *colly.HTMLElement) {
		avatar := e.ChildAttr("img", "src")
		user := e.ChildAttrs("a", "href")[0][1:]
		repo := e.ChildText("a.text-bold")
		repoUrl := fmt.Sprintf("https://github.com/%s/%s", user, repo)

		deps = append(deps, Dep{User: user,
			Repo:    repo,
			Avatar:  avatar,
			RepoUrl: repoUrl,
		})
	})

	err := c.Visit(url)
	if err != nil {
		return "", deps
	}
	return nextUrl, deps
}

func InitScrape() tea.Cmd {
	return func() tea.Msg {
		return InitScrapeTick("Tick")
	}
}

func Scrape(url string, deps *[]Dep) tea.Cmd {
	return func() tea.Msg {
		nextUrl, deps := scrapePage(url, *deps)

		return PageTick{
			nextUrl: nextUrl,
			deps:    deps,
		}
	}
}
