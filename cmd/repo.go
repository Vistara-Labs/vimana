package cmd

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"vimana/cli"
	"vimana/cmd/utils"
	"vimana/log"

	"github.com/BurntSushi/toml"

	"github.com/spf13/cobra"
)

func repoCommand() *cobra.Command {

	repoCmd := &cobra.Command{
		Use:   "repo",
		Short: "Add repo to vimana",
	}

	addCmd := &cobra.Command{
		Use:   "add",
		Short: "add x.y spacecore to vimana",
		Run: func(cmd *cobra.Command, args []string) {

			logger := log.GetLogger(cmd.Context())
			if len(args) == 0 {
				logger.Infof(
					"Please provide the spacecore and node type to add. e.g. vimana repo add x.y-node-type",
				)
			} else {
				logger.Info(args)
				param := strings.Split(args[0], ".")
				spacecore := param[0]
				param = strings.Split(param[1], "-")
				node_type := param[0]
				logger.Info(spacecore, node_type)

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
				logger.Info(repo_url)

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
				logger.Info(binary_file)

				configFile := os.Getenv("HOME") + "/.vimana/config.toml"
				var config cli.Config
				if _, err := toml.DecodeFile(configFile, &config); err != nil {
					return
				}
				for spacecore, nodeTypes := range config.Spacecores {
					logger.Info(spacecore)
					for nodeType := range nodeTypes {
						logger.Info(nodeType, nodeTypes[nodeType])
						logger.Info(nodeTypes[nodeType].Binary)
						logger.Info(nodeTypes[nodeType].Download)
					}
				}
				if _, ok := config.Spacecores[spacecore]; !ok {
					var m cli.Mode
					m.Binary = "/usr/local/bin/" + spacecore + "/" + binary_file
					m.Download = "/tmp/vimana/" + spacecore + "/init.sh"
					m.Install = "/tmp/vimana/" + spacecore + "/install.sh"
					m.Start = "/usr/local/bin/" + spacecore + "/start.sh"

					logger.Infof("m mode is  %v\n\n", m)
					res, err := http.Get(repo_url + "/init.sh")
					if err != nil {
						fmt.Errorf("file init.sh download error, check file address: %v", err)
						return
					}
					os.MkdirAll("/tmp/vimana/"+spacecore, 0755)
					f, err := os.Create(m.Download)
					if err != nil {
						logger.Info(f, err)
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
					os.MkdirAll("/usr/local/bin/"+spacecore, 0755)
					f, err = os.Create(m.Start)
					_, err = io.Copy(f, res.Body)

					//download stop.sh
					res, err = http.Get(repo_url + "/install.sh")
					if err != nil {
						fmt.Errorf("file start.sh download error, check file address: %v", err)
						return
					}
					os.MkdirAll("/tmp/vimana/"+spacecore, 0755)
					f, err = os.Create(m.Install)
					_, err = io.Copy(f, res.Body)

					//download Binary
					//res, err = http.Get(repo_url + binary_file)
					//if err != nil {
					//	fmt.Errorf("file start.sh download error, check file address: %v", err)
					//	return
					//}
					//os.MkdirAll("/usr/local/bin/" + spacecore + "/", 0755)
					//f, err = os.Create("/usr/local/bin/" + spacecore + "/" + binary_file)
					//_, err = io.Copy(f, res.Body)

					new_spacecore := make(map[string]cli.Mode, 1)
					new_spacecore[node_type] = m
					config.Spacecores[spacecore] = new_spacecore

				}

				for spacecore, nodeTypes := range config.Spacecores {
					logger.Infof("spacecore %s\n\n", spacecore)
					for nodeType := range nodeTypes {
						logger.Info(nodeType, nodeTypes[nodeType])
						logger.Info(nodeTypes[nodeType].Binary)
						logger.Info(nodeTypes[nodeType].Download)
					}
				}

				cli.WriteConf(config)
			}

		},
	}

	repoCmd.AddCommand(addCmd)

	// importCmd := &cobra.Command{
	// 	Use:   "import",
	// 	Short: "import repo from vimana",
	// 	Run: func(cmd *cobra.Command, args []string) {
	// 		logger := log.GetLogger(cmd.Context())
	// 		logger.Info("import repo")
	// 		// prompt user input repo url

	// 		prompter := utils.NewPrompter()
	// 		repo_url, err := prompter.InputString(
	// 			"Enter your github repo address:",
	// 			"",
	// 			"",
	// 			func(s string) error {
	// 				return nil
	// 			},
	// 		)
	// 		if err != nil {
	// 			return
	// 		}
	// 		logger.Info(repo_url)

	// 		// Download the spaceocre from the source url
	// 		// Verify the downloaded content (checksum)
	// 		// Register the spacecore in the Vimana system
	// 		// Add the spacecore to the config.toml file
	// 		// Return a success message to the user and any errors

	// 		// this is the approach used for vimana repo add cmd.
	// 		// res, err := http.Get(repo_url + "/init.sh")
	// 		// if err != nil {
	// 		// 	fmt.Errorf("file init.sh download error, check file address: %v", err)
	// 		// 	return
	// 		// }
	// 		// os.MkdirAll("/tmp/vimana/"+repo_url, 0755)

	// 	},
	// }

	// repoCmd.AddCommand(importCmd)
	return repoCmd
}
