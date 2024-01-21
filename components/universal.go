package components

import (
	"os/exec"
)

type UniveralComponent struct {
	Root      string
	ConfigDir string
}

func NewUniveralComponent(root string, home string, node string) *UniveralComponent {
	return &UniveralComponent{
		Root:      root,
		ConfigDir: home,
	}
}

func (c *UniveralComponent) InitializeConfig() error {
	// lightNodePath := filepath.Join(os.Getenv("HOME"), c.ConfigDir+"/"+c.NodeType+"-node")
	// mkdir -p ~/.vimana/celestia/light-node
	return nil
}

func (c *UniveralComponent) GetStartCmd() *exec.Cmd {
	args := []string{}
	// Univeralup.sh handles this.
	return exec.Command(
		c.Root, args...,
	)
}
