#!/bin/bash
set -e

INTERNAL_DIR="/usr/local/bin"

# check if the binary is already installed
if [ -f "$INTERNAL_DIR/celestia" ]; then
    # Capture the version output
    VERSION_OUTPUT=$("$INTERNAL_DIR/celestia" version)
    
    # Check if the version matches "v0.11.0-rc15-dev"
    if [[ $VERSION_OUTPUT == *"v0.12"* ]]; then
        echo "🚀 Celestia is already installed with the correct version." $VERSION_OUTPUT
        # celestia light init --p2p.network mocha
        celestia bridge start --core.ip rpc-mocha.pops.one --p2p.network mocha --node.store ~/.vimana/gmdcelestia/
        exit 0
    else
        echo "🚀 Celestia is installed but with a different version."
    fi
fi

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

# Download celestia binary

TGZ_URL="https://github.com/Vistara-Labs/vimana/releases/download/celestia-v0.12.0/${OS}_${ARCH}.zip"
sudo mkdir -p "/tmp/vimcel"
echo "💈 Downloading Celestia..."
sudo curl -o /tmp/vimcel/${OS}_${ARCH}.zip -L "$TGZ_URL" --progress-bar

sudo unzip -q /tmp/vimcel/${OS}_${ARCH}.zip -d /tmp/vimcel/
sudo mv "/tmp/vimcel/${OS}_${ARCH}"/* "$INTERNAL_DIR"
sudo chmod +x "$INTERNAL_DIR"
sudo rm -rf "/tmp/vimcel"

if [ ! -f "$HOME/.vimana/gmdcelestia/config.yml" ]; then
    # This should be handled in the InitializeConfig code
    celestia bridge init --core.ip rpc-mocha.pops.one --p2p.network mocha --node.store ~/.vimana/gmdcelestia/
fi

# Handle this in GetStartCmd
celestia bridge start --core.ip rpc-mocha.pops.one --p2p.network mocha --node.store ~/.vimana/gmdcelestia/
