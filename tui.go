package main

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
)

type ModelState struct {
    Github_Url string
    Dotfile_Dir string
    syncItems []SyncItem
}

type SyncItem struct {
    name string
    source PathItem
    destination PathItem
    linked bool
}

type PathItem struct {
    path string
    exists bool
}

type model struct {
    actions []string
    dotfiles []SyncItem
    cursor int
    actionStatus map[int]string
    modelState ModelState
}

func initialModel(config ConfigData) model {
    var items []SyncItem
    for linkName, val := range config.Links {
        items = append(items, SyncItem{
            name: linkName,
            source: PathItem{
                path: home_relative_path_to_abs(val.Source),
                exists: false,
            },
            destination: PathItem{
                path: home_relative_path_to_abs(val.Destination),
                exists: false,
            },
            linked: false,
        })
    }
    modelState := ModelState{
        Github_Url: config.Github_Url,
        Dotfile_Dir: config.Dotfile_Dir,
        syncItems: items,
    }
    dotfiles := modelState.syncItems
    dotfiles = check_if_paths_exist(dotfiles)

    return model{
        actions: []string{"Clone", "Pull", "Link"},
        dotfiles: dotfiles,
        cursor: 0,
        actionStatus: make(map[int]string),
        modelState: modelState,
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
            for i := range m.actionStatus {
                m.actionStatus[i] = ""
            }
        case "down", "j":
            if m.cursor < len(m.dotfiles)-1 {
                m.cursor++
            }
            for i := range m.actionStatus {
                m.actionStatus[i] = ""
            }
        case "enter", " ":
            switch m.actions[m.cursor] {
            case "Clone":
                res := strings.Split(git_clone("~", m.modelState.Github_Url, m.modelState.Dotfile_Dir), "\n")[0]
                m.actionStatus[m.cursor] = res
            case "Pull":
                res := strings.Split(git_pull(m.modelState.Dotfile_Dir), "\n")[0]
                m.actionStatus[m.cursor] = res
            case "Link":
                link_files()
            }
        }
    }
    return m, nil
}


func (m model) View() string {
    s := fmt.Sprintf("Github repo:\n %s\n", m.modelState.Github_Url)
    s += fmt.Sprintf("Dotfile dir:\n %s\n\n", m.modelState.Dotfile_Dir)
    s += fmt.Sprintf("Tracked config items:\n")
    items := m.modelState.syncItems
    for _, item := range items {
        s += fmt.Sprintf(" %s, linked: %t\n", item.name, item.linked)
        s += fmt.Sprintf("  exists: %t   src: %s  \n", item.source.exists, item.source.path)
        s += fmt.Sprintf("  exists: %t   dst: %s  \n", item.destination.exists, item.destination.path)
    }

    s += "\nChoose an action to perform\n\n"
    for i, choice := range m.actions {
        cursor := " "
        if m.cursor == i {
            cursor = ">"
        }
        status := m.actionStatus[i]
        if status != "" {
            status = "|| " + m.actionStatus[i]
        }
        s += fmt.Sprintf("%s %s %s\n", cursor, choice, status)
    }
    s += "\nPress q to quit.\n"
    return s
}

