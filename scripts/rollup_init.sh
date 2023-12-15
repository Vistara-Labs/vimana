#!/bin/bash
set -e

# init light node
# start light node
celestia light init --p2p.network mocha
celestia light start --core.ip rpc-mocha.pops.one --p2p.network mocha
