package cli

import (
	"fmt"
	"os/exec"
	"vimana/cmd/utils"

	"github.com/spf13/cobra"
)

type EigenOperatorCommander struct {
	BaseCommander
}

func NewEigenOperatorCommander() *EigenOperatorCommander {
	return &EigenOperatorCommander{
		BaseCommander: BaseCommander{NodeType: "operator"},
	}
}

func (a *EigenOperatorCommander) AddFlags(cmd *cobra.Command) {
}

func (a *EigenOperatorCommander) Init(cmd *cobra.Command, args []string, mode Mode) error {
	utils.ExecBashCmd(exec.Command("bash", mode.Download), utils.WithOutputToStdout(), utils.WithErrorsToStderr())
	a.initComponentManager("eigen", mode.Binary)
	return a.componentMgr.InitializeConfig()
}

func (a *EigenOperatorCommander) Run(cmd *cobra.Command, args []string, mode Mode) {
	cmdexecute := a.componentMgr.GetStartCmd()
	utils.ExecBashCmd(cmdexecute, utils.WithOutputToStdout(), utils.WithErrorsToStderr())
}

func (a *EigenOperatorCommander) Start(cmd *cobra.Command, args []string, mode Mode) {
	a.Init(cmd, args, mode)
	fmt.Println(a.componentMgr)
	fmt.Println("executing start command")
	cmdexecute := a.componentMgr.GetStartCmd()
	fmt.Println(cmdexecute)
	utils.ExecBashCmd(cmdexecute, utils.WithOutputToStdout(), utils.WithErrorsToStderr())
}

func (a *EigenOperatorCommander) Stop(cmd *cobra.Command, args []string, mode Mode) {
	fmt.Println("Stopping Eigen operator node")
}

func (a *EigenOperatorCommander) Status(cmd *cobra.Command, args []string, mode Mode) {
	fmt.Println("Getting status of Eigen operator node")
}
