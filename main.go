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
	state   int
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
	s.Spinner = spinner.Points
	s.Style = lipgloss.NewStyle().Foreground(lipgloss.Color("205"))

	return model{
		input:   ti,
		state:   0,
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
			if m.state == 0 {
				// todo: validate/clean input
				m.state = 1
				m.repo = m.input.Value()
				if m.repo == "" {
					hardcodeUrl := "https://github.com/aquasecurity/trivy/network/dependents"
					m.repo = hardcodeUrl
				}
				return m, tea.Batch(InitScrape(), m.spinner.Tick)
			}
		}

	case errMsg:
		SortByStars(&m.deps)
		WriteJson(m.deps)
		m.err = msg
		fmt.Printf("\nError: %v\nCrawled results are saved to deps.json\nQuitting...\n", msg)
		return m, tea.Quit

	case InitScrapeTick:
		return m, tea.Batch(Scrape(m.repo), m.spinner.Tick)

	case PageTick:
		nextUrl, deps := msg.nextUrl, msg.deps
		m.pages++
		m.deps = append(m.deps, deps...)
		m.count = len(m.deps)

		if nextUrl != "" {
			m.logs = fmt.Sprintf("📦 %s", nextUrl)
			return m, tea.Batch(Scrape(nextUrl), m.spinner.Tick)
		} else {
			SortByStars(&m.deps)
			WriteJson(m.deps)
			m.logs = "\nDone! 🧨 Write deps.json..."
			m.done = true
			return m, tea.Quit
		}
	}

	m.spinner, cmdSpinner = m.spinner.Update(msg)
	m.input, cmdInput = m.input.Update(msg)

	return m, tea.Batch(cmd, cmdSpinner, cmdInput)

}

func (m model) View() string {
	title := "🥦 Uby"
	body := ""
	footer := "(esc to quit)"

	switch m.state {
	case 0:
		body = m.input.View()
	case 1:
		if !m.done {
			body += m.spinner.View()
			body += fmt.Sprintf(" Fetching dependencies from %s...", m.repo)
		}
	}

	outText := fmt.Sprintf("%s\n\n%s\n\n%s\n\npages: %d\nrepos: %d\n%s\n", title, body, footer, m.pages, m.count, m.logs)
	return outText

}
