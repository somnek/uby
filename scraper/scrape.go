package main

import (
	"time"

	tea "github.com/charmbracelet/bubbletea"
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
