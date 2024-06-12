package cli

import (
	"io/ioutil"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"vimana/cmd/utils"
	"vimana/log"

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

func (a *EigenOperatorCommander) Install(cmd *cobra.Command, args []string, mode Mode, node_info string) {
	return
}

func (a *EigenOperatorCommander) Init(cmd *cobra.Command, args []string, mode Mode, node_info string) error {
	utils.ExecBashCmd(exec.Command("bash", mode.Download), node_info, utils.WithOutputToStdout(), utils.WithErrorsToStderr())
	a.initSpacecoreManager("eigen", mode.Binary)
	return a.spacecoresMgr.InitializeConfig()
}

func (a *EigenOperatorCommander) Run(cmd *cobra.Command, args []string, mode Mode, node_info string) {
	cmdexecute := a.spacecoresMgr.GetStartCmd()
	utils.ExecBinaryCmd(cmdexecute, node_info, utils.WithOutputToStdout(), utils.WithErrorsToStderr())
}

func (a *EigenOperatorCommander) Start(cmd *cobra.Command, args []string, mode Mode, node_info string) {
	logger := log.GetLogger(cmd.Context())
	node_info_arr := strings.Split(node_info, "-")
	a.Init(cmd, args, mode, node_info_arr[0])
	logger.Info(a.spacecoresMgr)
	logger.Info("executing start command")
	cmdexecute := a.spacecoresMgr.GetStartCmd()
	logger.Info(cmdexecute)
	utils.ExecBinaryCmd(cmdexecute, node_info, utils.WithOutputToStdout(), utils.WithErrorsToStderr())
}

func (a *EigenOperatorCommander) Stop(cmd *cobra.Command, args []string, mode Mode, node_info string) {
	logger := log.GetLogger(cmd.Context())
	logger.Info("Stopping Eigen operator node")
	PIDFile := utils.GetPIDFileName(node_info)
	if _, err := os.Stat(PIDFile); err == nil {
		data, err := ioutil.ReadFile(PIDFile)
		if err != nil {
			logger.Info("Not running")
			return
		}
		ProcessID, err := strconv.Atoi(string(data))

		if err != nil {
			logger.Info("Unable to read and parse process id found in ", PIDFile)
			return
		}

		process, err := os.FindProcess(ProcessID)

		if err != nil {
			logger.Infof("Unable to find process ID [%v] with error %v \n", ProcessID, err)
			return
		}
		// remove PID file
		os.Remove(PIDFile)

		node_info_arr := strings.Split(node_info, "-")
		logger.Info("Stopping " + node_info_arr[0] + " " + node_info_arr[1] + " node")
		// kill process and exit immediately
		err = process.Kill()

		if err != nil {
			logger.Infof("Unable to kill process ID [%v] with error %v \n", ProcessID, err)
		} else {
			logger.Infof("Killed process ID [%v]\n", ProcessID)
		}
	} else {
		logger.Info("Not running.")
	}
}

func (a *EigenOperatorCommander) Status(cmd *cobra.Command, args []string, mode Mode, node_info string) {
	logger := log.GetLogger(cmd.Context())

	node_info_arr := strings.Split(node_info, "-")
	logger.Info("Getting status of " + node_info_arr[0] + " " + node_info_arr[1] + " node")

	PIDFile := utils.GetPIDFileName(node_info)
	if _, err := os.Stat(PIDFile); err == nil {
		_, err := ioutil.ReadFile(PIDFile)
		if err != nil {
			logger.Info("Not running")
		} else {
			logger.Info(node_info_arr[0] + " " + node_info_arr[1] + " node is running")
		}
	} else {
		logger.Info("Not running.")
	}
}
