package cli

import (
	"fmt"
	"os/exec"
	"strings"
	"vimana/cmd/utils"

	"github.com/spf13/cobra"
)

type AvailLightCommander struct {
	BaseCommander
}

func NewAvailLightCommander(node_type string) *AvailLightCommander {
	return &AvailLightCommander{
		BaseCommander: BaseCommander{NodeType: "light"},
	}
}

func (a *AvailLightCommander) AddFlags(cmd *cobra.Command) {
}

func (a *AvailLightCommander) Init(cmd *cobra.Command, args []string, mode Mode, node_info string) error {
	utils.ExecBashCmd(exec.Command("bash", mode.Download), node_info, utils.WithOutputToStdout(), utils.WithErrorsToStderr())
	a.initComponentManager("avail", mode.Binary)
	return a.componentMgr.InitializeConfig()
}

func (a *AvailLightCommander) Run(cmd *cobra.Command, args []string, mode Mode, node_info string) {
	cmdexecute := a.componentMgr.GetStartCmd()
	utils.ExecBashCmd(cmdexecute, node_info, utils.WithOutputToStdout(), utils.WithErrorsToStderr())
}

func (a *AvailLightCommander) Start(cmd *cobra.Command, args []string, mode Mode, node_info string) {
	node_info_arr := strings.Split(node_info, "-")
	a.Init(cmd, args, mode, node_info_arr[0])
	fmt.Println(a.componentMgr)
	fmt.Println("executing start command")
	cmdexecute := a.componentMgr.GetStartCmd()
	fmt.Println(cmdexecute)
	utils.ExecBashCmd(cmdexecute, node_info, utils.WithOutputToStdout(), utils.WithErrorsToStderr())
}

func (a *AvailLightCommander) Stop(cmd *cobra.Command, args []string, mode Mode, node_info string) {
	fmt.Println("Stopping Celestia bridge node")
}

func (a *AvailLightCommander) Status(cmd *cobra.Command, args []string, mode Mode, node_info string) {
	fmt.Println("Getting status of Celestia bridge node")
}
