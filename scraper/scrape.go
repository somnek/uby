package main

import "time"

type Dep struct {
	User    string `json:"user"`
	Repo    string `json:"repo"`
	Stars   int    `json:"stars"`
	Avatar  string `json:"avatar"`
	RepoUrl string `json:"repoUrl"`
	DepUrl  string `json:"depUrl"`
}

func SomeLongTask() {
	// ...
	time.Sleep(2 * time.Second)
}
