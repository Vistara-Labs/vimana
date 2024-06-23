package cli

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"vimana/cmd/utils"
	"vimana/config"
	"vimana/log"
	"vimana/spacecores"

	"github.com/BurntSushi/toml"
	"github.com/spf13/cobra"
)

type Config struct {
	Spacecores map[string]Spacecore `toml:"spacecores"`
}

type Spacecore map[string]Mode

type Mode struct {
	Binary   string `toml:"binary"`
	Download string `toml:"download"`
	Install  string `toml:"install"`
	Start    string `toml:"start"`
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
	Logs(*cobra.Command, []string, Mode, string)
}

type BaseCommander struct {
	Name          string
	NodeType      string
	spacecoresMgr *spacecores.SpacecoreManager
	config        *spacecores.SpacecoreConfig
}

func (b *BaseCommander) initSpacecoreManager(spacecore config.SpacecoreType, binary string) {
	if b.spacecoresMgr == nil {
		b.spacecoresMgr = spacecores.NewSpacecoreManager(spacecore, binary, b.NodeType, b.config)
	}
}

func GetCommandsFromConfig(path string, commanderRegistry map[string]NodeCommander) ([]*cobra.Command, error) {
	var config Config
	if _, err := toml.DecodeFile(path, &config); err != nil {
		return nil, err
	}

	// update commanderRegistry
	for spacecore, nodeTypes := range config.Spacecores {
		for nodeType := range nodeTypes {
			cmd_key := spacecore + "-" + nodeType
			commander := commanderRegistry[cmd_key]
			if commander == nil {
				commanderRegistry[cmd_key] = NewUniversalCommander(nodeType)
			}
		}
	}

	var commands []*cobra.Command

	runCmd := &cobra.Command{
		Use:   "run",
		Short: "Run a spacecore",
	}
	initPath := filepath.Join(os.Getenv("HOME"), ".vimana", "init.toml")
	initConf, err := utils.LoadVimanaConfig(initPath)
	if err != nil {
		return nil, err
	}

	for spacecore, nodeTypes := range config.Spacecores {
		currentSpacecore := spacecore
		spacecoreCmd := &cobra.Command{
			Use:   currentSpacecore,
			Short: fmt.Sprintf("Run the %s spacecore", currentSpacecore),
		}

		for nodeType := range nodeTypes {
			ntype := nodeTypes[nodeType]
			currentNodeType := nodeType
			nodeCmd := &cobra.Command{
				Use:  nodeType + "-node",
				Args: cobra.NoArgs,
				Run: func(c *cobra.Command, args []string) {
					logger := log.GetLogger(c.Context())
					key := fmt.Sprintf("%s-%s", currentSpacecore, currentNodeType)
					logger.Info("commander spacecore", key)

					commander := commanderRegistry[key]
					if commander != nil {
						initConf.SpaceCore = currentSpacecore
						if initConf.Analytics.Enabled {
							// // go utils.SaveAnalyticsData(initConf)
						}
						commander.Start(c, args, ntype, key)
					} else {
						logger.Fatalf("Spacecores '%s' of type '%s' not recognized", spacecore, ntype)
					}
				},
			}
			if commander, ok := commanderRegistry[fmt.Sprintf("%s-%s", currentSpacecore, nodeType)]; ok {
				commander.AddFlags(nodeCmd)
			}
			spacecoreCmd.AddCommand(nodeCmd)
		}
		runCmd.AddCommand(spacecoreCmd)

	}
	commands = append(commands, runCmd)

	// start command
	startCmd := &cobra.Command{
		Use:   "start",
		Short: "Start a Spacecore",
	}
	for spacecore, nodeTypes := range config.Spacecores {
		currentSpacecore := spacecore
		spacecoreCmd := &cobra.Command{
			Use:   currentSpacecore,
			Short: fmt.Sprintf("start the %s spacecore", currentSpacecore),
		}

		for nodeType := range nodeTypes {
			ntype := nodeTypes[nodeType]
			currentNodeType := nodeType
			nodeCmd := &cobra.Command{
				Use:  nodeType + "-node",
				Args: cobra.NoArgs,
				Run: func(c *cobra.Command, args []string) {
					logger := log.GetLogger(c.Context())
					key := fmt.Sprintf("%s-%s", currentSpacecore, currentNodeType)
					logger.Info("commander spacecore", key)

					commander := commanderRegistry[key]
					if commander != nil {
						initConf.SpaceCore = currentSpacecore
						if initConf.Analytics.Enabled {
							// go utils.SaveAnalyticsData(initConf)
						}
						commander.Start(c, args, ntype, key)
					} else {
						logger.Fatalf("Spacecores '%s' of type '%s' not recognized", spacecore, ntype)
					}
				},
			}
			if commander, ok := commanderRegistry[fmt.Sprintf("%s-%s", currentSpacecore, nodeType)]; ok {
				commander.AddFlags(nodeCmd)
			}
			spacecoreCmd.AddCommand(nodeCmd)
		}
		startCmd.AddCommand(spacecoreCmd)

	}
	commands = append(commands, startCmd)

	// stop command
	stopCmd := &cobra.Command{
		Use:   "stop",
		Short: "Stop a Spacecore",
	}
	for spacecore, nodeTypes := range config.Spacecores {
		currentSpacecore := spacecore
		spacecoreCmd := &cobra.Command{
			Use:   currentSpacecore,
			Short: fmt.Sprintf("stop the %s spacecore", currentSpacecore),
		}

		for nodeType := range nodeTypes {
			ntype := nodeTypes[nodeType]
			currentNodeType := nodeType
			nodeCmd := &cobra.Command{
				Use:  nodeType + "-node",
				Args: cobra.NoArgs,
				Run: func(c *cobra.Command, args []string) {
					logger := log.GetLogger(c.Context())
					key := fmt.Sprintf("%s-%s", currentSpacecore, currentNodeType)
					logger.Info("commander spacecore", key)

					commander := commanderRegistry[key]
					if commander != nil {
						initConf.SpaceCore = currentSpacecore
						if initConf.Analytics.Enabled {
							// go utils.SaveAnalyticsData(initConf)
						}
						commander.Stop(c, args, ntype, key)
					} else {
						logger.Fatalf("Spacecores '%s' of type '%s' not recognized", spacecore, ntype)
					}
				},
			}
			if commander, ok := commanderRegistry[fmt.Sprintf("%s-%s", currentSpacecore, nodeType)]; ok {
				commander.AddFlags(nodeCmd)
			}
			spacecoreCmd.AddCommand(nodeCmd)
		}
		stopCmd.AddCommand(spacecoreCmd)

	}
	commands = append(commands, stopCmd)

	// status command
	statusCmd := &cobra.Command{
		Use:   "status",
		Short: "Show status of a Spacecore",
	}
	for spacecore, nodeTypes := range config.Spacecores {
		currentSpacecore := spacecore
		spacecoreCmd := &cobra.Command{
			Use:   currentSpacecore,
			Short: fmt.Sprintf("show status of the %s spacecore", currentSpacecore),
		}

		for nodeType := range nodeTypes {
			ntype := nodeTypes[nodeType]
			currentNodeType := nodeType
			nodeCmd := &cobra.Command{
				Use:  nodeType + "-node",
				Args: cobra.NoArgs,
				Run: func(c *cobra.Command, args []string) {
					logger := log.GetLogger(c.Context())
					key := fmt.Sprintf("%s-%s", currentSpacecore, currentNodeType)
					logger.Info("commander spacecore", key)

					commander := commanderRegistry[key]
					if commander != nil {
						initConf.SpaceCore = currentSpacecore
						if initConf.Analytics.Enabled {
							// go utils.SaveAnalyticsData(initConf)
						}
						commander.Status(c, args, ntype, key)
					} else {
						logger.Fatalf("Spacecores '%s' of type '%s' not recognized", spacecore, ntype)
					}
				},
			}
			if commander, ok := commanderRegistry[fmt.Sprintf("%s-%s", currentSpacecore, nodeType)]; ok {
				commander.AddFlags(nodeCmd)
			}
			spacecoreCmd.AddCommand(nodeCmd)
		}
		statusCmd.AddCommand(spacecoreCmd)

	}
	commands = append(commands, statusCmd)

	// logs command
	logsCmd := &cobra.Command{
		Use:   "logs",
		Short: "Show logs of running Spacecore",
	}
	for spacecore, nodeTypes := range config.Spacecores {
		currentSpacecore := spacecore
		spacecoreCmd := &cobra.Command{
			Use:   currentSpacecore,
			Short: fmt.Sprintf("show logs of the %s spacecore", currentSpacecore),
		}

		for nodeType := range nodeTypes {
			ntype := nodeTypes[nodeType]
			currentNodeType := nodeType
			nodeCmd := &cobra.Command{
				Use:  nodeType + "-node",
				Args: cobra.NoArgs,
				Run: func(c *cobra.Command, args []string) {
					logger := log.GetLogger(c.Context())
					key := fmt.Sprintf("%s-%s", currentSpacecore, currentNodeType)
					logger.Info("commander spacecore", key)

					commander := commanderRegistry[key]
					if commander != nil {
						initConf.SpaceCore = currentSpacecore
						if initConf.Analytics.Enabled {
							// go utils.SaveAnalyticsData(initConf)
						}
						// commander.Status(c, args, ntype, key)
						commander.Logs(c, args, ntype, key)
					} else {
						logger.Fatalf("Spacecores '%s' of type '%s' not recognized", spacecore, ntype)
					}
				},
			}
			if commander, ok := commanderRegistry[fmt.Sprintf("%s-%s", currentSpacecore, nodeType)]; ok {
				commander.AddFlags(nodeCmd)
			}
			spacecoreCmd.AddCommand(nodeCmd)
		}
		logsCmd.AddCommand(spacecoreCmd)
	}
	commands = append(commands, logsCmd)

	return commands, nil
}

// Implement Logs function for BaseCommander
func (b *BaseCommander) Logs(c *cobra.Command, args []string, mode Mode, node_info string) {
	logger := log.GetLogger(c.Context())
	logger.Info("Getting logs of " + node_info)
	PIDFile := utils.GetPIDFileName(node_info)
	_, err := os.Stat(PIDFile)
	if err == nil {
		_, err := ioutil.ReadFile(PIDFile)
		if err != nil {
			logger.Infof("%v not running\n", node_info)
		} else {
			logFile := "/tmp/" + node_info + ".log"
			cmd := exec.Command("tail", "-f", logFile)
			cmd.Stdout = os.Stdout
			cmd.Stderr = os.Stderr
			cmd.Run()
		}
	} else {
		logger.Infof("%v not running\n", node_info)
	}
}
