package components

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
)

// TODO: test how much is enough to run the LC for one day and set the minimum balance accordingly.
const (
	CelestiaRestApiEndpoint = "https://api-arabica-9.consensus.celestia-arabica.com"
	DefaultCelestiaRPC      = "consensus-full-arabica-9.celestia-arabica.com"
	DefaultCelestiaNetwork  = "arabica"
)

type CelestiaComponent struct {
	Root            string
	ConfigDir       string
	rpcEndpoint     string
	metricsEndpoint string
	RPCPort         string
	NamespaceID     string
}

func NewCelestiaComponent(root string, home string) *CelestiaComponent {
	return &CelestiaComponent{
		Root:      root,
		ConfigDir: home,
	}
}

func (c *CelestiaComponent) InitializeConfig() error {
	lightNodePath := filepath.Join(os.Getenv("HOME"), c.ConfigDir+"/light-node")

	path, err := filepath.Abs(filepath.Join(lightNodePath + "/config.toml"))
	if err != nil {
		fmt.Println("Error:", err)
		return err
	}
	if _, err := os.Stat(path); os.IsNotExist(err) {
		initLightNodeCmd := exec.Command(c.Root, "light", "init", "--p2p.network",
			DefaultCelestiaNetwork, "--node.store", lightNodePath)
		err := initLightNodeCmd.Run()
		fmt.Println("🚀 initLightNodeCmd", initLightNodeCmd)
		if err != nil {
			fmt.Println("Error initializing light node config", err)
			return err
		}
		fmt.Println("🚀 Celestia light node initialized: ", path)
	} else {
		fmt.Println("🚀 Celestia light node already initialized: ", path)
	}
	return nil
}

func (c *CelestiaComponent) GetStartCmd() *exec.Cmd {
	lightNodePath := filepath.Join(os.Getenv("HOME"), c.ConfigDir+"/light-node")
	args := []string{
		"light", "start",
		"--core.ip", c.rpcEndpoint,
		"--node.store", lightNodePath,
		"--gateway",
		"--gateway.deprecated-endpoints",
		"--p2p.network", DefaultCelestiaNetwork,
	}
	if c.metricsEndpoint != "" {
		args = append(args, "--metrics", "--metrics.endpoint", c.metricsEndpoint)
	}
	return exec.Command(
		c.Root, args...,
	)
}
