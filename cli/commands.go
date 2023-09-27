package cli

import (
	"fmt"
	"log"
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

type NodeCommander interface {
	Init(*cobra.Command, []string, Mode) error
	Start(*cobra.Command, []string, Mode)
	Stop(*cobra.Command, []string, Mode)
	Status(*cobra.Command, []string, Mode)
}

type BaseCommander struct {
	Name         string
	NodeType     string
	componentMgr *components.ComponentManager
}

func (b *BaseCommander) initComponentManager(component config.ComponentType, binary string) {
	if b.componentMgr == nil {
		b.componentMgr = components.NewComponentManager(component, binary, b.NodeType)
	}
	// if b.componentMgr == nil && component == "avail" {
	// 	b.componentMgr = components.NewComponentManager("avail", binary, b.NodeType)
	// }
}

func GetCommandsFromConfig(filepath string, commanderRegistry map[string]NodeCommander) ([]*cobra.Command, error) {
	var config Config
	if _, err := toml.DecodeFile(filepath, &config); err != nil {
		return nil, err
	}

	var commands []*cobra.Command

	runCmd := &cobra.Command{
		Use:   "run",
		Short: "Run a modular component",
	}

	for component, nodeTypes := range config.Components {
		for nodeType := range nodeTypes {
			currentComponent := component
			currentNodeType := nodeType

			subCmd := &cobra.Command{
				Use:  fmt.Sprintf("%s-%s", currentComponent, currentNodeType),
				Args: cobra.NoArgs,
				Run: func(c *cobra.Command, args []string) {
					ntype := nodeTypes[currentNodeType]

					key := fmt.Sprintf("%s-%s", currentComponent, currentNodeType)
					commander := commanderRegistry[key]
					if commander != nil {
						commander.Start(c, args, ntype)
					} else {
						log.Fatalf("Components '%s' of type '%s' not recognized", component, ntype)
					}
				},
			}
			runCmd.AddCommand(subCmd)
		}
	}
	commands = append(commands, runCmd)
	return commands, nil
}

func GetCommonSubCommands() map[string]func(*cobra.Command, []string, Mode) {
	return map[string]func(*cobra.Command, []string, Mode){
		"init": func(cmd *cobra.Command, args []string, mode Mode) {
			fmt.Println("Command:", cmd.Name(), "Args:", args, "Mode:", mode, "Expected: Initing")
		},
		"start": func(cmd *cobra.Command, args []string, mode Mode) {
			fmt.Println("Command:", cmd.Name(), "Args:", args, "Mode:", mode, "Expected: Starting")
		},
		"stop": func(cmd *cobra.Command, args []string, mode Mode) {
			fmt.Println("Command:", cmd.Name(), "Args:", args, "Mode:", mode, "Expected: Stopping")
		},
		"status": func(cmd *cobra.Command, args []string, mode Mode) {
			fmt.Println("Command:", cmd.Name(), "Args:", args, "Mode:", mode, "Expected: Status")
		},
		"list": func(cmd *cobra.Command, args []string, mode Mode) {
			fmt.Println("Command:", cmd.Name(), "Args:", args, "Mode:", mode, "Expected: Listing")
		},
	}
}
