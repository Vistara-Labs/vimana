package cli

import (
	"fmt"
	"os/exec"
	"vimana/cmd/utils"

	"github.com/spf13/cobra"
)

type GmdCommander struct {
	BaseCommander
}

func NewGmdCommander() *GmdCommander {
	return &GmdCommander{
		BaseCommander: BaseCommander{NodeType: "light"},
	}
}

func (a *GmdCommander) AddFlags(cmd *cobra.Command) {
}

func (a *GmdCommander) Init(cmd *cobra.Command, args []string, mode Mode) error {
	utils.ExecBashCmd(exec.Command("bash", mode.Download), utils.WithOutputToStdout(), utils.WithErrorsToStderr())
	a.initComponentManager("gmd", mode.Binary)
	return a.componentMgr.InitializeConfig()
}

func (a *GmdCommander) Run(cmd *cobra.Command, args []string, mode Mode) {
	cmdexecute := a.componentMgr.GetStartCmd()
	utils.ExecBashCmd(cmdexecute, utils.WithOutputToStdout(), utils.WithErrorsToStderr())
}

func (a *GmdCommander) Start(cmd *cobra.Command, args []string, mode Mode) {
	a.Init(cmd, args, mode)
	fmt.Println(a.componentMgr)
	fmt.Println("executing start command")
	cmdexecute := a.componentMgr.GetStartCmd()
	fmt.Println(cmdexecute)
	utils.ExecBashCmd(cmdexecute, utils.WithOutputToStdout(), utils.WithErrorsToStderr())
}

func (a *GmdCommander) Stop(cmd *cobra.Command, args []string, mode Mode) {
	fmt.Println("Stopping Celestia bridge node")
}

func (a *GmdCommander) Status(cmd *cobra.Command, args []string, mode Mode) {
	fmt.Println("Getting status of Celestia bridge node")
}
