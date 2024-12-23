#!/usr/bin/env bash

set -euo pipefail
sudo -i -u "$_CONTAINER_USER" "$PWD/scripts/install-dotfiles.sh"