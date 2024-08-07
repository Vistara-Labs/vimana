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

	logger "vimana/log"
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

	// configure logging
	logger.Configure(&logger.Config{
		Verbosity: logger.LogVerbosityInfo,
		Format:    logger.LogFormatText,
		Output:    "stderr",
	})

	rootCmd = &cobra.Command{
		Use:   "vimana",
		Short: "Orchestration client for running spacecores",
		Long:  "Vimana is an orchestration client and a cli for managing spacecore plugins.",
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) == 0 {
				vimanaFig := figure.NewFigure("vimana", "isometric1", true)
				vimanaFig.Print()
				cmd.Help()
			}
		},
	}

	commands, err := GetCommandsFromConfig(configFile, CommanderRegistry)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error:", err)
		return err
	}

	rootCmd.AddCommand(initVimana())
	rootCmd.AddCommand(commands...)
	rootCmd.AddCommand(ScaffoldNew)
	rootCmd.AddCommand(pluginCommand())
	// rootCmd.AddCommand(agentCommand())
	rootCmd.AddCommand(versionCommand())
	// rootCmd.AddCommand(migrateCommand())

	rootCmd.AddCommand(registryCommands()...)

	logger.AddFlagsToCommand(rootCmd, &logger.Config{})

	rootCmd.Aliases = []string{"v"}

	rootCmd.Example = `
	  vimana scaffold fancy-ai-agent
	  vimana register <spacecore> <plugin_path>
	  vimana search <spacecore>
	  vimana get <plugin_name>
	  vimana plugin <plugin_path> start
	  vimana plugin <plugin_path> status
	  vimana plugin <plugin_path> stop
	  vimana plugin <plugin_path> logs
	`
	return rootCmd.Execute()
}

func init() {
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
			fmt.Printf("vimana version: %s\n", Version)
		},
	}
	return versionCmd
}
