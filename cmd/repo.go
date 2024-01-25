package cmd

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"vimana/cli"
	"vimana/cmd/utils"

	"github.com/BurntSushi/toml"
	"github.com/spf13/cobra"
)

func repoCommand() *cobra.Command {
	repoCmd := &cobra.Command{
		Use:   "repo",
		Short: "add repo to vimana",
	}

	addCmd := &cobra.Command{
		Use:   "add",
		Short: "add x.y spacecore to vimana",
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) == 0 {
				fmt.Println("lack of parmater")
			} else {
				fmt.Println(args)
				param := strings.Split(args[0], ".")
				component := param[0]
				param = strings.Split(param[1], "-")
				node_type := param[0]
				fmt.Println(component, node_type)

				// prompt user input repo url
				prompter := utils.NewPrompter()
				repo_url, err := prompter.InputString(
					"Enter your github repo address:",
					"",
					"",
					func(s string) error {
						return nil
					},
				)
				if err != nil {
					return
				}
				fmt.Println(repo_url)

				// prompt user input node binary
				binary_file, err := prompter.InputString(
					"Enter your binary file name:",
					"",
					"",
					func(s string) error {
						return nil
					},
				)
				if err != nil {
					return
				}
				fmt.Println(binary_file)

				configFile := os.Getenv("HOME") + "/.vimana/config.toml"
				var config cli.Config
				if _, err := toml.DecodeFile(configFile, &config); err != nil {
					return
				}
				for component, nodeTypes := range config.Components {
					fmt.Println(component)
					for nodeType := range nodeTypes {
						fmt.Println(nodeType, nodeTypes[nodeType])
						fmt.Println(nodeTypes[nodeType].Binary)
						fmt.Println(nodeTypes[nodeType].Download)
					}
				}
				if _, ok := config.Components[component]; !ok {
					var m cli.Mode
					m.Binary = "/usr/local/bin/" + component + "/" + binary_file
					m.Download = "/tmp/vimana/" + component + "/init.sh"

					res, err := http.Get(repo_url + "/init.sh")
					if err != nil {
						fmt.Errorf("file init.sh download error, check file address: %v", err)
						return
					}
					os.MkdirAll("/tmp/vimana/"+component, 0755)
					f, err := os.Create(m.Download)
					if err != nil {
						fmt.Println(f, err)
						return
					}
					_, err = io.Copy(f, res.Body)
					if err != nil {
						fmt.Errorf("file save error: %v", err)
						return
					}
					//download start.sh
					res, err = http.Get(repo_url + "/start.sh")
					if err != nil {
						fmt.Errorf("file start.sh download error, check file address: %v", err)
						return
					}
					os.MkdirAll("/tmp/vimana/"+component, 0755)
					f, err = os.Create("/tmp/vimana/" + component + "/start.sh")
					_, err = io.Copy(f, res.Body)

					//download stop.sh
					res, err = http.Get(repo_url + "/install.sh")
					if err != nil {
						fmt.Errorf("file start.sh download error, check file address: %v", err)
						return
					}
					os.MkdirAll("/tmp/vimana/"+component, 0755)
					f, err = os.Create("/tmp/vimana/" + component + "/install.sh")
					_, err = io.Copy(f, res.Body)

					//download Binary
					//res, err = http.Get(repo_url + binary_file)
					//if err != nil {
					//	fmt.Errorf("file start.sh download error, check file address: %v", err)
					//	return
					//}
					//os.MkdirAll("/usr/local/bin/" + component + "/", 0755)
					//f, err = os.Create("/usr/local/bin/" + component + "/" + binary_file)
					//_, err = io.Copy(f, res.Body)

					new_component := make(map[string]cli.Mode, 1)
					new_component[node_type] = m
					config.Components[component] = new_component

				}

				for component, nodeTypes := range config.Components {
					fmt.Println(component)
					for nodeType := range nodeTypes {
						fmt.Println(nodeType, nodeTypes[nodeType])
						fmt.Println(nodeTypes[nodeType].Binary)
						fmt.Println(nodeTypes[nodeType].Download)
					}
				}

				cli.WriteConf(config)
			}

		},
	}

	repoCmd.AddCommand(addCmd)

	return repoCmd
}
