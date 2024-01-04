#!/bin/sh

echo "Fund your celestia account"
# cel-key list --node.type light --keyring-backend test --p2p.network mocha

## install the gmd binary
INTERNAL_DIR="/usr/local/bin"

# check if the binary is already installed
if [ -f "$INTERNAL_DIR/gmd" ]; then
    # Capture the version output
    VERSION_OUTPUT=$("$INTERNAL_DIR/gmd" version)
    
    # Check if the version matches "v0.11.0-rc15-dev"
    echo "🚀 gmd is already installed with the correct version." $VERSION_OUTPUT
else
    echo "🚀 gmd is not installed..."
    echo "🔍 Determining OS and architecture..."
    OS=$(uname -s | tr '[:upper:]' '[:lower:]')
    ARCH=$(uname -m)

    if [[ "$ARCH" == "x86_64" ]]; then
        ARCH="amd64"
    elif [[ "$ARCH" == "arm64" ]] || [[ "$ARCH" == "aarch64" ]]; then
        ARCH="arm64"
    fi

    echo "💻  OS: $OS"
    echo "🏗️  ARCH: $ARCH"

    # if OS is linux then install unzip
    if [[ "$OS" == "linux" ]]; then
        if which apt > /dev/null; then
            sudo apt-get update > /dev/null
            sudo apt-get install unzip > /dev/null
        elif which apk > /dev/null; then
            sudo apk update > /dev/null
            sudo apk add unzip > /dev/null
            ARCH="arm64_alpine"
        fi
    fi

    TGZ_URL="https://github.com/Vistara-Labs/vimana/releases/download/gmd-v0.0.1/gmd-${OS}-${ARCH}.zip"
    sudo mkdir -p "/tmp/vimgmd"
    echo "💈 Downloading gmd..."
    sudo curl -o /tmp/vimgmd/gmd-${OS}-${ARCH}.zip -L "$TGZ_URL" --progress-bar

    sudo unzip -q /tmp/vimgmd/gmd-${OS}-${ARCH}.zip -d /tmp/vimgmd/
    sudo mv "/tmp/vimgmd/gmd-${OS}-${ARCH}"/* "$INTERNAL_DIR"
    # sudo chmod +x "$INTERNAL_DIR"
    sudo rm -rf "/tmp/vimgmd"
fi



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
# ignite chain build

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
export AUTH_TOKEN=$(celestia light auth write --p2p.network mocha --node.store ~/.vimana/gmdcelestia/)

# start the chain
# gmd start --rollkit.aggregator true --rollkit.lazy_aggregator --rollkit.da_layer celestia --rollkit.da_config='{"base_url":"http://localhost:26658","timeout":60000000000,"fee":600000,"gas_limit":6000000,"auth_token":"'$AUTH_TOKEN'"}' --rollkit.namespace_id $NAMESPACE --rollkit.da_start_height $DA_BLOCK_HEIGHT
gmd start --rollkit.aggregator true --rollkit.lazy_aggregator --rollkit.da_layer celestia --rollkit.da_config='{"base_url":"http://localhost:26658","timeout":60000000000,"fee":600000,"gas_limit":6000000,"auth_token":"'$AUTH_TOKEN'"}' --rollkit.namespace_id $NAMESPACE --rollkit.da_start_height $DA_BLOCK_HEIGHT
