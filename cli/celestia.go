package cli

import (
	"fmt"
	"os/exec"
	"strings"
	"vimana/cmd/utils"
	"vimana/components"

	"github.com/spf13/cobra"
)

type CelestiaLightCommander struct {
	CelestiaNetwork string
	CelestiaRPC     string
	BaseCommander
}

type CelestiaBridgeCommander struct {
	CelestiaNetwork string
	CelestiaRPC     string
	BaseCommander
}

// Reference from roller
const (
	DefaultCelestiaRPC     = "rpc.celestia.pops.one"
	DefaultCelestiaNetwork = "celestia"
)

func NewCelestiaLightCommander(node_type string) *CelestiaLightCommander {
	return &CelestiaLightCommander{
		BaseCommander: BaseCommander{NodeType: "light"},
	}
}

func NewCelestiaBridgeCommander(node_type string) *CelestiaBridgeCommander {
	return &CelestiaBridgeCommander{
		BaseCommander: BaseCommander{NodeType: "bridge"},
	}
}
func (c *CelestiaLightCommander) AddFlags(cmd *cobra.Command) {
	cmd.Flags().StringVar(&c.CelestiaNetwork, "network", DefaultCelestiaNetwork, "Specifies the Celestia network")
	cmd.Flags().StringVar(&c.CelestiaRPC, "rpc", DefaultCelestiaRPC, "Specifies the Celestia RPC endpoint")
}

func (c *CelestiaLightCommander) Init(cmd *cobra.Command, args []string, mode Mode, node_info string) error {
	utils.ExecBashCmd(exec.Command("bash", mode.Download), node_info, utils.WithOutputToStdout(), utils.WithErrorsToStderr())

	c.config = &components.ComponentConfig{
		RPC:     c.CelestiaRPC,
		Network: c.CelestiaNetwork,
	}

	// c.componentMgr = components.NewComponentManager("celestia", mode.Binary, c.NodeType, config)
	c.initComponentManager("celestia", mode.Binary)

	// c.initComponentManager("celestia", mode.Binary, c.CelestiaNetwork, c.CelestiaRPC)
	return c.componentMgr.InitializeConfig()
}

func (c *CelestiaLightCommander) Run(cmd *cobra.Command, args []string, mode Mode, node_info string) {
	node_info_arr := strings.Split(node_info, "-")
	c.Init(cmd, args, mode, node_info_arr[0])
	cmdexecute := c.componentMgr.GetStartCmd()
	utils.ExecBashCmd(cmdexecute, node_info, utils.WithOutputToStdout(), utils.WithErrorsToStderr())
}

func (c *CelestiaLightCommander) Start(cmd *cobra.Command, args []string, mode Mode, node_info string) {
	node_info_arr := strings.Split(node_info, "-")
	c.Init(cmd, args, mode, node_info_arr[0])
	cmdexecute := c.componentMgr.GetStartCmd()
	fmt.Println(cmdexecute)
	utils.ExecBashCmd(cmdexecute, node_info, utils.WithOutputToStdout(), utils.WithErrorsToStderr())
}

func (c *CelestiaLightCommander) Stop(cmd *cobra.Command, args []string, mode Mode, node_info string) {
	// Implementation for "start" command for Celestia light node
	fmt.Println("To implement: Celestia light node stop")
}

func (c *CelestiaLightCommander) Status(cmd *cobra.Command, args []string, mode Mode, node_info string) {
	// Implementation for "start" command for Celestia light node
	fmt.Println("To implement: Celestia light node status")
}

func (c *CelestiaBridgeCommander) AddFlags(cmd *cobra.Command) {
	cmd.Flags().StringVar(&c.CelestiaNetwork, "network", DefaultCelestiaNetwork, "Specifies the Celestia network")
	cmd.Flags().StringVar(&c.CelestiaRPC, "rpc", DefaultCelestiaRPC, "Specifies the Celestia RPC endpoint")
}

func (c *CelestiaBridgeCommander) Init(cmd *cobra.Command, args []string, mode Mode, node_info string) error {
	utils.ExecBashCmd(exec.Command("bash", mode.Download), node_info, utils.WithOutputToStdout(), utils.WithErrorsToStderr())
	c.config = &components.ComponentConfig{
		RPC:     c.CelestiaRPC,
		Network: c.CelestiaNetwork,
	}

	c.initComponentManager("celestia", mode.Binary)

	return c.componentMgr.InitializeConfig()
}

func (c *CelestiaBridgeCommander) Start(cmd *cobra.Command, args []string, mode Mode, node_info string) {
	node_info_arr := strings.Split(node_info, "-")
	c.Init(cmd, args, mode, node_info_arr[0])
	// fmt.Println("Starting Celestia bridge node", c)
	cmdexecute := c.componentMgr.GetStartCmd()
	fmt.Println(cmdexecute)
	utils.ExecBashCmd(cmdexecute, node_info, utils.WithOutputToStdout(), utils.WithErrorsToStderr())
}

func (c *CelestiaBridgeCommander) Stop(cmd *cobra.Command, args []string, mode Mode, node_info string) {
	// Implementation for "start" command for Celestia light node
	fmt.Println("Stopping Celestia bridge node")
}

func (c *CelestiaBridgeCommander) Status(cmd *cobra.Command, args []string, mode Mode, node_info string) {
	// Implementation for "start" command for Celestia light node
	fmt.Println("Getting status of Celestia bridge node")
}
