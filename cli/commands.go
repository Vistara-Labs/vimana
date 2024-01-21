package cli

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"vimana/cmd/utils"
	"vimana/components"
	"vimana/config"

	"github.com/BurntSushi/toml"
	"github.com/spf13/cobra"
)

type Config struct {
	Components map[string]Component `toml:"components"`
}

type Component map[string]Mode

type Mode struct {
	Binary   string `toml:"binary"`
	Download string `toml:"download"`
}

func WriteConf(conf Config) error {
	// open the file
	configFile := os.Getenv("HOME") + "/.vimana/config.toml"
	file, err := os.OpenFile(configFile, os.O_WRONLY|os.O_TRUNC, 0644)
	if err != nil {
		return err
	}

	defer file.Close()

	if err := toml.NewEncoder(file).Encode(conf); err != nil {
		return err
	}

	return nil
}

type NodeCommander interface {
	AddFlags(*cobra.Command)
	Init(*cobra.Command, []string, Mode, string) error
	Start(*cobra.Command, []string, Mode, string)
	Stop(*cobra.Command, []string, Mode, string)
	Status(*cobra.Command, []string, Mode, string)
}

type BaseCommander struct {
	Name         string
	NodeType     string
	componentMgr *components.ComponentManager
	config       *components.ComponentConfig
}

func (b *BaseCommander) initComponentManager(component config.ComponentType, binary string) {
	if b.componentMgr == nil {
		b.componentMgr = components.NewComponentManager(component, binary, b.NodeType, b.config)
	}
}

