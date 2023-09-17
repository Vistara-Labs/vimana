package cli

import (
	"fmt"
	"os/exec"
	"vimana/cmd/utils"
	"vimana/components"

	"github.com/spf13/cobra"
)

type CelestiaLightCommander struct {
	Name string
}

func (c *CelestiaLightCommander) Init(cmd *cobra.Command, args []string, mode Mode) error {
	// if isValidUrl(mode.Download) {
	// 	fmt.Println("Downloading Celestia init script from ", mode.Download)
	// 	if err := downloadFile(mode.Download); err != nil {
	// 		return fmt.Errorf("failed to download file: %v", err)
	// 	}
	// } else {
	// 	fmt.Println("Skipping download of Celestia init script from ", mode.Download)
	// }
	fmt.Println("Downloading Celestia from init script ", mode.Download)

	utils.ExecBashCmd(exec.Command("bash", mode.Download), utils.WithOutputToStdout(), utils.WithErrorsToStderr())
	fmt.Println("After from init script ", mode.Download)
	compmanager := components.NewComponentManager("celestia", mode.Binary)
	err := compmanager.InitializeConfig()
	if err != nil {
		return err
	}
	return nil
}

func (c *CelestiaLightCommander) Start(cmd *cobra.Command, args []string, mode Mode) {
	fmt.Println("Starting Celestia light node ", mode, args)

	c.Init(cmd, args, mode)
	compmanager := components.NewComponentManager("celestia", mode.Binary)

	cmdexecute := compmanager.GetStartCmd()

	// args = []string{
	// 	"light", "start",
	// 	"--core.ip", "consensus-full-arabica-9.celestia-arabica.com",
	// 	"--node.store", filepath.Join("~/.vimana/celestia/", "da-light-node"),
	// 	"--gateway",
	// 	"--gateway.deprecated-endpoints",
	// 	"--p2p.network", "arabica",
	// }

	// // RunCmd(mode.Binary, args...) // add start command
	// cmdexecute := exec.Command(
	// 	mode.Binary, args...,
	// )
	fmt.Println(cmdexecute, args)
	utils.ExecBashCmd(cmdexecute, utils.WithOutputToStdout(), utils.WithErrorsToStderr())
}

// Move exec Command snippet to a separate function
func RunCmd(execCmd string, args ...string) {
	cm := exec.Command(execCmd, args...)
	fmt.Println("Running command: ", cm)
	output, err := cm.CombinedOutput()
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(string(output))
}

func (c *CelestiaLightCommander) Stop(cmd *cobra.Command, args []string, mode Mode) {
	// Implementation for "start" command for Celestia light node
	fmt.Println("Stopping Celestia light node")
}

func (c *CelestiaLightCommander) Status(cmd *cobra.Command, args []string, mode Mode) {
	// Implementation for "start" command for Celestia light node
	fmt.Println("Getting status of Celestia light node")
}

type CelestiaBridgeCommander struct {
	Name string
}

func (c *CelestiaBridgeCommander) Init(cmd *cobra.Command, args []string, mode Mode) error {
	// Implementation for "init" command for Celestia light node
	fmt.Println("Initializing Celestia bridge node")
	return nil
}

func (c *CelestiaBridgeCommander) Start(cmd *cobra.Command, args []string, mode Mode) {
	// Implementation for "start" command for Celestia light node
	fmt.Println("Starting Celestia bridge node")
}

func (c *CelestiaBridgeCommander) Stop(cmd *cobra.Command, args []string, mode Mode) {
	// Implementation for "start" command for Celestia light node
	fmt.Println("Stopping Celestia bridge node")
}

func (c *CelestiaBridgeCommander) Status(cmd *cobra.Command, args []string, mode Mode) {
	// Implementation for "start" command for Celestia light node
	fmt.Println("Getting status of Celestia bridge node")
}
