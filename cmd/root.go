/*
Copyright © 2023 Vistara Labs mayur@vistara.dev
*/
package cmd

import (
	"fmt"
	"log"
	"os"
	"vimana/cli"

	"github.com/common-nighthawk/go-figure"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use: "vimana",
	// Short: "  A Hardware Availability Network Orchestrator",
}

func Execute() {
	// err := rootCmd.Execute()
	// if err != nil {
	// 	os.Exit(1)
	// }
}

// Define function type for dependency injection
type userHomeDirFunc func() (string, error)
type getCommandsFromConfigFunc func(string, map[string]cli.NodeCommander) ([]*cobra.Command, error)

var OsUserHomeDir userHomeDirFunc = os.UserHomeDir
var GetCommandsFromConfig getCommandsFromConfigFunc = cli.GetCommandsFromConfig

var force bool
var noTrack bool

func InitCLI() error {
	home, err := OsUserHomeDir()
	if err != nil {
		log.Fatal(err)
	}
	configFile := home + "/.vimana/config.toml"
	rootCmd = &cobra.Command{Use: "vimana"}

	commands, err := GetCommandsFromConfig(configFile, CommanderRegistry)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error:", err)
		return err
	}
	rootCmd.AddCommand(commands...)
	rootCmd.AddCommand(initVimana())
	rootCmd.AddCommand(versionCommand())

	return rootCmd.Execute()
}

func init() {
	vimanaFig := figure.NewFigure("vimana", "", true)
	vimanaFig.Print()

	if err := InitCLI(); err != nil {
		log.Fatalf("Failed to initialize CLI: %s", err)
	}
}

func initVimana() *cobra.Command {
	initCmd := &cobra.Command{
		Use:   "init",
		Short: "Initializes and checks system resources",
		RunE: func(cmd *cobra.Command, args []string) error {
			return InitializeSystem(force, noTrack)
		},
	}
	initCmd.PersistentFlags().BoolVarP(&force, "force", "f", false, "Force initialization")
	initCmd.PersistentFlags().BoolVarP(&noTrack, "no-track", "n", false, "Opt out of anonymous usage tracking")
	return initCmd
}

func versionCommand() *cobra.Command {
	versionCmd := &cobra.Command{
		Use:   "version",
		Short: "Print the version of vimana",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("vimana version: ", Version)
		},
	}
	return versionCmd
}

func printASCIIArt() {
	art := `
___    ________________  __________ _____   _________ 
__ |  / /____  _/___   |/  /___    |___  | / /___    |
__ | / /  __  /  __  /|_/ / __  /| |__   |/ / __  /| |
__ |/ /  __/ /   _  /  / /  _  ___ |_  /|  /  _  ___ |
_____/   /___/   /_/  /_/   /_/  |_|/_/ |_/   /_/  |_|
`
	fmt.Println(art)
}
