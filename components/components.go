package components

import (
	"os/exec"
	"vimana/config"
)

type Component interface {
	InitializeConfig() error
	GetStartCmd() *exec.Cmd
}

type ComponentManager struct {
	ComponentType config.ComponentType
	Component
}

type ComponentConfig struct {
	RPC     string
	Network string
}

func NewComponentManager(componentType config.ComponentType, root string, nodeType string, c *ComponentConfig) *ComponentManager {
	var component Component

	switch componentType {
	case config.Celestia:
		component = NewCelestiaComponent(root, ".vimana/celestia", nodeType, c.RPC, c.Network)
	case config.Avail:
		component = NewAvailComponent(root, ".vimana/avail", nodeType)
	// case config.Berachain:
	// 	component = berachain.NewBerachainComponent(home)
	case config.Eigen:
		component = NewEigenComponent(root, ".vimana/eigen", nodeType)
	default:
		panic("Unknown component type")
	}

	return &ComponentManager{
		ComponentType: componentType,
		Component:     component,
	}
}
