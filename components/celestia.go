package components

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
)

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
	NodeType        string
	NodeStorePath   string
	celestiaNetwork string
}

func NewCelestiaComponent(root string, home string, node string) *CelestiaComponent {
	return &CelestiaComponent{
		Root:            root,
		ConfigDir:       home,
		NodeType:        node,
		NodeStorePath:   filepath.Join(os.Getenv("HOME"), home+"/"+node+"-node"),
		rpcEndpoint:     DefaultCelestiaRPC,
		celestiaNetwork: DefaultCelestiaNetwork,
	}
}

func (c *CelestiaComponent) InitializeConfig() error {
	// lightNodePath := filepath.Join(os.Getenv("HOME"), c.ConfigDir+"/"+c.NodeType+"-node")
	// mkdir -p ~/.vimana/celestia/light-node
	fmt.Println("🚀 Creating Celestia ", c.NodeType, " node config dir: ", c.NodeStorePath)
	if _, err := os.Stat(c.NodeStorePath); os.IsNotExist(err) {
		err := os.MkdirAll(c.NodeStorePath, 0755)
		if err != nil {
			fmt.Println("Error creating ", c.NodeType, " node config dir", err)
			return err
		}
		fmt.Println("🚀 Celestia ", c.NodeType, " node config dir created: ", c.NodeStorePath)
	} else {
		fmt.Println("🚀 Celestia ", c.NodeType, " node config dir already exists: ", c.NodeStorePath)
	}

	path, err := filepath.Abs(filepath.Join(c.NodeStorePath + "/config.toml"))
	if err != nil {
		fmt.Println("Error:", err)
		return err
	}
	if _, err := os.Stat(path); os.IsNotExist(err) {
		initLightNodeCmd := exec.Command(c.Root, c.NodeType, "init", "--p2p.network",
			DefaultCelestiaNetwork, "--node.store", c.NodeStorePath)
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
	// nodePath := filepath.Join(os.Getenv("HOME"), c.ConfigDir+"/"+c.NodeType+"/light-node")
	args := []string{
		c.NodeType, "start",
		"--core.ip", c.rpcEndpoint,
		"--node.store", c.NodeStorePath,
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
