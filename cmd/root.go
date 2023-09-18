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

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "vimana",
	Short: "A Hardware Availability Network Orchestrator",
	// Long:  `CLI to create and manage nodes on the Vistara Network.`,
	// Run: func(cmd *cobra.Command, args []string) {
	// 	fmt.Println("vimana: A Hardware Availability Network Manager")
	// },
}

func Execute() {
	vimanaFig := figure.NewFigure("vimana", "", true)
	vimanaFig.Print()
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	// var configFile string = "$HOME/.vimana/config.toml"
	home, err := os.UserHomeDir()
	if err != nil {
		log.Fatal(err)
	}
	configFile := home + "/.vimana/config.toml"
	rootCmd := &cobra.Command{Use: "vimana"}

	// rootCmd.PersistentFlags().StringVar(&configFile, "config", "$HOME/.vimana/.config.toml", "config (default is $HOME/.vimana/.config.toml)")

	commands, err := cli.GetCommandsFromConfig(configFile, CommanderRegistry)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error:", err)
		return
	}

	for _, cmd := range commands {
		rootCmd.AddCommand(cmd)
	}
	rootCmd.AddCommand(versionCommand())
	// fmt.Print(rootCmd.UsageString())

	if err := rootCmd.Execute(); err != nil {
		log.Fatalf("Failed to execute command: %s", err)
	}
}

func versionCommand() *cobra.Command {
	versionCmd := &cobra.Command{
		Use:   "version",
		Short: "Print the version of vimana",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("vimana version", Version)
		},
	}
	return versionCmd
}
