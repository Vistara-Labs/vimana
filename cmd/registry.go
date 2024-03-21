package cmd

import (
	"fmt"
	"os"
	"strings"
	"vimana/cli"

	"github.com/BurntSushi/toml"
	"github.com/asmcos/requests"
	"github.com/spf13/cobra"
)

// CommanderRegistry maps node types to their corresponding NodeCommander implementations.
var CommanderRegistry = map[string]cli.NodeCommander{
	"celestia-light":  cli.NewCelestiaLightCommander("light"),
	"celestia-bridge": cli.NewCelestiaBridgeCommander("bridge"),
	"avail-light":     cli.NewAvailLightCommander("light"),
	"gmworld-da":      cli.NewGmworldDaCommander("da"),
	"gmworld-rollup":  cli.NewGmworldRollupCommander("rollup"),
	"eigen-operator":  cli.NewEigenOperatorCommander("operator"),
}

func registryCommand() *cobra.Command {
	type spacecore struct {
		Spacecore string `json:"spacecore"`
		Repo      string `json:"repo"`
	}

	//url := "https://raw.githubusercontent.com/Vistara-Labs/vimana/spacecores.json"
	url := "https://raw.githubusercontent.com/zhangwenqiangnb/vimana/dev/spacecores.json"

	registryCmd := &cobra.Command{
		Use:   "registry",
		Short: "registry search/list command",
	}

	searchCmd := &cobra.Command{
		Use:   "search",
		Short: "search x",
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) == 0 {
				fmt.Println("lack of parmater")
			} else {
				// 1. check config.toml
				configFile := os.Getenv("HOME") + "/.vimana/config.toml"
				var config cli.Config
				if _, err := toml.DecodeFile(configFile, &config); err != nil {
					return
				}
				for component := range config.Components {
					if strings.Contains(strings.ToLower(component), strings.ToLower(args[0])) {
						fmt.Println(component)
					}
				}
				// 2. check spacecores.json
				resp, err := requests.Get(url)
				if err != nil {
					return
				}
				// Status code
				if resp.R.StatusCode != 200 {
					fmt.Println("Get Vistara-Labs/vimana/spacecores.json error: status code =", resp.R.StatusCode)
					return
				}

				var Spacecores []spacecore
				resp.Json(&Spacecores)

				for _, spacecore := range Spacecores {
					if strings.Contains(strings.ToLower(spacecore.Spacecore), strings.ToLower(args[0])) {
						fmt.Println(spacecore.Spacecore)
					}
				}
			}

		},
	}

	listCmd := &cobra.Command{
		Use:   "list",
		Short: "list",
		Run: func(cmd *cobra.Command, args []string) {
			// 1. check config.toml
			configFile := os.Getenv("HOME") + "/.vimana/config.toml"
			var config cli.Config
			if _, err := toml.DecodeFile(configFile, &config); err != nil {
				return
			}
			for component := range config.Components {
				fmt.Println(component)
			}

			// 2. check spacecores.json
			resp, err := requests.Get(url)
			if err != nil {
				return
			}
			// Status code
			if resp.R.StatusCode != 200 {
				fmt.Println("Get Vistara-Labs/vimana/spacecores.json error: status code =", resp.R.StatusCode)
				return
			}

			var Spacecores []spacecore
			resp.Json(&Spacecores)

			for _, spacecore := range Spacecores {
				fmt.Println(spacecore.Spacecore)
			}

		},
	}

	registryCmd.AddCommand(searchCmd)
	registryCmd.AddCommand(listCmd)

	return registryCmd
}
