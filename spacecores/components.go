package spacecores

import (
	"os/exec"
	"vimana/config"
)

type Spacecore interface {
	InitializeConfig() error
	GetStartCmd() *exec.Cmd
}

type SpacecoreManager struct {
	SpacecoreType config.SpacecoreType
	Spacecore
}

type SpacecoreConfig struct {
	RPC     string
	Network string
}

func NewSpacecoreManager(spacecoreType config.SpacecoreType, root string, nodeType string, c *SpacecoreConfig) *SpacecoreManager {
	var spacecore Spacecore

	switch spacecoreType {
	case config.Celestia:
		spacecore = NewCelestiaSpacecore(root, ".vimana/celestia", nodeType, c.RPC, c.Network)
	case config.Avail:
		spacecore = NewAvailSpacecore(root, ".vimana/avail", nodeType)
	case config.Gmworld:
		spacecore = NewGmworldSpacecore(root, ".vimana/gmd", nodeType)
	// case config.Berachain:
	// 	spacecore = berachain.NewBerachainSpacecore(home)
	case config.Eigen:
		spacecore = NewEigenSpacecore(root, ".vimana/eigen", nodeType)
	default:
		//panic("Unknown spacecore type")
		spacecore = NewUniversalSpacecore(root, ".vimana/"+string(spacecoreType), nodeType)
	}

	return &SpacecoreManager{
		SpacecoreType: spacecoreType,
		Spacecore:     spacecore,
	}
}
