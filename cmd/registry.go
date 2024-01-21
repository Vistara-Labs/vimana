package cmd

import "vimana/cli"

// CommanderRegistry maps node types to their corresponding NodeCommander implementations.
var CommanderRegistry = map[string]cli.NodeCommander{
	"celestia-light":  cli.NewCelestiaLightCommander("light"),
	"celestia-bridge": cli.NewCelestiaBridgeCommander("bridge"),
	"avail-light":     cli.NewAvailLightCommander("light"),
	"gmworld-da":      cli.NewGmworldDaCommander("da"),
	"gmworld-rollup":  cli.NewGmworldRollupCommander("rollup"),
	"eigen-operator":  cli.NewEigenOperatorCommander("operator"),
}
