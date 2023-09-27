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
		currentComponent := component
		componentCmd := &cobra.Command{
			Use:   currentComponent,
			Short: fmt.Sprintf("Run the %s component", component),
		}

		for nodeType := range nodeTypes {
			ntype := nodeTypes[nodeType]
			currentNodeType := nodeType
			nodeCmd := &cobra.Command{
				Use:  nodeType + "-node",
				Args: cobra.NoArgs,
				Run: func(c *cobra.Command, args []string) {
					key := fmt.Sprintf("%s-%s", currentComponent, currentNodeType)
					commander := commanderRegistry[key]
					if commander != nil {
						commander.Start(c, args, ntype)
					} else {
						log.Fatalf("Components '%s' of type '%s' not recognized", component, ntype)
					}
				},
			}
			componentCmd.AddCommand(nodeCmd)
		}
		runCmd.AddCommand(componentCmd)

	}
	commands = append(commands, runCmd)
	return commands, nil
}
