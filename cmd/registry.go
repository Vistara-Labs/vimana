package cmd

import "vimana/cli"

// CommanderRegistry maps node types to their corresponding NodeCommander implementations.
var CommanderRegistry = map[string]cli.NodeCommander{
	"celestia-light":  cli.NewCelestiaLightCommander(),
	"celestia-bridge": cli.NewCelestiaBridgeCommander(),
	"avail-light":     cli.NewAvailLightCommander(),
	"gmworld-da":     cli.NewGmworldDaCommander(),
	"gmworld-rollup":     cli.NewGmworldRollupCommander(),
}
