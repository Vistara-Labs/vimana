package cli

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
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

func (a *CelestiaLightCommander) Install(cmd *cobra.Command, args []string, mode Mode, node_info string) {
	_, err := os.Stat(mode.Download)
	if err == nil {
		// true
		utils.ExecBinaryCmd(exec.Command("bash", mode.Download), node_info, utils.WithOutputToStdout(), utils.WithErrorsToStderr())
	} else if os.IsNotExist(err) {
		// false
		currentDir, err := os.Getwd()
		if err != nil {
			fmt.Println("Error getting current directory:", err)
			return
		}
		parentDir := filepath.Dir(currentDir)
		utils.ExecBinaryCmd(exec.Command("bash", parentDir+"/scripts/init.sh"), node_info, utils.WithOutputToStdout(), utils.WithErrorsToStderr())
	} else {

		fmt.Printf("error：%v\n", err)
	}
	return
	return
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
	utils.ExecBinaryCmd(cmdexecute, node_info, utils.WithOutputToStdout(), utils.WithErrorsToStderr())
}

func (c *CelestiaLightCommander) Start(cmd *cobra.Command, args []string, mode Mode, node_info string) {

	// check if daemon already running.
	PIDFile := utils.GetPIDFileName(node_info)
	if _, err := os.Stat(PIDFile); err == nil {
		fmt.Println("Already running or " + PIDFile + " file exist.")
		return
	}

	node_info_arr := strings.Split(node_info, "-")
	c.Init(cmd, args, mode, node_info_arr[0])
	cmdexecute := c.componentMgr.GetStartCmd()
	fmt.Println(cmdexecute)
	utils.ExecBinaryCmd(cmdexecute, node_info, utils.WithOutputToStdout(), utils.WithErrorsToStderr())
}

func (c *CelestiaLightCommander) Stop(cmd *cobra.Command, args []string, mode Mode, node_info string) {
	// Implementation for "start" command for Celestia light node
	PIDFile := utils.GetPIDFileName(node_info)
	if _, err := os.Stat(PIDFile); err == nil {
		data, err := ioutil.ReadFile(PIDFile)
		if err != nil {
			fmt.Println("Not running")
			return
		}
		ProcessID, err := strconv.Atoi(string(data))

		if err != nil {
			fmt.Println("Unable to read and parse process id found in ", PIDFile)
			return
		}

		process, err := os.FindProcess(ProcessID)

		if err != nil {
			fmt.Printf("Unable to find process ID [%v] with error %v \n", ProcessID, err)
			return
		}
		// remove PID file
		os.Remove(PIDFile)

		node_info_arr := strings.Split(node_info, "-")
		fmt.Println("Stopping " + node_info_arr[0] + " " + node_info_arr[1] + " node")
		// kill process and exit immediately
		err = process.Kill()

		if err != nil {
			fmt.Printf("Unable to kill process ID [%v] with error %v \n", ProcessID, err)
		} else {
			fmt.Printf("Killed process ID [%v]\n", ProcessID)
		}
	} else {
		fmt.Println("Not running.")
	}
}

func (c *CelestiaLightCommander) Status(cmd *cobra.Command, args []string, mode Mode, node_info string) {
	// Implementation for "start" command for Celestia light node
	node_info_arr := strings.Split(node_info, "-")
	fmt.Println("Getting status of " + node_info_arr[0] + " " + node_info_arr[1] + " node")

	PIDFile := utils.GetPIDFileName(node_info)
	if _, err := os.Stat(PIDFile); err == nil {
		_, err := ioutil.ReadFile(PIDFile)
		if err != nil {
			fmt.Println("Not running")
		} else {
			fmt.Println(node_info_arr[0] + " " + node_info_arr[1] + " node is running")
		}
	} else {
		fmt.Println("Not running.")
	}
}

func (c *CelestiaBridgeCommander) AddFlags(cmd *cobra.Command) {
	cmd.Flags().StringVar(&c.CelestiaNetwork, "network", DefaultCelestiaNetwork, "Specifies the Celestia network")
	cmd.Flags().StringVar(&c.CelestiaRPC, "rpc", DefaultCelestiaRPC, "Specifies the Celestia RPC endpoint")
}

func (a *CelestiaBridgeCommander) Install(cmd *cobra.Command, args []string, mode Mode, node_info string) {
	return
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

	// check if daemon already running.
	PIDFile := utils.GetPIDFileName(node_info)
	if _, err := os.Stat(PIDFile); err == nil {
		fmt.Println("Already running or " + PIDFile + " file exist.")
		return
	}

	node_info_arr := strings.Split(node_info, "-")
	c.Init(cmd, args, mode, node_info_arr[0])
	// fmt.Println("Starting Celestia bridge node", c)
	cmdexecute := c.componentMgr.GetStartCmd()
	fmt.Println("Start: ", cmdexecute)
	utils.ExecBinaryCmd(cmdexecute, node_info, utils.WithOutputToStdout(), utils.WithErrorsToStderr())
}

func (c *CelestiaBridgeCommander) Stop(cmd *cobra.Command, args []string, mode Mode, node_info string) {
	// Implementation for "start" command for Celestia light node
	fmt.Println("Stopping Celestia bridge node")
	PIDFile := utils.GetPIDFileName(node_info)
	if _, err := os.Stat(PIDFile); err == nil {
		data, err := ioutil.ReadFile(PIDFile)
		if err != nil {
			fmt.Println("Not running")
			return
		}
		ProcessID, err := strconv.Atoi(string(data))

		if err != nil {
			fmt.Println("Unable to read and parse process id found in ", PIDFile)
			return
		}

		process, err := os.FindProcess(ProcessID)

		if err != nil {
			fmt.Printf("Unable to find process ID [%v] with error %v \n", ProcessID, err)
			return
		}
		// remove PID file
		os.Remove(PIDFile)

		node_info_arr := strings.Split(node_info, "-")
		fmt.Println("Stopping " + node_info_arr[0] + " " + node_info_arr[1] + " node")
		// kill process and exit immediately
		err = process.Kill()

		if err != nil {
			fmt.Printf("Unable to kill process ID [%v] with error %v \n", ProcessID, err)
		} else {
			fmt.Printf("Killed process ID [%v]\n", ProcessID)
		}
	} else {
		fmt.Println("Not running.")
	}
}

func (c *CelestiaBridgeCommander) Status(cmd *cobra.Command, args []string, mode Mode, node_info string) {
	// Implementation for "start" command for Celestia light node
	node_info_arr := strings.Split(node_info, "-")
	fmt.Println("Getting status of " + node_info_arr[0] + " " + node_info_arr[1] + " node")

	PIDFile := utils.GetPIDFileName(node_info)
	if _, err := os.Stat(PIDFile); err == nil {
		_, err := ioutil.ReadFile(PIDFile)
		if err != nil {
			fmt.Println("Not running")
		} else {
			fmt.Println(node_info_arr[0] + " " + node_info_arr[1] + " node is running")
		}
	} else {
		fmt.Println("Not running.")
	}
}
