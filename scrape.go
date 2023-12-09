package main

import (
	"fmt"
	"strconv"
	"strings"

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

type PageTick struct {
	nextUrl string
	deps    []Dep
}
type InitScrapeTick string

func ToNum(s string) int {
	// remove commas and spaces
	s = strings.ReplaceAll(s, ",", "")
	s = strings.ReplaceAll(s, " ", "")

	// convert to int
	i, _ := strconv.Atoi(s)
	return i
}

func extractStars(e *colly.HTMLElement) int {
	parent := e.ChildText("span.color-fg-muted.text-bold.pl-3")
	split := strings.TrimSpace(strings.Split(parent, " ")[0])
	stars := ToNum(split)
	return stars
}

type ScrapeModel struct {
	nextUrl string
	deps    []Dep
}

func scrapePage(url string) (ScrapeModel, error) {
	c := colly.NewCollector()
	var nextUrl string

	// pagination
	c.OnHTML("a.btn.BtnGroup-item", func(e *colly.HTMLElement) {
		if e.Text == "Next" {
			nextUrl = e.Attr("href")
		} else {
			nextUrl = ""
		}
	})

	// deps
	var deps []Dep
	c.OnHTML("div.Box-row", func(e *colly.HTMLElement) {
		avatar := e.ChildAttr("img", "src")
		user := e.ChildAttrs("a", "href")[0][1:]
		repo := e.ChildText("a.text-bold")
		stars := extractStars(e)
		repoUrl := fmt.Sprintf("https://github.com/%s/%s", user, repo)

		deps = append(deps, Dep{User: user,
			Repo:    repo,
			Stars:   stars,
			Avatar:  avatar,
			RepoUrl: repoUrl,
		})
	})

	err := c.Visit(url)
	if err != nil {
		return ScrapeModel{}, err
	}
	return ScrapeModel{nextUrl: nextUrl, deps: deps}, nil
}

func InitScrape() tea.Cmd {
	return func() tea.Msg {
		return InitScrapeTick("Tick")
	}
}

func Scrape(url string) tea.Cmd {
	return func() tea.Msg {
		scrapeResult, err := scrapePage(url)
		if err != nil {
			return errMsg(err)
		}
		nextUrl := scrapeResult.nextUrl
		newDeps := scrapeResult.deps

		return PageTick{
			nextUrl: nextUrl,
			deps:    newDeps,
		}
	}
}
