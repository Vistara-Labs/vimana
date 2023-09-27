package components

import (
	"os/exec"
)

type AvailComponent struct {
	Root      string
	ConfigDir string
}

func NewAvailComponent(root string, home string, node string) *AvailComponent {
	return &AvailComponent{
		Root:      root,
		ConfigDir: home,
	}
}

func (c *AvailComponent) InitializeConfig() error {
	// lightNodePath := filepath.Join(os.Getenv("HOME"), c.ConfigDir+"/"+c.NodeType+"-node")
	// mkdir -p ~/.vimana/celestia/light-node
	return nil
}

func (c *AvailComponent) GetStartCmd() *exec.Cmd {
	args := []string{}
	// availup.sh handles this.
	return exec.Command(
		c.Root, args...,
	)
}
