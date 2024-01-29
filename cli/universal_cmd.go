package cli

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"vimana/cmd/utils"
	"vimana/config"

	"github.com/spf13/cobra"
)

type UniversalCommander struct {
	BaseCommander
}

func NewUniversalCommander(node_type string) *UniversalCommander {
	return &UniversalCommander{
		BaseCommander: BaseCommander{NodeType: node_type},
	}
}

func (a *UniversalCommander) AddFlags(cmd *cobra.Command) {
}

func (a *UniversalCommander) Install(cmd *cobra.Command, args []string, mode Mode, node_info string) {

	fmt.Println(a.componentMgr)
	fmt.Println("executing install command")
	//cmdexecute := a.componentMgr.GetStartCmd()
	//fmt.Println(cmdexecute)
	utils.ExecBashCmd(exec.Command("bash", mode.Install), node_info, utils.WithOutputToStdout(), utils.WithErrorsToStderr())
}

func (a *UniversalCommander) Init(cmd *cobra.Command, args []string, mode Mode, node_info string) error {

	utils.ExecBashCmd(exec.Command("bash", mode.Download), node_info, utils.WithOutputToStdout(), utils.WithErrorsToStderr())
	node_info_arr := strings.Split(node_info, "-")
	a.initComponentManager(config.ComponentType(node_info_arr[0]), mode.Binary)
	return a.componentMgr.InitializeConfig()
}

func (a *UniversalCommander) Run(cmd *cobra.Command, args []string, mode Mode, node_info string) {
	cmdexecute := a.componentMgr.GetStartCmd()
	utils.ExecBinaryCmd(cmdexecute, node_info, utils.WithOutputToStdout(), utils.WithErrorsToStderr())
}

func (a *UniversalCommander) Start(cmd *cobra.Command, args []string, mode Mode, node_info string) {

	// check if daemon already running.
	PIDFile := utils.GetPIDFileName(node_info)
	if _, err := os.Stat(PIDFile); err == nil {
		fmt.Println("Already running or " + PIDFile + " file exist.")
		return
	}

	//node_info_arr := strings.Split(node_info, "-")
	//a.Init(cmd, args, mode, node_info_arr[0])
	fmt.Println(a.componentMgr)
	fmt.Println("executing start command")
	//cmdexecute := a.componentMgr.GetStartCmd()
	//fmt.Println(cmdexecute)
	utils.ExecBinaryCmd(exec.Command("bash", mode.Start), node_info, utils.WithOutputToStdout(), utils.WithErrorsToStderr())
}

func (a *UniversalCommander) Stop(cmd *cobra.Command, args []string, mode Mode, node_info string) {
	PIDFile := utils.GetPIDFileName(node_info)
	if _, err := os.Stat(PIDFile); err == nil {
		data, err := ioutil.ReadFile(PIDFile)
		if err != nil {
			fmt.Println("Not running")
			return
		}
		ProcessID, err := strconv.Atoi(string(data))
		ProcessID = ProcessID + 1
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

func (a *UniversalCommander) Status(cmd *cobra.Command, args []string, mode Mode, node_info string) {
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
