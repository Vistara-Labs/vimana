package cli

import (
	"fmt"
	"os/exec"
	"strings"
	"vimana/cmd/utils"

	"github.com/spf13/cobra"
)

type GmworldDaCommander struct {
	BaseCommander
}

func NewGmworldDaCommander(node_type string) *GmworldDaCommander {
	return &GmworldDaCommander{
		BaseCommander: BaseCommander{NodeType: "init"},
	}
}

func (a *GmworldDaCommander) AddFlags(cmd *cobra.Command) {
}

func (a *GmworldDaCommander) Init(cmd *cobra.Command, args []string, mode Mode, node_info string) error {
	utils.ExecBashCmd(exec.Command("bash", mode.Download), node_info, utils.WithOutputToStdout(), utils.WithErrorsToStderr())
	a.initComponentManager("gmworld", mode.Binary)
	return a.componentMgr.InitializeConfig()
}

func (a *GmworldDaCommander) Run(cmd *cobra.Command, args []string, mode Mode, node_info string) {
	cmdexecute := a.componentMgr.GetStartCmd()
	utils.ExecBashCmd(cmdexecute, node_info, utils.WithOutputToStdout(), utils.WithErrorsToStderr())
}

func (a *GmworldDaCommander) Start(cmd *cobra.Command, args []string, mode Mode, node_info string) {
	node_info_arr := strings.Split(node_info, "-")
	a.Init(cmd, args, mode, node_info_arr[0])
	fmt.Println("executing start command")
	cmdexecute := a.componentMgr.GetStartCmd()
	fmt.Println(cmdexecute)
	utils.ExecBashCmd(cmdexecute, node_info, utils.WithOutputToStdout(), utils.WithErrorsToStderr())
}

func (a *GmworldDaCommander) Stop(cmd *cobra.Command, args []string, mode Mode, node_info string) {
	fmt.Println("Stopping Celestia bridge node")
}

func (a *GmworldDaCommander) Status(cmd *cobra.Command, args []string, mode Mode, node_info string) {
	fmt.Println("Getting status of Celestia bridge node")
}

type GmworldRollupCommander struct {
	BaseCommander
}

func NewGmworldRollupCommander(node_type string) *GmworldRollupCommander {
	return &GmworldRollupCommander{
		BaseCommander: BaseCommander{NodeType: "start"},
	}
}

func (a *GmworldRollupCommander) AddFlags(cmd *cobra.Command) {
}

func (a *GmworldRollupCommander) Init(cmd *cobra.Command, args []string, mode Mode, node_info string) error {
	utils.ExecBashCmd(exec.Command("bash", mode.Download), node_info, utils.WithOutputToStdout(), utils.WithErrorsToStderr())
	a.initComponentManager("rollup", mode.Binary)
	return a.componentMgr.InitializeConfig()
}

func (a *GmworldRollupCommander) Run(cmd *cobra.Command, args []string, mode Mode, node_info string) {
	cmdexecute := a.componentMgr.GetStartCmd()
	utils.ExecBashCmd(cmdexecute, node_info, utils.WithOutputToStdout(), utils.WithErrorsToStderr())
}

func (a *GmworldRollupCommander) Start(cmd *cobra.Command, args []string, mode Mode, node_info string) {
	node_info_arr := strings.Split(node_info, "-")
	a.Init(cmd, args, mode, node_info_arr[0])
	fmt.Println(a.componentMgr)
	fmt.Println("executing start command")
	cmdexecute := a.componentMgr.GetStartCmd()
	fmt.Println(cmdexecute)
	utils.ExecBashCmd(cmdexecute, node_info, utils.WithOutputToStdout(), utils.WithErrorsToStderr())
}

func (a *GmworldRollupCommander) Stop(cmd *cobra.Command, args []string, mode Mode, node_info string) {
	fmt.Println("Stopping Celestia bridge node")
}

func (a *GmworldRollupCommander) Status(cmd *cobra.Command, args []string, mode Mode, node_info string) {
	fmt.Println("Getting status of Celestia bridge node")
}
