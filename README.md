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
  - [List Nodes](#list-nodes)
- [Flag Details](#flag-details)
- [Potential Enhancements](#potential-enhancements)
- [Support & Feedback](#support--feedback)

## Installation

*Instructions for installation will go here.*

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
- `list`: List all nodes of a particular type.

## Command API

### Create Nodes

**Syntax**: 
vimana create [NODE_TYPE] [OPTIONS]

**Example**: 
vimana create light-node --port=30303 --datadir=/path/to/data

### Start Nodes

**Syntax**: 
vimana start [NODE_TYPE] [OPTIONS]

**Example**: 
vimana start light-node --id=node123

### Stop Nodes

**Syntax**: 
vimana stop [NODE_TYPE] [OPTIONS]

**Example**: 
vimana stop full-node --id=node456

### Node Status

**Syntax**: 
vimana status [NODE_TYPE] [OPTIONS]

**Example**: 
vimana status bridge-node --id=node789

### List Nodes

**Syntax**: 
vimana list [NODE_TYPE]

**Example**: 
vimana list beechain


## Flag Details

- `--port`: Specify the port for the node.
- `--datadir`: Specify the directory for node's data.
- `--id`: Identify a specific node instance.

## Potential Enhancements

1. **Config Files**: Consider a configuration file approach (`.yaml` or `.json`) for detailing settings.
2. **Logging**: Implement a `--log` or `--verbose` flag for detailed operations logs.
3. **Node Updates**: Introduce functionality to update or modify node configurations post-creation.

## Support & Feedback

For any issues, questions, or feedback, please contact *support-email@example.com*.
