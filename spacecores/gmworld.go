package spacecores

import (
	"os/exec"
)

type GmworldSpacecore struct {
	Root      string
	ConfigDir string
}

func NewGmworldSpacecore(root string, home string, node string) *GmworldSpacecore {
	return &GmworldSpacecore{
		Root:      root,
		ConfigDir: home,
	}
}

func (c *GmworldSpacecore) InitializeConfig() error {
	return nil
}

func (c *GmworldSpacecore) GetStartCmd() *exec.Cmd {
	args := []string{}
	// Gmworldup.sh handles this. TODO: Move the init and start command here.
	return exec.Command(
		c.Root, args...,
	)
}
