#!/bin/bash
set -e
OS=$(uname -s | tr '[:upper:]' '[:lower:]')
ARCH=$(uname -m)
if [[ "$ARCH" == "x86_64" ]]; then
    ARCH="amd64"
elif [[ "$ARCH" == "arm64" ]] || [[ "$ARCH" == "aarch64" ]]; then
    ARCH="arm64"
fi
VIMANA_RELEASE_TAG="v0.0.12"
GZ_URL="https://github.com/Vistara-Labs/vimana/releases/download/vimana-${VIMANA_RELEASE_TAG}/vimana-${OS}-${ARCH}.tar.gz"

INTERNAL_DIR="/usr/local/bin/"
VIMANA_BIN_PATH="/usr/local/bin/vimana"
if [ -f "$VIMANA_BIN_PATH" ] || [ -f "$INTERNAL_DIR" ]; then
    sudo rm -f "$VIMANA_BIN_PATH"
    sudo rm -rf "$INTERNAL_DIR"
fi
sudo mkdir -p "$INTERNAL_DIR"
sudo mkdir -p "/tmp/vimana_bins"
echo "💿 Downloading vimana..."
curl -O -L $GZ_URL --progress-bar

sudo tar -xzf vimana-${OS}-${ARCH}.tar.gz -C "/tmp/vimana_bins" 2>/dev/null

echo "🔨 Installing vimana..."
sudo cp "/tmp/vimana_bins/vimana-${OS}-${ARCH}/vimana" "$INTERNAL_DIR/vimana"
sudo chmod +x "$INTERNAL_DIR/vimana"
sudo rm -rf "/tmp/vimana_bins"
curl -O https://vistara-labs.github.io/vimana/config.toml 2>/dev/null
mkdir -p ~/.vimana && cp config.toml ~/.vimana/config.toml
curl -O https://vistara-labs.github.io/vimana/scripts/init.sh 2>/dev/null
sudo mkdir -p /tmp/vimana/celestia && sudo mv init.sh /tmp/vimana/celestia/init.sh

# Step 1: Get availup script from repo
curl -O https://vistara-labs.github.io/vimana/scripts/availup.sh 2>/dev/null
sudo mkdir -p /tmp/vimana/avail && sudo mv availup.sh /tmp/vimana/avail/init.sh

sudo chmod +x /tmp/vimana/celestia/init.sh
sudo chmod +x /tmp/vimana/avail/init.sh
mkdir -p ~/.vimana/celestia/light-node
chmod +x ~/.vimana/celestia/light-node
echo "✅ vimana installed!"
