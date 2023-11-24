package components

import (
	"os/exec"
)

type NostreamComponent struct {
	Root      string
	ConfigDir string
}

func NewNostreamComponent(root string, home string, node string) *NostreamComponent {
	return &NostreamComponent{
		Root:      root,
		ConfigDir: home,
	}
}

func (c *NostreamComponent) InitializeConfig() error {
	// lightNodePath := filepath.Join(os.Getenv("HOME"), c.ConfigDir+"/"+c.NodeType+"-node")
	// mkdir -p ~/.vimana/celestia/light-node
	return nil
}

func (c *NostreamComponent) GetStartCmd() *exec.Cmd {
	args := []string{}
	// nostreamup.sh handles this.
	return exec.Command(
		c.Root, args...,
	)
}
