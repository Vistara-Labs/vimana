package spacecores

import (
	"os/exec"
)

type UniversalSpacecore struct {
	Root      string
	ConfigDir string
}

func NewUniversalSpacecore(root string, home string, node string) *UniversalSpacecore {
	return &UniversalSpacecore{
		Root:      root,
		ConfigDir: home,
	}
}

func (c *UniversalSpacecore) InitializeConfig() error {
	// lightNodePath := filepath.Join(os.Getenv("HOME"), c.ConfigDir+"/"+c.NodeType+"-node")
	// mkdir -p ~/.vimana/celestia/light-node
	return nil
}

func (c *UniversalSpacecore) GetStartCmd() *exec.Cmd {
	args := []string{}
	// Univeralup.sh handles this.
	return exec.Command(
		c.Root, args...,
	)
}
