package components

import (
	"log"
	"os"
	"os/exec"
	"path/filepath"
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

func NewCelestiaComponent(root string, home string, node string, celestiaRPC, celestiaNetwork string) *CelestiaComponent {
	return &CelestiaComponent{
		Root:            root,
		ConfigDir:       home,
		NodeType:        node,
		NodeStorePath:   filepath.Join(os.Getenv("HOME"), home+"/"+node+"-node"+"/"+celestiaNetwork),
		rpcEndpoint:     celestiaRPC,
		celestiaNetwork: celestiaNetwork,
	}
}

func (c *CelestiaComponent) InitializeConfig() error {
	log.Println("🚀 Creating Celestia ", c.NodeType, " node config dir: ", c.NodeStorePath)
	if _, err := os.Stat(c.NodeStorePath); os.IsNotExist(err) {
		err := os.MkdirAll(c.NodeStorePath, 0755)
		if err != nil {
			log.Println("Error creating ", c.NodeType, " node config dir", err)
			return err
		}
		log.Println("🚀 Celestia ", c.NodeType, " node config dir created: ", c.NodeStorePath)
	} else {
		log.Println("🚀 Celestia ", c.NodeType, " node config dir already exists: ", c.NodeStorePath)
	}

	path, err := filepath.Abs(filepath.Join(c.NodeStorePath + "/config.toml"))
	if err != nil {
		log.Println("Error:", err)
		return err
	}
	if _, err := os.Stat(path); os.IsNotExist(err) {
		initLightNodeCmd := exec.Command(c.Root, c.NodeType, "init", "--p2p.network",
			c.celestiaNetwork, "--node.store", c.NodeStorePath)
		err := initLightNodeCmd.Run()
		log.Println("🚀 initLightNodeCmd", initLightNodeCmd)
		if err != nil {
			log.Println("Error initializing light node config", err)
			return err
		}
		log.Println("🚀 Celestia light node initialized: ", path)
	} else {
		log.Println("🚀 Celestia light node already initialized: ", path)
	}
	return nil
}

func (c *CelestiaComponent) GetStartCmd() *exec.Cmd {
	args := []string{
		c.NodeType, "start",
		"--core.ip", c.rpcEndpoint,
		"--node.store", c.NodeStorePath,
		"--gateway",
		// "--gateway.deprecated-endpoints",
		"--p2p.network", c.celestiaNetwork,
	}
	if c.metricsEndpoint != "" {
		args = append(args, "--metrics", "--metrics.endpoint", c.metricsEndpoint)
	}
	return exec.Command(
		c.Root, args...,
	)
}
