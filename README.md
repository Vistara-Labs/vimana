# Vimana CLI

Vimana CLI is a comprehensive tool designed to simplify the creation and management of different types of nodes, including the data availability layer light node, full node, bridge node, and full nodes for Ethereum-like berachain.

## Table of Contents

- [Installation](#installation)
- [Command Structure](#command-structure)
- [Command API](#command-api)
  - [Run Nodes](#run-nodes)
  - [Stop Nodes](#stop-nodes)
  - [Node Status](#node-status)
- [Support & Feedback](#support--feedback)

## Quick Start

Open terminal and run:

```bash
  curl -L https://vistara-labs.github.io/vimana/install.sh | bash && vimana init
  vimana run celestia light-node
```

## Installation

Install Binary:

`curl -L https://vistara-labs.github.io/vimana/install.sh | bash`

Run celestia light node:

`vimana run celestia light-node`

Install from Source:

For Linux:

`make build-linux`

For MacOS:

`make build`

Run celestia light node:

`./vimana-linux-amd64/vimana run celestia light-node`

See options for a specific node type:

`vimana run celestia light-node --help`

```
Usage:
  vimana run celestia bridge-node [flags]

Flags:
  -h, --help             help for bridge-node
      --network string   Specifies the Celestia network (default "arabica")
      --rpc string       Specifies the Celestia RPC endpoint (default "consensus-validator.celestia-arabica-10.com")
```

You can pass in the network and rpc endpoint as flags, or default (arabica) is used if not specified.

`vimana run celestia light-node --network arabica --rpc consensus-validator.celestia-arabica-10.com`

## Command Structure

Main command: `vimana`

Subcommand:

- `run`: Initialize and run the different types of nodes.

With this setup, when developers want to support new node types or components, they:

1. Add the configuration to config.toml.
2. Implement the NodeCommander interface for that component and mode.
3. Register their implementation in the commanderRegistry.
   This provides a modular and expandable CLI framework.

## Command API

## Run Nodes

**Syntax**:
vimana run [NODE_TYPE] [OPTIONS]

**Example**:

```
vimana run celestia light-node
vimana run celestia bridge-node
```

**Launch via service**: <br/>
service creation for light-node

```
tee /etc/systemd/system/vinama.service > /dev/null <<EOF
[Unit]
Description=Vinama
After=network.target
[Service]
Type=simple
User=$USER
ExecStart=vimana run celestia light-node
Restart=on-failure
RestartSec=10
LimitNOFILE=65535
[Install]
WantedBy=multi-user.target
EOF
```

launch

```
systemctl daemon-reload
systemctl enable vinama
systemctl restart vinama && journalctl -u vinama -f -o cat
```

service creation for bridge-node

```
tee /etc/systemd/system/vinama-bridge.service > /dev/null <<EOF
[Unit]
Description=Vinama
After=network.target
[Service]
Type=simple
User=$USER
ExecStart=vimana run celestia bridge-node
Restart=on-failure
RestartSec=10
LimitNOFILE=65535
[Install]
WantedBy=multi-user.target
EOF
```

launch

```
systemctl daemon-reload
systemctl enable vinama-bridge
systemctl restart vinama-bridge && journalctl -u vinama-bridge -f -o cat
```

### Create a new component, avail

Follow #CREATE_COMPONENT.md

```
vimana run avail light-node
```

## Stop Nodes

```
sudo systemctl stop vinama
```

## Node Status

```
sudo systemctl status vinama
```

### Node log

```
sudo journalctl -u vinama.service -f

```

## Support & Feedback

For any issues, questions, or feedback, please contact *mayur@vistara.dev*.
