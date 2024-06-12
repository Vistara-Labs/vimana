package spacecores

import (
	"os/exec"
)

type EigenSpacecore struct {
	Root      string
	ConfigDir string
}

func NewEigenSpacecore(root string, home string, node string) *EigenSpacecore {
	return &EigenSpacecore{
		Root:      root,
		ConfigDir: home,
	}
}

func (c *EigenSpacecore) InitializeConfig() error {
	// lightNodePath := filepath.Join(os.Getenv("HOME"), c.ConfigDir+"/"+c.NodeType+"-node")
	// mkdir -p ~/.vimana/celestia/light-node
	return nil
}

func (c *EigenSpacecore) GetStartCmd() *exec.Cmd {
	args := []string{}
	// eigenup.sh handles this.
	return exec.Command(
		c.Root, args...,
	)
}