func GetCommandsFromConfig(path string, commanderRegistry map[string]NodeCommander) ([]*cobra.Command, error) {
	var config Config
	if _, err := toml.DecodeFile(path, &config); err != nil {
		return nil, err
	}

	// update commanderRegistry
	for component, nodeTypes := range config.Components {
		fmt.Println(component)
		for nodeType := range nodeTypes {
			fmt.Println(nodeType, nodeTypes[nodeType])
			fmt.Println(nodeTypes[nodeType].Binary)
			fmt.Println(nodeTypes[nodeType].Download)
			cmd_key := component + "-" + nodeType
			fmt.Println(cmd_key)

			commander := commanderRegistry[cmd_key]
			if commander == nil {
				commanderRegistry[cmd_key] = NewUniversalCommander(nodeType)
			}
		}
	}

	fmt.Println(len(commanderRegistry))

	//
	var commands []*cobra.Command

	runCmd := &cobra.Command{
		Use:   "run",
		Short: "Run a modular component",
	}
	initPath := filepath.Join(os.Getenv("HOME"), ".vimana", "init.toml")
	initConf, err := utils.LoadVimanaConfig(initPath)
	if err != nil {
		return nil, err
	}

	for component, nodeTypes := range config.Components {
		currentComponent := component
		componentCmd := &cobra.Command{
			Use:   currentComponent,
			Short: fmt.Sprintf("Run the %s component", currentComponent),
		}

		for nodeType := range nodeTypes {
			ntype := nodeTypes[nodeType]
			currentNodeType := nodeType
			nodeCmd := &cobra.Command{
				Use:  nodeType + "-node",
				Args: cobra.NoArgs,
				Run: func(c *cobra.Command, args []string) {
					key := fmt.Sprintf("%s-%s", currentComponent, currentNodeType)
					fmt.Println("commander component", key)

					commander := commanderRegistry[key]
					if commander != nil {
						initConf.SpaceCore = currentComponent
						if initConf.Analytics.Enabled {
							go utils.SaveAnalyticsData(initConf)
						}
						commander.Start(c, args, ntype, currentComponent)
					} else {
						log.Fatalf("Components '%s' of type '%s' not recognized", component, ntype)
					}
				},
			}
			if commander, ok := commanderRegistry[fmt.Sprintf("%s-%s", currentComponent, nodeType)]; ok {
				commander.AddFlags(nodeCmd)
			}
			componentCmd.AddCommand(nodeCmd)
		}
		runCmd.AddCommand(componentCmd)

	}
	commands = append(commands, runCmd)

	// start command
	startCmd := &cobra.Command{
		Use:   "start",
		Short: "start a modular component",
	}
	for component, nodeTypes := range config.Components {
		currentComponent := component
		componentCmd := &cobra.Command{
			Use:   currentComponent,
			Short: fmt.Sprintf("start the %s component", currentComponent),
		}

		for nodeType := range nodeTypes {
			ntype := nodeTypes[nodeType]
			currentNodeType := nodeType
			nodeCmd := &cobra.Command{
				Use:  nodeType + "-node",
				Args: cobra.NoArgs,
				Run: func(c *cobra.Command, args []string) {
					key := fmt.Sprintf("%s-%s", currentComponent, currentNodeType)
					fmt.Println("commander component", key)

					commander := commanderRegistry[key]
					if commander != nil {
						initConf.SpaceCore = currentComponent
						if initConf.Analytics.Enabled {
							go utils.SaveAnalyticsData(initConf)
						}
						commander.Start(c, args, ntype, key)
					} else {
						log.Fatalf("Components '%s' of type '%s' not recognized", component, ntype)
					}
				},
			}
			if commander, ok := commanderRegistry[fmt.Sprintf("%s-%s", currentComponent, nodeType)]; ok {
				commander.AddFlags(nodeCmd)
			}
			componentCmd.AddCommand(nodeCmd)
		}
		startCmd.AddCommand(componentCmd)

	}
	commands = append(commands, startCmd)

	// stop command
	stopCmd := &cobra.Command{
		Use:   "stop",
		Short: "stop a modular component",
	}
	for component, nodeTypes := range config.Components {
		currentComponent := component
		componentCmd := &cobra.Command{
			Use:   currentComponent,
			Short: fmt.Sprintf("stop the %s component", currentComponent),
		}

		for nodeType := range nodeTypes {
			ntype := nodeTypes[nodeType]
			currentNodeType := nodeType
			nodeCmd := &cobra.Command{
				Use:  nodeType + "-node",
				Args: cobra.NoArgs,
				Run: func(c *cobra.Command, args []string) {
					key := fmt.Sprintf("%s-%s", currentComponent, currentNodeType)
					fmt.Println("commander component", key)

					commander := commanderRegistry[key]
					if commander != nil {
						initConf.SpaceCore = currentComponent
						if initConf.Analytics.Enabled {
							go utils.SaveAnalyticsData(initConf)
						}
						commander.Stop(c, args, ntype, key)
					} else {
						log.Fatalf("Components '%s' of type '%s' not recognized", component, ntype)
					}
				},
			}
			if commander, ok := commanderRegistry[fmt.Sprintf("%s-%s", currentComponent, nodeType)]; ok {
				commander.AddFlags(nodeCmd)
			}
			componentCmd.AddCommand(nodeCmd)
		}
		stopCmd.AddCommand(componentCmd)

	}
	commands = append(commands, stopCmd)

	// status command
	statusCmd := &cobra.Command{
		Use:   "status",
		Short: "show status of a modular component",
	}
	for component, nodeTypes := range config.Components {
		currentComponent := component
		componentCmd := &cobra.Command{
			Use:   currentComponent,
			Short: fmt.Sprintf("show status of the %s component", currentComponent),
		}

		for nodeType := range nodeTypes {
			ntype := nodeTypes[nodeType]
			currentNodeType := nodeType
			nodeCmd := &cobra.Command{
				Use:  nodeType + "-node",
				Args: cobra.NoArgs,
				Run: func(c *cobra.Command, args []string) {
					key := fmt.Sprintf("%s-%s", currentComponent, currentNodeType)
					fmt.Println("commander component", key)

					commander := commanderRegistry[key]
					if commander != nil {
						initConf.SpaceCore = currentComponent
						if initConf.Analytics.Enabled {
							go utils.SaveAnalyticsData(initConf)
						}
						commander.Status(c, args, ntype, key)
					} else {
						log.Fatalf("Components '%s' of type '%s' not recognized", component, ntype)
					}
				},
			}
			if commander, ok := commanderRegistry[fmt.Sprintf("%s-%s", currentComponent, nodeType)]; ok {
				commander.AddFlags(nodeCmd)
			}
			componentCmd.AddCommand(nodeCmd)
		}
		statusCmd.AddCommand(componentCmd)

	}
	commands = append(commands, statusCmd)

	return commands, nil
}
