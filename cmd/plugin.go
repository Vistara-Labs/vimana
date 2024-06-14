package cmd

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"vimana/plugins"
	"vimana/plugins/proto"

	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

type PluginCommander struct {
	Name       string
	PluginPath string
}
type PluginCommanders int

const (
	PluginCmdName        = "plugin"
	StopPluginCmdName    = "stop"
	StartPluginCmdName   = "start"
	RestartPluginCmdName = "restart"
	StatusPluginCmdName  = "status"
	LogsPluginCmdName    = "logs"
	MetricsPluginCmdName = "metrics"
)

func pluginCommand() *cobra.Command {
	pluginCmd := &cobra.Command{
		Use:   "plugin [command] [plugin]",
		Short: "Manage spacecore plugins",
		Args:  cobra.MinimumNArgs(2),
		Run: func(c *cobra.Command, args []string) {
			commands := make(chan PluginCommander)
			responses := make(chan string)

			pluginPath := args[0]
			pluginCmd := args[1]
			pluginName := filepath.Base(pluginPath)
			logger := setupLogger(pluginName)

			// Start the goroutine to hande the plugin command
			go func() {
				for cmd := range commands {
					response := executePluginCommand(cmd, logger)
					responses <- response
				}
			}()

			// pass the plugin name i.e. plugin path to the goroutine

			// Send the plugin command to the channel
			commands <- PluginCommander{Name: pluginCmd, PluginPath: pluginPath}

			// Wait for a response from the plugin
			response := <-responses
			logger.Info(response)
			fmt.Println(response)

		},
	}
	return pluginCmd

}

func executePluginCommand(cmd PluginCommander, logger *logrus.Entry) string {

	client := plugins.GetPluginClient(cmd.PluginPath)
	spacecore, err := plugins.SpacecoreGRPCClient(client)
	if err != nil {
		logger.Error("Error getting spacecore plugin s: ", err)
		return "Error getting spacecore plugin: " + err.Error()
	}

	switch cmd.Name {
	case StartPluginCmdName:
		logger.Infof("Starting plugin %s", cmd.PluginPath)
		msg, err := spacecore.Start(context.Background(), &proto.StartRequest{})
		if err != nil {
			return "Error starting plugin: " + err.Error()
		}
		logger.Infof("Plugin ID: %s, status: %s", msg.GetPluginId(), msg.GetStatus())

		return msg.GetStatus()
	case StopPluginCmdName:
		logger.Infof("Stopping plugin %s", cmd.PluginPath)
		msg, err := spacecore.Stop(context.Background(), &proto.StopRequest{})
		if err != nil {
			return "Error stopping plugin: " + err.Error()
		}
		logger.Infof("Plugin status: %s", msg.GetStatus())
		return "Plugin stopped"
	case StatusPluginCmdName:
		logger.Infof("Checking status of plugin %s", cmd.PluginPath)
		msg, err := spacecore.Status(context.Background(), &proto.StatusRequest{})
		if err != nil {
			return "Error checking plugin status: " + err.Error()
		}
		logger.Infof("Plugin status: %s", msg.GetStatus())
		return "Plugin status: " + msg.GetStatus()
	case LogsPluginCmdName:
		logger.Infof("Fetching logs of plugin %s", cmd.PluginPath)
		msg, err := spacecore.Logs(context.Background(), &proto.LogsRequest{})
		if err != nil {
			return "Error fetching plugin logs: " + err.Error()
		}
		logs := concatenateLogs(msg.GetLogs())
		logger.Infof("Plugin logs: \n%s", logs)
		return logs
	default:
		return "Unknown command: " + cmd.Name
	}
}

func concatenateLogs(logs []string) string {
	return strings.Join(logs, "\n")
}

func setupLogger(pluginName string) *logrus.Entry {
	logFileName := fmt.Sprintf("/tmp/%s.log", pluginName)
	logFile, err := os.OpenFile(logFileName, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Printf("Error opening log file: %v\n", err)
		os.Exit(1)
	}

	// our log.Configure is throwing bad file descriptor, fix later
	logger := logrus.New()
	logger.SetOutput(logFile)
	logger.SetLevel(logrus.DebugLevel)
	logger.SetFormatter(&logrus.TextFormatter{
		FullTimestamp: true,
	})

	return logger.WithFields(logrus.Fields{"plugin": pluginName})
}
