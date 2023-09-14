package cmd

import "vimana/vimana/cli"

// CommanderRegistry maps node types to their corresponding NodeCommander implementations.
var CommanderRegistry = map[string]cli.NodeCommander{
	"celestia-light":  &cli.CelestiaLightCommander{},
	"celestia-bridge": &cli.CelestiaBridgeCommander{},
}
