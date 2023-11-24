package cli

import (
	"fmt"
	"os/exec"
	"vimana/cmd/utils"

	"github.com/spf13/cobra"
)

type NostreamCommander struct {
	BaseCommander
}

func NewNostreamCommander() *NostreamCommander {
	return &NostreamCommander{
		BaseCommander: BaseCommander{NodeType: "light"},
	}
}

func (a *NostreamCommander) AddFlags(cmd *cobra.Command) {
}

func (a *NostreamCommander) Init(cmd *cobra.Command, args []string, mode Mode) error {
	utils.ExecBashCmd(exec.Command("bash", mode.Download), utils.WithOutputToStdout(), utils.WithErrorsToStderr())
	a.initComponentManager("mostream", mode.Binary)
	return a.componentMgr.InitializeConfig()
}

func (a *NostreamCommander) Run(cmd *cobra.Command, args []string, mode Mode) {
	cmdexecute := a.componentMgr.GetStartCmd()
	utils.ExecBashCmd(cmdexecute, utils.WithOutputToStdout(), utils.WithErrorsToStderr())
}

func (a *NostreamCommander) Start(cmd *cobra.Command, args []string, mode Mode) {
	a.Init(cmd, args, mode)
	fmt.Println(a.componentMgr)
	fmt.Println("executing start command")
	cmdexecute := a.componentMgr.GetStartCmd()
	fmt.Println(cmdexecute)
	utils.ExecBashCmd(cmdexecute, utils.WithOutputToStdout(), utils.WithErrorsToStderr())
}

func (a *NostreamCommander) Stop(cmd *cobra.Command, args []string, mode Mode) {
	fmt.Println("Stopping Celestia bridge node")
}

func (a *NostreamCommander) Status(cmd *cobra.Command, args []string, mode Mode) {
	fmt.Println("Getting status of Celestia bridge node")
}
