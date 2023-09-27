package cli

import (
	"fmt"
	"os/exec"
	"vimana/cmd/utils"

	"github.com/spf13/cobra"
)

type CelestiaLightCommander struct {
	BaseCommander
}

type CelestiaBridgeCommander struct {
	BaseCommander
}

func NewCelestiaLightCommander() *CelestiaLightCommander {
	return &CelestiaLightCommander{
		BaseCommander: BaseCommander{NodeType: "light"},
	}
}

func NewCelestiaBridgeCommander() *CelestiaBridgeCommander {
	return &CelestiaBridgeCommander{
		BaseCommander: BaseCommander{NodeType: "bridge"},
	}
}

func (c *CelestiaLightCommander) Init(cmd *cobra.Command, args []string, mode Mode) error {
	utils.ExecBashCmd(exec.Command("bash", mode.Download), utils.WithOutputToStdout(), utils.WithErrorsToStderr())
	c.initComponentManager("celestia", mode.Binary)
	return c.componentMgr.InitializeConfig()
}

func (c *CelestiaLightCommander) Run(cmd *cobra.Command, args []string, mode Mode) {
	c.Init(cmd, args, mode)
	cmdexecute := c.componentMgr.GetStartCmd()
	utils.ExecBashCmd(cmdexecute, utils.WithOutputToStdout(), utils.WithErrorsToStderr())
}

func (c *CelestiaLightCommander) Start(cmd *cobra.Command, args []string, mode Mode) {
	c.Init(cmd, args, mode)
	cmdexecute := c.componentMgr.GetStartCmd()
	fmt.Println(cmdexecute)
	utils.ExecBashCmd(cmdexecute, utils.WithOutputToStdout(), utils.WithErrorsToStderr())
}

func (c *CelestiaLightCommander) Stop(cmd *cobra.Command, args []string, mode Mode) {
	// Implementation for "start" command for Celestia light node
	fmt.Println("To implement: Celestia light node stop")
}

func (c *CelestiaLightCommander) Status(cmd *cobra.Command, args []string, mode Mode) {
	// Implementation for "start" command for Celestia light node
	fmt.Println("To implement: Celestia light node status")
}

func (c *CelestiaBridgeCommander) Init(cmd *cobra.Command, args []string, mode Mode) error {
	utils.ExecBashCmd(exec.Command("bash", mode.Download), utils.WithOutputToStdout(), utils.WithErrorsToStderr())
	c.initComponentManager("celestia", mode.Binary)
	return c.componentMgr.InitializeConfig()
}

func (c *CelestiaBridgeCommander) Start(cmd *cobra.Command, args []string, mode Mode) {
	c.Init(cmd, args, mode)
	// fmt.Println("Starting Celestia bridge node", c)
	cmdexecute := c.componentMgr.GetStartCmd()
	fmt.Println(cmdexecute)
	utils.ExecBashCmd(cmdexecute, utils.WithOutputToStdout(), utils.WithErrorsToStderr())
}

func (c *CelestiaBridgeCommander) Stop(cmd *cobra.Command, args []string, mode Mode) {
	// Implementation for "start" command for Celestia light node
	fmt.Println("Stopping Celestia bridge node")
}

func (c *CelestiaBridgeCommander) Status(cmd *cobra.Command, args []string, mode Mode) {
	// Implementation for "start" command for Celestia light node
	fmt.Println("Getting status of Celestia bridge node")
}
