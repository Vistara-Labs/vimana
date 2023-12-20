package components

import (
	"os/exec"
)

type GmworldComponent struct {
	Root      string
	ConfigDir string
}

func NewGmworldComponent(root string, home string, node string) *GmworldComponent {
	return &GmworldComponent{
		Root:      root,
		ConfigDir: home,
	}
}

func (c *GmworldComponent) InitializeConfig() error {
	// lightNodePath := filepath.Join(os.Getenv("HOME"), c.ConfigDir+"/"+c.NodeType+"-node")
	// mkdir -p ~/.vimana/celestia/light-node
	return nil
}

func (c *GmworldComponent) GetStartCmd() *exec.Cmd {
	args := []string{}
	// Gmworldup.sh handles this.
	return exec.Command(
		c.Root, args...,
	)
}
