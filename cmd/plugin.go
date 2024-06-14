package cmd

import (
	"context"
	"vimana/plugins"
	"vimana/plugins/proto"

	"vimana/log"

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
	// logger := log.GetLogger(ctx)

	// setup logging
	logger := log.GetLogger(context.Background())

	// create a channel for sending and receiving commands to the plugin
	commands := make(chan PluginCommander)
	responses := make(chan string)

	// Get the plugin in a separate goroutine and pass the commands channel to it
	go func() {

		for cmd := range commands {
			client := plugins.GetPluginClient(cmd.PluginPath)
			spacecore, err := plugins.SpacecoreGRPCClient(client)
			if err != nil {
				logger.Error("Error getting spacecore plugin:", err)
				responses <- "Error getting spacecore plugin" + err.Error()
			}

			switch cmd.Name {
			case StartPluginCmdName:
				logger.Infof("Starting plugin inside goroutine %s", cmd.PluginPath)
				msg, err := spacecore.Start(context.Background(), &proto.StartRequest{})

				if err != nil {
					responses <- "Error starting plugin" + err.Error()
				}

				logger.Infof("Plugin ID: %s, status: %s", msg.GetPluginId(), msg.GetStatus())
				logger.Info("Plugin is now running. Press CTRL+C to exit")
				responses <- "Plugin started"

			case StopPluginCmdName:
				logger.Info("Stopping plugin")
				client.Kill()

				responses <- "Plugin stopped"

			case RestartPluginCmdName:
				logger.Info("Restarting plugin")
				responses <- "Plugin restarted"

			case LogsPluginCmdName:
				logger.Info("Getting logs of plugin")
				msg, err := spacecore.Logs(context.Background(), &proto.LogsRequest{})
				if err != nil {
					logger.Infof("Error getting logs: %s", err)
					responses <- "Error getting logs" + err.Error()
				}

				logger.Infof("Plugin ID: %s, logs: %s", msg.GetPluginId(), msg.GetLogs())
				responses <- "Plugin logs"
				// select {}

			case StatusPluginCmdName:
				logger.Info("Getting status of plugin")

				msg, err := spacecore.Status(context.Background(), &proto.StatusRequest{})
				if err != nil {
					logger.Infof("Error getting status: %s", err)
					responses <- "Error getting status" + err.Error()
				}

				logger.Infof("Plugin ID: %s, status: %s", msg.GetPluginId(), msg.GetStatus())
				responses <- "Plugin status"
			}
		}
	}()

	pluginCmd := &cobra.Command{
		Use:   "plugin [plugin] [action]",
		Short: "Run a spacecore plugin",
		Args:  cobra.MinimumNArgs(2),
		Run: func(c *cobra.Command, args []string) {
			if len(args) == 0 {
				logger.Info("Please provide a plugin name")
				return
			}
			action := args[1]

			requestedAction := ""
			switch action {
			case "start":
				requestedAction = StartPluginCmdName
				// commands <- PluginCommander{Name: StartPluginCmdName, PluginPath: args[0]}
			case "stop":
				requestedAction = StopPluginCmdName
				// commands <- PluginCommander{Name: StopPluginCmdName, PluginPath: args[0]}
			case "restart":
				requestedAction = RestartPluginCmdName
				// commands <- PluginCommander{Name: RestartPluginCmdName, PluginPath: args[0]}
			case "status":
				requestedAction = StatusPluginCmdName
				// commands <- PluginCommander{Name: StatusPluginCmdName, PluginPath: args[0]}
			case "logs":
				requestedAction = LogsPluginCmdName
				// commands <- PluginCommander{Name: LogsPluginCmdName, PluginPath: args[0]}
			}

			// pass the plugin name i.e. plugin path to the goroutine
			commands <- PluginCommander{Name: requestedAction, PluginPath: args[0]}

			// wait for a response from the plugin
			resp := <-responses
			logger.Info("Response from plugin:", resp)
		},
	}
	return pluginCmd
}
