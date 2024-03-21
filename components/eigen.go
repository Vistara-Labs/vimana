package components

import (
	"os/exec"
)

type EigenComponent struct {
	Root      string
	ConfigDir string
}

func NewEigenComponent(root string, home string, node string) *EigenComponent {
	return &EigenComponent{
		Root:      root,
		ConfigDir: home,
	}
}

func (c *EigenComponent) InitializeConfig() error {
	// lightNodePath := filepath.Join(os.Getenv("HOME"), c.ConfigDir+"/"+c.NodeType+"-node")
	// mkdir -p ~/.vimana/celestia/light-node
	return nil
}

func (c *EigenComponent) GetStartCmd() *exec.Cmd {
	args := []string{}
	// eigenup.sh handles this.
	return exec.Command(
		c.Root, args...,
	)
}
