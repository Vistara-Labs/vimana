package components

import (
	"os/exec"
)

type RollupComponent struct {
	Root      string
	ConfigDir string
}

func NewRollupComponent(root string, home string, node string) *RollupComponent {
	return &RollupComponent{
		Root:      root,
		ConfigDir: home,
	}
}

func (c *RollupComponent) InitializeConfig() error {
	// lightNodePath := filepath.Join(os.Getenv("HOME"), c.ConfigDir+"/"+c.NodeType+"-node")
	// mkdir -p ~/.vimana/celestia/light-node
	return nil
}

func (c *RollupComponent) GetStartCmd() *exec.Cmd {
	args := []string{}
	// Rollupup.sh handles this.
	return exec.Command(
		c.Root, args...,
	)
}
