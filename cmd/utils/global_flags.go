package utils

import (
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
)

func AddGlobalFlags(name string, command *cobra.Command) {
	command.PersistentFlags().StringP(
		name, "", GetVimanaConfig(), "The directory of the vimana config files")
}

var FlagNames = struct {
	Home string
}{
	Home: "home",
}

func GetVimanaConfig() string {
	return filepath.Join(os.Getenv("HOME"), ".vimana")
}
