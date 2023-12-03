#!/bin/bash
set -e

docker run -i -t --platform linux/amd64 -p 26657:26657 -p 26658:26658 -p 26659:26659 ghcr.io/rollkit/local-celestia-devnet:v0.11.0