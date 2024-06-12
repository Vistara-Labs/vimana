package spacecores

import (
	"os/exec"
)

type AvailSpacecore struct {
	Root      string
	ConfigDir string
}

func NewAvailSpacecore(root string, home string, node string) *AvailSpacecore {
	return &AvailSpacecore{
		Root:      root,
		ConfigDir: home,
	}
}

func (c *AvailSpacecore) InitializeConfig() error {
	// lightNodePath := filepath.Join(os.Getenv("HOME"), c.ConfigDir+"/"+c.NodeType+"-node")
	// mkdir -p ~/.vimana/celestia/light-node
	return nil
}

func (c *AvailSpacecore) GetStartCmd() *exec.Cmd {
	args := []string{}
	// availup.sh handles this.
	return exec.Command(
		c.Root, args...,
	)
}
