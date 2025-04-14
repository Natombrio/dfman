package main

import (
	"errors"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

func relative_path_to_abs(given_path string) string {
    if strings.HasPrefix(given_path, "~") {
        home, err := os.UserHomeDir()
        if err != nil {
            log.Fatal("Error finding user home path.")
        }
        return filepath.Join(home, given_path[1:])
    }
    if !filepath.IsAbs(given_path) {
        log.Fatalf("Path should be absolute. Was %s", given_path)
    }
    return given_path
}

func run_command(workspace_dir string, binary_name string, args ...string) string {
    cmd := exec.Command(binary_name, args...)
    //fmt.Printf("\nRunning -- %s\n", cmd.String())
    cmd.Dir = relative_path_to_abs(workspace_dir)
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

func link_files() {
    fmt.Printf("TODO\n")
}

func is_symlinked(symlink_source string, symlink_target string) bool {
    real_target, err := os.Readlink(symlink_source)
    if err != nil {
        log.Println("YAY!")
    }
    if symlink_target != real_target {
        return false
    }
    return true
}

func check_if_paths_exist(items []SyncItem) {
    for idx := range items {
        _, srcErr := os.Stat(items[idx].source.path)
        _, dstErr := os.Stat(items[idx].destination.path)
        items[idx].source.exists = !errors.Is(srcErr, os.ErrNotExist)
        items[idx].destination.exists = !errors.Is(dstErr, os.ErrNotExist)
    }
}

func check_if_symlinked(items []SyncItem) {
    for idx := range items {
        info, err := os.Stat(items[idx].source.path)
        if err != nil {
            log.Fatalf("Error processing symlink %s, %s", items[idx].name, err)
        }
        linkSource := ""
        if info.IsDir() {
            linkSource = items[idx].destination.path+"/"+items[idx].name
        } else {
            linkSource = items[idx].destination.path+"/"+info.Name()
        }
        linkDestination := items[idx].source.path
        items[idx].linked = is_symlinked(linkSource, linkDestination)
        _, srcErr := os.Stat(items[idx].source.path)
        _, dstErr := os.Stat(items[idx].destination.path)
        items[idx].source.exists = !errors.Is(srcErr, os.ErrNotExist)
        items[idx].destination.exists = !errors.Is(dstErr, os.ErrNotExist)
    }
}
