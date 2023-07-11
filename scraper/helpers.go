package main

import (
	"encoding/json"
	"log"
	"os"
	"sort"
)

func WriteJson(deps []Dep) {
	jsonData, err := json.MarshalIndent(deps, "", "  ")
	if err != nil {
		log.Fatal(err)
	}

	err = os.WriteFile("deps.json", jsonData, 0644)
	if err != nil {
		log.Fatal(err)
	}
}

func SortByStars(deps *[]Dep) {
	sort.Slice(*deps, func(i, j int) bool {
		return (*deps)[i].Stars > (*deps)[j].Stars
	})
}
