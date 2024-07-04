package cmd

import (
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"vimana/cli"

	"vimana/log"

	"github.com/BurntSushi/toml"
	"github.com/spf13/cobra"
	"google.golang.org/grpc"

	"vimana/pb"
)

type RegistryConfig struct {
	Registry Registry `toml:"registry"`
}

type Registry struct {
	URL  string `toml:"url"`
	Port string `toml:"port"`
}

func registryCommands() []*cobra.Command {
	registerCmd := &cobra.Command{
		Use:   "register",
		Short: "Register a spacecore plugin",
		Run: func(cmd *cobra.Command, args []string) {
			logger := log.GetLogger(cmd.Context())
			if len(args) == 0 {
				logger.Info("Usage register <spacecore> <pluginPath>")
			} else {
				name := args[0]
				version := "1.0"
				pluginPath := args[1]
				registerPlugin(name, version, pluginPath)
			}
		},
	}

	searchCmd := &cobra.Command{
		Use:   "search",
		Short: "Search for a spacecore plugin",
		Run: func(cmd *cobra.Command, args []string) {
			logger := log.GetLogger(cmd.Context())
			if len(args) == 0 {
				logger.Info("Usage search <spacecore>\n Listing all spacecore plugins")
				searchPlugins("")
			} else {
				logger.Info("Searching for spacecore plugin: ", args[0])
				searchPlugins(args[0])
			}
		},
	}

	getCmd := &cobra.Command{
		Use:   "get",
		Short: "Get a spacecore plugin",
		Run: func(cmd *cobra.Command, args []string) {
			logger := log.GetLogger(cmd.Context())
			if len(args) == 0 {
				logger.Info("Usage get <spacecore> <version> - default version is 1.0")
			} else {
				name := os.Args[2]
				version := "1.0"
				if len(os.Args) == 4 {
					version = os.Args[3]
				}
				logger.Infof("Getting spacecore plugin %s, %s", name, version)
				pluginPath := getPlugin(name, version)
				logger.Infof("Saving plugin: %s", pluginPath)
				savePlugin(pluginPath, name)
			}
		},
	}
	return []*cobra.Command{registerCmd, searchCmd, getCmd}
}

func getPlugin(name, version string) string {
	logger := log.GetLogger(context.Background())
	conn, client := connect()
	defer conn.Close()

	req := &pb.GetPluginRequest{
		Name:    name,
		Version: version,
	}

	resp, err := client.GetPlugin(context.Background(), req)
	if err != nil {
		logger.Fatalf("Failed to get plugin: %v", err)
	}

	if resp.Plugin == nil {
		logger.Fatalf("Plugin not found")
	}

	fmt.Printf("Plugin: \nName: %s, Version: %s, CID: %s\n", resp.Plugin.Name, resp.Plugin.Version, resp.Plugin.Cid)
	return resp.Plugin.Cid
}
func savePlugin(url, name string) {
	logger := log.GetLogger(context.Background())

	pinataUrl := fmt.Sprintf("https://gateway.pinata.cloud%s", url)
	// logger.Printf("Downloading plugin from %s", pinataUrl)
	resp, err := http.Get(pinataUrl)

	if err != nil {
		logger.Fatalf("Failed to download plugin: %v", err)
	}
	defer resp.Body.Close()

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		logger.Fatalf("Failed to read plugin data: %v", err)
	}

	filePath := filepath.Join("/tmp", name)
	if err := ioutil.WriteFile(filePath, data, 0644); err != nil {
		logger.Fatalf("Failed to save plugin: %v", err)
	}

	logger.Printf("Plugin %s version 0.1 saved to %s", name, filePath)
}

func connect() (*grpc.ClientConn, pb.PluginRegistryClient) {
	logger := log.GetLogger(context.Background())
	// get url and port from config file. TODO Create it if it doesn't exist
	configFile := os.Getenv("HOME") + "/.vimana/registry.toml"
	// set to localhost if running registry locally
	url := os.Getenv("PLUGIN_REGISTRY_URL")
	port := "50051"

	// if url is not set, use default value
	if url == "" {
		url = "registry.vistara.dev"
	}

	if _, err := os.Stat(configFile); os.IsExist(err) {

		var config RegistryConfig
		if _, err := toml.DecodeFile(configFile, &config); err != nil {
			return nil, nil
		}
		url = config.Registry.URL
		port = config.Registry.Port
	}

	target := fmt.Sprintf("%s:%s", url, port)
	conn, err := grpc.NewClient(target, grpc.WithInsecure())
	// logger.Infof("Connecting to gRPC server at %s", target)
	if err != nil {
		logger.Fatalf("Failed to connect to gRPC server: %v", err)
	}

	client := pb.NewPluginRegistryClient(conn)
	return conn, client
}

func registerPlugin(name, version, pluginPath string) {
	logger := log.GetLogger(context.Background())
	logger.Info("Registering plugin", name, version, pluginPath)
	conn, client := connect()
	defer conn.Close()

	req := &pb.RegisterPluginRequest{
		Name:    name,
		Version: version,
		Plugin:  pluginPath,
	}

	resp, err := client.RegisterPlugin(context.Background(), req)
	if err != nil {
		logger.Fatalf("Failed to register plugin: %v", err)
	}

	fmt.Printf("Plugin registration: %s\nCid Gateway: %s", resp.Message, resp.Cid)
}

func searchPlugins(name string) {
	logger := log.GetLogger(context.Background())
	conn, client := connect()
	defer conn.Close()

	req := &pb.DiscoverPluginsRequest{
		Name: &name,
	}
	resp, err := client.DiscoverPlugins(context.Background(), req)
	if err != nil {
		logger.Fatalf("Failed to discover plugins: %v", err)
	}

	for _, plugin := range resp.Plugins {
		fmt.Printf("Name: %s, Version: %s, CID: %s\n", plugin.Name, plugin.Version, plugin.Cid)
	}
}

var CommanderRegistry = map[string]cli.NodeCommander{
	"celestia-light":  cli.NewCelestiaLightCommander("light"),
	"celestia-bridge": cli.NewCelestiaBridgeCommander("bridge"),
	"avail-light":     cli.NewAvailLightCommander("light"),
}
