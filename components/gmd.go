package components

import (
	"os/exec"
)

type GmdComponent struct {
	Root      string
	ConfigDir string
}

func NewGmdComponent(root string, home string, node string) *GmdComponent {
	return &GmdComponent{
		Root:      root,
		ConfigDir: home,
	}
}

func (c *GmdComponent) InitializeConfig() error {
	// lightNodePath := filepath.Join(os.Getenv("HOME"), c.ConfigDir+"/"+c.NodeType+"-node")
	// mkdir -p ~/.vimana/celestia/light-node
	return nil
}

func (c *GmdComponent) GetStartCmd() *exec.Cmd {
	args := []string{}
	// Gmdup.sh handles this.
	return exec.Command(
		c.Root, args...,
	)
}
