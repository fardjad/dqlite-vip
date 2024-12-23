#!/usr/bin/env bash

set -exuo pipefail

source "$(dirname "$0")/bootstrap.bash" || exit 1

sudo apt-get update
sudo apt-get install -y curl locales iproute2 iputils-ping iputils-arping

test -d ~/.linuxbrew && eval "$(~/.linuxbrew/bin/brew shellenv)"
test -d /home/linuxbrew/.linuxbrew && eval "$(/home/linuxbrew/.linuxbrew/bin/brew shellenv)"

if check_command brew; then
  exit 0
fi

/bin/bash -c "$(curl -fsSL https://raw.githubusercontent.com/Homebrew/install/HEAD/install.sh)"
eval "$(/home/linuxbrew/.linuxbrew/bin/brew shellenv)"

LOCALE="en_US.UTF-8"
sudo locale-gen "$LOCALE"

brew install gcc

git clone https://github.com/fardjad/dotfiles.git ~/.dotfiles
directories_to_keep=("zsh" "starship" "bat" "eza" "mcfly" "tldr" "helper-scripts" "script")
for directory in ~/.dotfiles/*/; do
  dir_name=$(basename "$directory")
  if [[ ! " ${directories_to_keep[@]} " =~ " ${dir_name} " ]]; then
    rm -rf "$directory"
  fi
done
~/.dotfiles/script/setup
touch ~/.zsh_history

brew install vim neovim