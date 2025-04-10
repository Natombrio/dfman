package main

import (
	"context"
	"log"
	"os"

	"github.com/BurntSushi/toml"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/urfave/cli/v3"
)

type Config struct {
	Github_Url string
	Links      map[string]DirectoryLink
}

type DirectoryLink struct {
	Source      string
	Destination string
}

func launch_tui(config Config) {
	p := tea.NewProgram(initialModel(config))
	if _, err := p.Run(); err != nil {
		log.Printf("Error: %v", err)
		os.Exit(1)
	}
}

func parse_config(config_file string) Config {
	f, err := os.ReadFile(config_file)
	if err != nil {
		log.Fatal("Error reading config file!\n\t", err)
	}
	config_data := string(f)
	var config Config
	_, err = toml.Decode(config_data, &config)
	if err != nil {
		log.Fatal("Error decoding config file: ", err)
	}
    return config

}

func main() {
	cmd := &cli.Command{
		Flags: []cli.Flag{
			&cli.BoolFlag{
				Name:  "cli",
				Value: false,
				Usage: "Set false to run as cli",
			},
			&cli.StringFlag{
				Name:  "config_file",
				Value: "config.toml",
				Usage: "Name of toml file containing config.",
			},
		},
		Action: func(ctx context.Context, cmd *cli.Command) error {
            cli_flag := cmd.Bool("cli")
            config_file := cmd.String("config_file")
            config := parse_config(config_file)

            if cli_flag {
                log.Println("cli mode")
            } else {
                launch_tui(config)
            }
            return nil
		},
	}
	cmd.Run(context.Background(), os.Args)
}
