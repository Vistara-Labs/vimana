package components

import (
	"os/exec"
)

type GmworldComponent struct {
	Root      string
	ConfigDir string
}

func NewGmworldComponent(root string, home string, node string) *GmworldComponent {
	return &GmworldComponent{
		Root:      root,
		ConfigDir: home,
	}
}

func (c *GmworldComponent) InitializeConfig() error {
	return nil
}

func (c *GmworldComponent) GetStartCmd() *exec.Cmd {
	args := []string{}
	// Gmworldup.sh handles this. TODO: Move the init and start command here.
	return exec.Command(
		c.Root, args...,
	)
}
