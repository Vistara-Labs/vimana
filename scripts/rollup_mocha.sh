#!/bin/sh

# Define the paths
gmd_path="$HOME/go/bin/gmd"
gm_path="$HOME/.gm"

# Check if $HOME/go/bin/gmd exists before attempting to remove
if [ -e "$gmd_path" ]; then
    rm -r "$gmd_path"
    echo "Removed $gmd_path"
else
    echo "$gmd_path does not exist"
fi

# Check if $HOME/.gm exists before attempting to remove
if [ -e "$gm_path" ]; then
    rm -rf "$gm_path"
    echo "Removed $gm_path"
else
    echo "$gm_path does not exist"
fi


# Function to install jq if not installed
install_jq() {
  if command -v jq >/dev/null 2>&1; then
    echo "jq is already installed."
  else
    echo "Installing jq..."
    # Add installation command based on your operating system
    # For example, on Debian-based systems (like Ubuntu):
    sudo apt-get update
    sudo apt-get install -y jq
    # Add similar commands for other package managers if needed
  fi
}

# Check for dependencies
install_jq

# Exit if a command fails
# set -e

# set variables for the chain
VALIDATOR_NAME=validator1
CHAIN_ID=gm
KEY_NAME=gm-key
KEY_2_NAME=gm-key-2
CHAINFLAG="--chain-id ${CHAIN_ID}"
TOKEN_AMOUNT="10000000000000000000000000stake"
STAKING_AMOUNT="1000000000stake"

# create a random Namespace ID for your rollup to post blocks to
NAMESPACE=$(openssl rand -hex 8)
echo 'NAMESPACE' $NAMESPACE

# query the DA Layer start height, in this case we are querying
# an RPC endpoint provided by Celestia Labs. The RPC endpoint is
# to allow users to interact with Celestia's core network by querying
# the node's state and broadcasting transactions on the Celestia
# network. This is for mocha, if using another network, change the RPC.
DA_BLOCK_HEIGHT=$(curl  https://rpc-mocha.pops.one/block |jq -r '.result.block.header.height')
echo 'DA_BLOCK_HEIGHT' $DA_BLOCK_HEIGHT

# build the gm chain with Rollkit
ignite chain build

# reset any existing genesis/chain data
gmd tendermint unsafe-reset-all

# initialize the validator with the chain ID you set
gmd init $VALIDATOR_NAME --chain-id $CHAIN_ID

# add keys for key 1 and key 2 to keyring-backend test
gmd keys add $KEY_NAME --keyring-backend test
gmd keys add $KEY_2_NAME --keyring-backend test

# add these as genesis accounts
gmd add-genesis-account $KEY_NAME $TOKEN_AMOUNT --keyring-backend test
gmd add-genesis-account $KEY_2_NAME $TOKEN_AMOUNT --keyring-backend test

# set the staking amounts in the genesis transaction
gmd gentx $KEY_NAME $STAKING_AMOUNT --chain-id $CHAIN_ID --keyring-backend test

# collect genesis transactions
gmd collect-gentxs

# copy centralized sequencer address into genesis.json
# Note: validator and sequencer are used interchangeably here
ADDRESS=$(sudo jq -r '.address' ~/.gm/config/priv_validator_key.json)
PUB_KEY=$(sudo jq -r '.pub_key' ~/.gm/config/priv_validator_key.json)
jq --argjson pubKey "$PUB_KEY" '. + {"validators": [{"address": "'$ADDRESS'", "pub_key": $pubKey, "power": "1000", "name": "Rollkit Sequencer"}]}' ~/.gm/config/genesis.json > temp.json && mv temp.json ~/.gm/config/genesis.json

# export the Celestia light node's auth token to allow you to submit
# PayForBlobs to Celestia's data availability network
# this is for mocha, if using another network, change the network name
export AUTH_TOKEN=$(celestia light auth write --p2p.network mocha)

# start the chain
# gmd start --rollkit.aggregator true --rollkit.lazy_aggregator --rollkit.da_layer celestia --rollkit.da_config='{"base_url":"http://localhost:26658","timeout":60000000000,"fee":600000,"gas_limit":6000000,"auth_token":"'$AUTH_TOKEN'"}' --rollkit.namespace_id $NAMESPACE --rollkit.da_start_height $DA_BLOCK_HEIGHT
gmd start --rollkit.aggregator true  --rollkit.da_layer celestia --rollkit.da_config='{"base_url":"http://localhost:26658","timeout":60000000000,"fee":600000,"gas_limit":6000000,"auth_token":"'$AUTH_TOKEN'"}' --rollkit.namespace_id $NAMESPACE --rollkit.da_start_height $DA_BLOCK_HEIGHT