package main

import (
	"fmt"
	"log"

	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

func main() {
	p := tea.NewProgram(initialModel())
	if _, err := p.Run(); err != nil {
		log.Fatal(err)
	}
}

type model struct {
	tab     int
	repo    string
	input   textinput.Model
	deps    []Dep
	err     error
	spinner spinner.Model
	done    bool
	pages   int
	logs    string
	count   int
}

type (
	errMsg error
)

func initialModel() model {
	ti := textinput.New()
	ti.Placeholder = "Enter a dependency graph URL..."
	ti.Focus()
	ti.CharLimit = 200
	ti.Width = 200

	s := spinner.New()
	s.Spinner = spinner.Dot
	s.Style = lipgloss.NewStyle().Foreground(lipgloss.Color("205"))

	return model{
		input:   ti,
		tab:     0,
		err:     nil,
		spinner: s,
		pages:   0,
		deps:    []Dep{},
	}
}

func (m model) Init() tea.Cmd {
	return tea.Batch(textinput.Blink)
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	var cmdSpinner tea.Cmd
	var cmdInput tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {

		case tea.KeyCtrlC, tea.KeyEsc:
			return m, tea.Quit

		case tea.KeyEnter:
			if m.tab == 0 {
				// todo: validate/clean input
				m.tab = 1
				m.repo = m.input.Value()
				// <-- ️ HARD CODED URL FOR TESTING 🚨
				hardcodeUrl := "https://github.com/aquasecurity/trivy/network/dependents"
				m.repo = hardcodeUrl
				// 🚨️ HARD CODED URL FOR TESTING -->
				return m, tea.Batch(InitScrape(), m.spinner.Tick)
			}

		}

	case errMsg:
		m.err = msg
		return m, nil

	case InitScrapeTick:
		return m, tea.Batch(Scrape(m.repo, &m.deps), m.spinner.Tick)

	case PageTick:
		nextUrl, deps := msg.nextUrl, msg.deps
		m.count += len(deps)
		if nextUrl != "" {
			// m.logs += concatDeps(deps)
			m.pages++
			return m, tea.Batch(Scrape(nextUrl, &m.deps), m.spinner.Tick)
		} else {
			m.logs += "🧨"
			return m, tea.Quit
		}
	}

	m.spinner, cmdSpinner = m.spinner.Update(msg)
	m.input, cmdInput = m.input.Update(msg)

	return m, tea.Batch(cmd, cmdSpinner, cmdInput)

}

func concatDeps(deps []Dep) string {
	var out string
	for _, dep := range deps {
		out += fmt.Sprintf("📦 %s\n", dep.RepoUrl)
	}
	return out
}

func (m model) View() string {
	title := "🥦 Uby"
	body := ""
	footer := "(esc to quit)"

	switch m.tab {
	case 0:
		body = m.input.View()
	case 1:
		if !m.done {
			body += m.spinner.View()
			body += fmt.Sprintf("Fetching dependencies from %s...", m.repo)
		} else {
			body = "Done!"
		}
	}

	outText := fmt.Sprintf("%s\n\n%s\n\n%s\n\npages: %d\n%d\n%s\n", title, body, footer, m.pages, m.count, m.logs)
	return outText

}
