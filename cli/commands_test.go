package cli

import (
	"fmt"
	"github.com/spf13/cobra"
	"io/ioutil"
	"os"
	"testing"
)

func TestGetCommandsFromConfig(t *testing.T) {
	// Mocking the TOML config data
	mockData := `
[components]

[components.celestia]

[components.celestia.light]
binary = "/usr/local/bin/celestia/celestia"
download = "/tmp/vimana/celestia/init.sh"

[components.celestia.bridge]
binary = "/usr/local/bin/celestia/celestia"
download = "/tmp/vimana/celestia/init.sh"

[components.avail]

[components.avail.light]
binary = "avail-light"
download = "/tmp/vimana/avail/init.sh"


[components.gmworld]

[components.gmworld.da]
binary = "gmworld-da"
download = "/tmp/vimana/gmd/rollup_init.sh"

[components.gmworld.rollup]
binary = "gmworld-rollup"
download = "/tmp/vimana/gmd/rollup_mocha.sh"
`
	// Write mockData to a temporary file
	tmpfile, err := ioutil.TempFile("", "example.toml")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tmpfile.Name()) // cleanup
	if _, err := tmpfile.Write([]byte(mockData)); err != nil {
		t.Fatal(err)
	}
	if err := tmpfile.Close(); err != nil {
		t.Fatal(err)
	}

	// Define a mock commanderRegistry
	var mockCommanderRegistry = map[string]NodeCommander{
		"celestia-light":  NewMockCommander("light"),
		"celestia-bridge": NewMockCommander("bridge"),
		"avail-light":     NewMockCommander("light"),
		"eigen-operator":  NewMockCommander("operator"),
	}

	// Call GetCommandsFromConfig
	commands, err := GetCommandsFromConfig(tmpfile.Name(), mockCommanderRegistry)

	if err != nil {
		t.Fatal(err)
	}

	/*	if len(commands) != 1 {
		print(commands)
		t.Fatalf("Expected 1 main command but got %d", len(commands))
	}*/

	runCmd := commands[0]
	if runCmd.Use != "run" {
		t.Fatalf("Expected 'run' command but got '%s'", runCmd.Use)
	}

	/*
		if len(runCmd.Commands()) != 2 {
			t.Fatalf("Expected 2 component commands but got %d", len(runCmd.Commands()))
		}
	*/
}

type MockCommander struct{ BaseCommander }

func NewMockCommander(node_type string) *MockCommander {
	return &MockCommander{
		BaseCommander: BaseCommander{NodeType: node_type},
	}
}

func (m *MockCommander) Init(cmd *cobra.Command, args []string, mode Mode, node_info string) error {
	fmt.Println("MockCommander.Init() called")
	return nil
}

func (m *MockCommander) Start(cmd *cobra.Command, args []string, mode Mode, node_info string) {
}

func (m *MockCommander) Status(cmd *cobra.Command, args []string, mode Mode, node_info string) {
}

func (m *MockCommander) Stop(cmd *cobra.Command, args []string, mode Mode, node_info string) {
}

func (m *MockCommander) Install(cmd *cobra.Command, args []string, mode Mode, node_info string) {
}

func (m *MockCommander) AddFlags(cmd *cobra.Command) {
}
