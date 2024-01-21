package cli

import (
	"fmt"
	"os/exec"
	"strings"
	"vimana/cmd/utils"

	"github.com/spf13/cobra"
)

type EigenOperatorCommander struct {
	BaseCommander
}

func NewEigenOperatorCommander(node_type string) *EigenOperatorCommander {
	return &EigenOperatorCommander{
		BaseCommander: BaseCommander{NodeType: "operator"},
	}
}

func (a *EigenOperatorCommander) AddFlags(cmd *cobra.Command) {
}

func (a *EigenOperatorCommander) Init(cmd *cobra.Command, args []string, mode Mode, node_info string) error {
	utils.ExecBashCmd(exec.Command("bash", mode.Download), node_info, utils.WithOutputToStdout(), utils.WithErrorsToStderr())
	a.initComponentManager("eigen", mode.Binary)
	return a.componentMgr.InitializeConfig()
}

func (a *EigenOperatorCommander) Run(cmd *cobra.Command, args []string, mode Mode, node_info string) {
	cmdexecute := a.componentMgr.GetStartCmd()
	utils.ExecBashCmd(cmdexecute, node_info, utils.WithOutputToStdout(), utils.WithErrorsToStderr())
}

func (a *EigenOperatorCommander) Start(cmd *cobra.Command, args []string, mode Mode, node_info string) {
	node_info_arr := strings.Split(node_info, "-")
	a.Init(cmd, args, mode, node_info_arr[0])
	fmt.Println(a.componentMgr)
	fmt.Println("executing start command")
	cmdexecute := a.componentMgr.GetStartCmd()
	fmt.Println(cmdexecute)
	utils.ExecBashCmd(cmdexecute, node_info, utils.WithOutputToStdout(), utils.WithErrorsToStderr())
}

func (a *EigenOperatorCommander) Stop(cmd *cobra.Command, args []string, mode Mode, node_info string) {
	fmt.Println("Stopping Eigen operator node")
}

func (a *EigenOperatorCommander) Status(cmd *cobra.Command, args []string, mode Mode, node_info string) {
	fmt.Println("Getting status of Eigen operator node")
}
