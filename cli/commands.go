package cli

import (
	"fmt"
	"log"
	"vimana/components"

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

func (b *BaseCommander) initComponentManager(binary string) {
	if b.componentMgr == nil {
		b.componentMgr = components.NewComponentManager("celestia", binary, b.NodeType)
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
		Short: "Run a component",
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
						log.Fatalf("Component '%s' of type '%s' not recognized", component, nodeType)
					}
				},
			}
			runCmd.AddCommand(subCmd)
		}
	}
	commands = append(commands, runCmd)

	// for componentName, modes := range config.Components {
	// 	componentCmd := &cobra.Command{Use: componentName}
	// 	commands = append(commands, componentCmd)

	// 	for modeName, modeData := range modes {
	// 		modeCmd := &cobra.Command{Use: modeName}
	// 		componentCmd.AddCommand(modeCmd)
	// 		key := componentName + "-" + modeName
	// 		commander := commanderRegistry[key]
	// 		mData := modeData

	// 		if commander != nil {
	// 			modeCmd.AddCommand(&cobra.Command{
	// 				Use: "init",
	// 				Run: func(c *cobra.Command, args []string) {
	// 					commander.Init(c, args, mData)
	// 				},
	// 			})

	// 			modeCmd.AddCommand(&cobra.Command{
	// 				Use: "start",
	// 				Run: func(c *cobra.Command, args []string) {
	// 					commander.Start(c, args, mData)
	// 				},
	// 			})

	// 			modeCmd.AddCommand(&cobra.Command{
	// 				Use: "stop",
	// 				Run: func(c *cobra.Command, args []string) {
	// 					commander.Stop(c, args, mData)
	// 				},
	// 			})

	// 			modeCmd.AddCommand(&cobra.Command{
	// 				Use: "status",
	// 				Run: func(c *cobra.Command, args []string) {
	// 					commander.Status(c, args, mData)
	// 				},
	// 			})

	// 			// modeCmd.AddCommand(&cobra.Command{
	// 			// 	Use: "list",
	// 			// 	Run: func(c *cobra.Command, args []string) {
	// 			// 		commander.List(c, args, modeData)
	// 			// 	},
	// 			// })
	// 		}

	// 		// for cmdName, cmdFunc := range commonSubCommands {
	// 		// 	sc := &cobra.Command{
	// 		// 		Use: cmdName,
	// 		// 		Run: makeRunFunc(cmdFunc, modeData), // <-- Here's the change, we use makeRunFunc to generate the Run function
	// 		// 	}
	// 		// 	modeCmd.AddCommand(sc)
	// 		// }

	// 	}
	// }
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
