#!/bin/bash
set -e

INTERNAL_DIR="/tmp/vimana/celestias"

# check if the binary is already installed
if [ -f "$INTERNAL_DIR/celestia" ]; then
    echo "🚀 Celestia is already installed."
    exit 0
fi

echo "🔍  Determining OS and architecture..."

OS=$(uname -s | tr '[:upper:]' '[:lower:]')
ARCH=$(uname -m)

if [[ "$ARCH" == "x86_64" ]]; then
    ARCH="amd64"
elif [[ "$ARCH" == "arm64" ]] || [[ "$ARCH" == "aarch64" ]]; then
    ARCH="arm64"
fi

echo "💻  OS: $OS"
echo "🏗️  ARCH: $ARCH"

TGZ_URL="https://github.com/Vistara-Labs/vimana/releases/download/celestia-v0.10.4/${OS}_${ARCH}.zip"

sudo mkdir -p "$INTERNAL_DIR"
sudo mkdir -p "/tmp/vimcel"
echo "💈 Downloading vimana..."
# Replace this with vistara-labs repo
sudo curl -o /tmp/vimcel/${OS}_${ARCH}.zip -L "$TGZ_URL" --progress-bar
sudo unzip -q /tmp/vimcel/${OS}_${ARCH}.zip -d /tmp/vimcel/
sudo mv "/tmp/vimcel/${OS}_${ARCH}"/* "$INTERNAL_DIR"
sudo chmod +x "$INTERNAL_DIR"
sudo rm -rf "/tmp/vimcel"
${INTERNAL_DIR}/celestia version
echo "💈 Celestia light node version installed!"
