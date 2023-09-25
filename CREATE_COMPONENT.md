# Vimana Components 

Vimana CLI is a comprehensive tool designed to simplify the creation and management of different types of nodes, including the data availability layer light node, full node, bridge node, and full nodes for Ethereum-like berachain.

How to add a new component to Vimana CLI:

1. Add the configuration to config.toml.

e.g. 
```
[components]
[components.dymension]

[components.dymension.full]
binary = "/tmp/vimana/dymension/dymd"
download = "/tmp/vimana/dymension/init.sh"
```

2. Implement the NodeCommander interface for that component and mode.

e.g. 
```
type DymensionFullCommander struct {
  *BaseCommander
}

func NewDymensionFullCommander() *DymensionFullCommander {
	return &DymensionFullCommander{
		BaseCommander: BaseCommander{NodeType: "full"},
	}
}

func (c *DymensionFullCommander) Init(cmd *cobra.Command, args []string, mode Mode) error {
	utils.ExecBashCmd(exec.Command("bash", mode.Download), utils.WithOutputToStdout(), utils.WithErrorsToStderr())
	c.initComponentManager(mode.Binary)
	return c.componentMgr.InitializeConfig()
}

func (c *DymensionFullCommander) Run(cmd *cobra.Command, args []string, mode Mode) {
	c.Init(cmd, args, mode)
	cmdexecute := c.componentMgr.GetStartCmd()
	utils.ExecBashCmd(cmdexecute, utils.WithOutputToStdout(), utils.WithErrorsToStderr())
}
```

3. Update initComponentManager() and NewComponentManager() to support the new component.

e.g. 
```
func (c *BaseCommander) initComponentManager(binary string) {
	if b.componentMgr == nil {
		b.componentMgr = components.NewComponentManager("dymension", binary, b.NodeType)
	}
}
// In NewComponentManager():
	case config.Celestia:
		component = NewDymensionComponent(root, ".vimana/dymension", nodeType)
```

4. Implement NewDymensionComponent and it's methods InitializeConfig() and GetStartCmd().

e.g. 
```
// implement how to initialize dymd
// implement the start command for dymd
```

5. Register their implementation in the commanderRegistry.

e.g. 
```
commanderRegistry.Register("dymension", "full", &DymensionFullCommander{})
	"celestia-light":  cli.NewCelestiaLightCommander(),
```

6. Add the component to the README.md.