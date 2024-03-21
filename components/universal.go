package components

import (
	"os/exec"
)

type UniversalComponent struct {
	Root      string
	ConfigDir string
}

func NewUniversalComponent(root string, home string, node string) *UniversalComponent {
	return &UniversalComponent{
		Root:      root,
		ConfigDir: home,
	}
}

func (c *UniversalComponent) InitializeConfig() error {
	// lightNodePath := filepath.Join(os.Getenv("HOME"), c.ConfigDir+"/"+c.NodeType+"-node")
	// mkdir -p ~/.vimana/celestia/light-node
	return nil
}

func (c *UniversalComponent) GetStartCmd() *exec.Cmd {
	args := []string{}
	// Univeralup.sh handles this.
	return exec.Command(
		c.Root, args...,
	)
}
