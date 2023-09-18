# Vimana CLI 

Vimana CLI is a comprehensive tool designed to simplify the creation and management of different types of nodes, including the data availability layer light node, full node, bridge node, and full nodes for Ethereum-like berachain.

## Table of Contents

- [Installation](#installation)
- [Command Structure](#command-structure)
- [Command API](#command-api)
  - [Create Nodes](#create-nodes)
  - [Start Nodes](#start-nodes)
  - [Stop Nodes](#stop-nodes)
  - [Node Status](#node-status)
- [Support & Feedback](#support--feedback)

## Installation

Install Binary:

`curl -L https://vistara-labs.github.io/vimana/install.sh | bash`

Install from Source:

`make build`

Run celestia light node:

`vimana celestia light start`

## Command Structure

Main command: `vimana`

With this setup, when developers want to support new node types or components, they:

1. Add the configuration to config.toml.
2. Implement the NodeCommander interface for that component and mode.
3. Register their implementation in the commanderRegistry.
This provides a modular and expandable CLI framework.

Available subcommands:

- `init`: Initialize and set up the different types of nodes.
- `start`: Start the nodes after creation.
- `stop`: Stop running nodes.
- `status`: Get the status of a node.

## Command API

### Create Nodes

**Syntax**: 
vimana [NODE_TYPE] [OPTIONS]

**Example**: 
vimana celestia light init

### Start Nodes

**Syntax**: 
vimana celestia [NODE_TYPE] start [OPTIONS]

**Example**: 
vimana celestia light start

### Stop Nodes

**Syntax**: 
vimana celestia [NODE_TYPE] stop [OPTIONS]

**Example**: 
vimana celestia light stop

### Node Status

**Syntax**: 
vimana celestia [NODE_TYPE] status [OPTIONS]

**Example**: 
vimana celestia light status

## Support & Feedback

For any issues, questions, or feedback, please contact *mayur@vistara.dev*.
