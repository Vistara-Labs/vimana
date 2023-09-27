package cli

import (
	"fmt"
	"vimana/cmd/utils"

	"github.com/spf13/cobra"
)

type AvailLightCommander struct {
	BaseCommander
}

func NewAvailLightCommander() *AvailLightCommander {
	return &AvailLightCommander{
		BaseCommander: BaseCommander{NodeType: "light"},
	}
}

func (a *AvailLightCommander) Init(cmd *cobra.Command, args []string, mode Mode) error {
	// utils.ExecBashCmd(exec.Command("bash", mode.Download), utils.WithOutputToStdout(), utils.WithErrorsToStderr())
	a.initComponentManager("avail", mode.Binary)
	// return c.componentMgr.InitializeConfig()
	// Initialization is in availup.sh
	return nil
}

func (a *AvailLightCommander) Run(cmd *cobra.Command, args []string, mode Mode) {
	cmdexecute := a.componentMgr.GetStartCmd()
	utils.ExecBashCmd(cmdexecute, utils.WithOutputToStdout(), utils.WithErrorsToStderr())
}

func (a *AvailLightCommander) Start(cmd *cobra.Command, args []string, mode Mode) {
	a.Init(cmd, args, mode)
	fmt.Println(a.componentMgr)
	fmt.Println(a)
	cmdexecute := a.componentMgr.GetStartCmd()
	fmt.Println(cmdexecute)
	utils.ExecBashCmd(cmdexecute, utils.WithOutputToStdout(), utils.WithErrorsToStderr())
}

func (a *AvailLightCommander) Stop(cmd *cobra.Command, args []string, mode Mode) {
	fmt.Println("Stopping Celestia bridge node")
}

func (a *AvailLightCommander) Status(cmd *cobra.Command, args []string, mode Mode) {
	fmt.Println("Getting status of Celestia bridge node")
}
