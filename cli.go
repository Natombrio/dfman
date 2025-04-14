package main

import (
	"fmt"
	"log"
    "os"
	"os/exec"
	"os/user"
	"strings"
)

func home_relative_path_to_abs(given_path string) string {
    path_parts := strings.Split(given_path, "~")
    if len(path_parts) > 2 {
        log.Fatal("Error: too many '~' in path")
    }
    usr, err := user.Current()
    if err != nil {
        log.Fatal("Error converting path: %s\n", err)
    }
    abs_path := strings.Replace(given_path, "~", usr.HomeDir, 1)
    return abs_path
}

func run_command(workspace_dir string, binary_name string, args ...string) string {
    cmd := exec.Command(binary_name, args...)
    //fmt.Printf("\nRunning -- %s\n", cmd.String())
    cmd.Dir = home_relative_path_to_abs(workspace_dir)
    output, err := cmd.CombinedOutput()
    if err != nil {
        return fmt.Sprintf("Error running command: %s\n", err)
    }
    return string(output)
}

func git_clone(workspace_dir string, git_repository_url string, dotfile_dir string) string {
    res := run_command(workspace_dir, "git", "clone", git_repository_url, dotfile_dir)
    return fmt.Sprintf("Done -- %s\n", res)
}

func git_pull(dotfile_dir string) string {
    res := run_command(dotfile_dir, "git", "pull")
    return fmt.Sprintf("Done -- %s\n", res)
}

func git_push() {
    fmt.Printf("TODO\n")
}

func link_files() {
    fmt.Printf("TODO\n")
}

func check_if_paths_exist(items []SyncItem) []SyncItem {
    results := make([]SyncItem, len(items))
    for _, item := range items {
        _, srcErr := os.Stat(item.source.path)
        _, dstErr := os.Stat(item.destination.path)
        item.source.exists = !os.IsNotExist(srcErr)
        item.destination.exists = !os.IsNotExist(dstErr)
        results = append(results, item)

    }
    return results
}

func validate_dotfile_destination() {
    // Check if the destination locations already exist
}
