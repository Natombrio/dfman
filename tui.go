package main

import (
	"fmt"
	"slices"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
)

type model struct {
    actions []string
    dotfiles []string
    cursor int
    status map[int]string
    data Config
}

func sortMap(config_map map[string]DirectoryLink) []string {
    keys := make([]string, 0, len(config_map))
    for k := range config_map {
        keys = append(keys, k)
    }
    slices.Sort(keys)
    return keys
}

func initialModel(config Config) model {
    config.Dotfile_Dir = home_relative_path_to_abs(config.Dotfile_Dir)
    for k, val := range config.Links {
        config.Links[k] = DirectoryLink{
            Source: home_relative_path_to_abs(val.Source),
            Destination: home_relative_path_to_abs(val.Destination),
        }
    }
    keys := sortMap(config.Links)

    return model{
        actions: []string{"Clone", "Pull", "Link"},
        dotfiles: keys,
        cursor: 0,
        status: make(map[int]string),
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
            for i := range m.status {
                m.status[i] = ""
            }
        case "down", "j":
            if m.cursor < len(m.dotfiles)-1 {
                m.cursor++
            }
            for i := range m.status {
                m.status[i] = ""
            }
        case "enter", " ":
            switch m.actions[m.cursor] {
            case "Clone":
                res := strings.Split(git_clone("~", m.data.Github_Url, m.data.Dotfile_Dir), "\n")[0]
                m.status[m.cursor] = res
            case "Pull":
                res := strings.Split(git_pull(m.data.Dotfile_Dir), "\n")[0]
                m.status[m.cursor] = res
            case "Link":
                link_files()
            }
        }
    }
    return m, nil
}


func (m model) View() string {
    s := fmt.Sprintf("Github repo:\n %s\n", m.data.Github_Url)
    s += fmt.Sprintf("Dotfile dir:\n %s\n\n", m.data.Dotfile_Dir)
    s += fmt.Sprintf("Tracked config files:\n")
    files := sortMap(m.data.Links)
    for _, name := range files {
        s += fmt.Sprintf(" %s\n", name)
    }

    s += "\nChoose an action to perform\n\n"
    for i, choice := range m.actions {
        cursor := " "
        if m.cursor == i {
            cursor = ">"
        }
        status := m.status[i]
        if status != "" {
            status = "|| " + m.status[i]
        }
        s += fmt.Sprintf("%s %s %s\n", cursor, choice, status)
    }
    s += "\nPress q to quit.\n"
    return s
}

