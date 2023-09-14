#!/bin/bash
set -e

INTERNAL_DIR="/tmp/vimana/celestia"

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

API_URL="https://api.github.com/repos/dymensionxyz/roller/releases/latest"
if [ -z "$ROLLER_RELEASE_TAG" ]; then
  TGZ_URL=$(curl -s $API_URL \
      | grep "browser_download_url.*_${OS}_${ARCH}.tar.gz" \
      | cut -d : -f 2,3 \
      | tr -d \" \
      | tr -d ' ' )
else
  TGZ_URL="https://github.com/dymensionxyz/roller/releases/download/$ROLLER_RELEASE_TAG/roller_${ROLLER_RELEASE_TAG}_${OS}_${ARCH}.tar.gz"
fi

sudo mkdir -p "$INTERNAL_DIR"
sudo mkdir -p "/tmp/vimcel"
echo "💈 Downloading vimana..."
# Replace this with vistara-labs repo
sudo curl -L "$TGZ_URL" --progress-bar | sudo tar -xz -C "/tmp/vimcel"
sudo mv "/tmp/vimcel/roller_bins/lib"/* "$INTERNAL_DIR"
sudo chmod +x "$INTERNAL_DIR"
sudo rm -rf "/tmp/vimcel"
echo "💈 Celestia light node installed!"
