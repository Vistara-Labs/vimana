package cli

import (
	"io/ioutil"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"vimana/cmd/utils"
	"vimana/config"
	"vimana/log"

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
	logger := log.GetLogger(cmd.Context())
	logger.Info(a.spacecoresMgr)
	logger.Info("executing install command")
	binaryPath := string([]rune(mode.Binary)[0:strings.LastIndex(mode.Binary, "/")])
	binaryName := string([]rune(mode.Binary)[strings.LastIndex(mode.Binary, "/")+1:])
	utils.ExecBashCmd(exec.Command("bash", mode.Install, binaryPath, binaryName), node_info, utils.WithOutputToStdout(), utils.WithErrorsToStderr())
}

func (a *UniversalCommander) Init(cmd *cobra.Command, args []string, mode Mode, node_info string) error {

	utils.ExecBashCmd(exec.Command("bash", mode.Download), node_info, utils.WithOutputToStdout(), utils.WithErrorsToStderr())
	node_info_arr := strings.Split(node_info, "-")
	a.initSpacecoreManager(config.SpacecoreType(node_info_arr[0]), mode.Binary)
	return a.spacecoresMgr.InitializeConfig()
}

func (a *UniversalCommander) Run(cmd *cobra.Command, args []string, mode Mode, node_info string) {
	cmdexecute := a.spacecoresMgr.GetStartCmd()
	utils.ExecBinaryCmd(cmdexecute, node_info, utils.WithOutputToStdout(), utils.WithErrorsToStderr())
}

func (a *UniversalCommander) Start(cmd *cobra.Command, args []string, mode Mode, node_info string) {
	logger := log.GetLogger(cmd.Context())
	PIDFile := utils.GetPIDFileName(node_info)
	if _, err := os.Stat(PIDFile); err == nil {
		logger.Info("Already running or " + PIDFile + " file exist.")
		return
	}

	logger.Info("executing start command")
	binaryPath := string([]rune(mode.Binary)[0:strings.LastIndex(mode.Binary, "/")])
	binaryName := string([]rune(mode.Binary)[strings.LastIndex(mode.Binary, "/")+1:])
	logger.Info(binaryPath)

	utils.ExecBinaryCmd(exec.Command("bash", mode.Start, binaryPath, binaryName), node_info, utils.WithOutputToStdout(), utils.WithErrorsToStderr())
}

func (a *UniversalCommander) Stop(cmd *cobra.Command, args []string, mode Mode, node_info string) {
	logger := log.GetLogger(cmd.Context())
	PIDFile := utils.GetPIDFileName(node_info)
	if _, err := os.Stat(PIDFile); err == nil {
		data, err := ioutil.ReadFile(PIDFile)
		if err != nil {
			logger.Info("Not running")
			return
		}
		ProcessID, err := strconv.Atoi(string(data))
		ProcessID = ProcessID + 1
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

func (a *UniversalCommander) Logs(cmd *cobra.Command, args []string, mode Mode, node_info string) {
	logger := log.GetLogger(cmd.Context())
	PIDFile := utils.GetPIDFileName(node_info)
	_, err := os.Stat(PIDFile)
	if err == nil {
		_, err := ioutil.ReadFile(PIDFile)
		if err != nil {
			logger.Infof("%v not running\n", node_info)
		} else {
			logFile := "/tmp/" + node_info + ".log"
			cmd := exec.Command("tail", "-f", logFile)
			cmd.Stdout = os.Stdout
			cmd.Stderr = os.Stderr
			cmd.Run()
		}
	} else {
		logger.Infof("%v not running\n", node_info)
	}
}

func (a *UniversalCommander) Status(cmd *cobra.Command, args []string, mode Mode, node_info string) {
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
