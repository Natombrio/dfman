package main

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
)

type model struct {
    choices []string
    cursor int
    status map[int]bool
    data Config
}

func initialModel(config Config) model {
    keys := make([]string, 0, len(config.Links))
    for k := range config.Links {
        keys = append(keys, k)
    }

    return model{
        choices: keys,
        cursor: 0,
        status: make(map[int]bool),
        data: config,
    }
}

func (m model) Init() tea.Cmd {
    return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
    switch msg := msg.(type) {
    case tea.KeyMsg:
        switch msg.String() {
        case "ctrl+c", "q":
            return m, tea.Quit
        case "up", "k":
            if m.cursor > 0 {
                m.cursor--
            }
        case "down", "j":
            if m.cursor < len(m.choices)-1 {
                m.cursor++
            }
        case "enter", " ":
            m.status[m.cursor] = !m.status[m.cursor]
        }
    }
    return m, nil
}


func (m model) View() string {
    s := "Choose an action to perform\n\n"
    for i, choice := range m.choices {
        cursor := " "
        if m.cursor == i {
            cursor = ">"
        }
        status := ":("
        if m.status[i] {
            status = ":)"
        }
        s += fmt.Sprintf("%s %s \t|| %s\n", cursor, choice, status)
    }
    s += "\nPress q to quit.\n"
    return s
}

