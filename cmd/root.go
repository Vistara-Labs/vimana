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
	// Long:  `CLI to create and manage nodes on the Vistara Network.`,
	// Run: func(cmd *cobra.Command, args []string) {
	// 	fmt.Println("vimana: A Hardware Availability Network Manager")
	// },
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	home, err := os.UserHomeDir()
	if err != nil {
		log.Fatal(err)
	}
	configFile := home + "/.vimana/config.toml"
	rootCmd := &cobra.Command{Use: "vimana"}
	vimanaFig := figure.NewFigure("vimana", "", true)
	vimanaFig.Print()

	commands, err := cli.GetCommandsFromConfig(configFile, CommanderRegistry)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error:", err)
		return
	}

	rootCmd.AddCommand(commands...)
	rootCmd.AddCommand(versionCommand())

	if err := rootCmd.Execute(); err != nil {
		log.Fatalf("Failed to execute command: %s", err)
	}
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
