#!/bin/bash
set -e
OS=$(uname -s | tr '[:upper:]' '[:lower:]')
ARCH=$(uname -m)
if [[ "$ARCH" == "x86_64" ]]; then
    ARCH="amd64"
elif [[ "$ARCH" == "arm64" ]] || [[ "$ARCH" == "aarch64" ]]; then
    ARCH="arm64"
fi
VIMANA_RELEASE_TAG="0.0.1"
GZ_URL="https://github.com/Vistara-Labs/vimana/releases/download/vimana-${VIMANA_RELEASE_TAG}/vimana-${OS}-${ARCH}.tar.gz"
# echo "💻  OS: $OS"
# echo "🏗️  ARCH: $ARCH"

INTERNAL_DIR="/usr/local/bin/"
VIMANA_BIN_PATH="/usr/local/bin/vimana"
if [ -f "$VIMANA_BIN_PATH" ] || [ -f "$INTERNAL_DIR" ]; then
    sudo rm -f "$VIMANA_BIN_PATH"
    sudo rm -rf "$INTERNAL_DIR"
fi
sudo mkdir -p "$INTERNAL_DIR"
sudo mkdir -p "/tmp/vimana_bins"
echo "💿Downloading vimana..."
curl -O -L $GZ_URL --progress-bar

sudo tar -xzf vimana-${OS}-${ARCH}.tar.gz -C "/tmp/vimana_bins"

echo "🔨Installing vimana..."
sudo cp "/tmp/vimana_bins/vimana-${OS}-${ARCH}/vimana" "$INTERNAL_DIR/vimana"
sudo chmod +x "$INTERNAL_DIR/vimana"
sudo rm -rf "/tmp/vimana_bins"
echo "✅ vimana installed!"