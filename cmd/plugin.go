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
			switch action {
			case "start":
				commands <- PluginCommander{Name: StartPluginCmdName, PluginPath: args[0]}
			case "stop":
				commands <- PluginCommander{Name: StopPluginCmdName, PluginPath: args[0]}
			case "restart":
				commands <- PluginCommander{Name: RestartPluginCmdName, PluginPath: args[0]}
			case "status":
				commands <- PluginCommander{Name: StatusPluginCmdName, PluginPath: args[0]}
			}

			// pass the plugin name i.e. plugin path to the goroutine
			// commands <- PluginCommander{Name: StartPluginCmdName, PluginPath: args[0]}

			// wait for a response from the plugin
			resp := <-responses
			logger.Info("Response from plugin:", resp)
		},
	}
	return pluginCmd
}

// started, err := pg.Start(context.Background(), &proto.StartRequest{})

// if err != nil {
// 	logger.Info("Error starting plugin:", err)
// }

// logger.Info("Plugin is now running. Press CTRL+C to exit")

// should there be a way to keep the plugin running?
// if so, how?
// should we just keep the plugin running until the user exits?
// or should we have a way to keep the plugin running in the background?

// we should have a way to keep the plugin running in the background
// so that the user can interact with it

// we should also have a way to stop the plugin

// we should also have a way to restart the plugin

// we should also have a way to get the status of the plugin

// we should also have a way to get the logs of the plugin

// we should also have a way to get the metrics of the plugin

// the plugin needs to keep running in the background

// if started.GetStatus() == "started" {
// 	logger.Info("Plugin started successfully")

// 	// keep the plugin running
// 	<-make(chan struct{})

// } else {
// 	logger.Info("Error starting plugin")
// }

// msg, err := pg.Status(context.Background(), &proto.StatusRequest{})
// if err != nil {
// 	logger.Info("Error getting status:", err)
// }
// // show status of the plugin
// logger.Infof("PG ID: %s, status: %s", msg.GetPluginId(), msg.GetStatus())
