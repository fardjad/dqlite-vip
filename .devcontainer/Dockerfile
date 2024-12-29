# syntax=docker/dockerfile:1

FROM ubuntu:noble

# Install dependencies

RUN apt-get update && apt-get install -y \
  software-properties-common \
  apt-utils \
  build-essential \
  automake \
  libtool \
  gettext \
  autopoint \
  tclsh \
  tcl \
  libsqlite3-dev \
  pkg-config \
  git \
  curl \
  locales \
  iproute2 \
  iputils-ping \
  iputils-arping \
  sudo \
  && rm -rf /var/lib/apt/lists/*

RUN sudo locale-gen "en_US.UTF-8" && sudo dpkg-reconfigure locales

# Configure non-root user

# https://code.visualstudio.com/remote/advancedcontainers/add-nonroot-user#_creating-a-nonroot-user
ARG USER_UID=1000
ARG USER_GID=$USER_UID
ARG USERNAME=dev

COPY <<EOF /delete-default-user.sh
#!/usr/bin/env bash
set -euo pipefail

username=$(awk -F: '\$3 == $USER_UID {print \$1}' /etc/passwd)
if [ -n "\$username" ]; then userdel -f -r "\$username"; fi

groupname=\$(awk -F: '\$3 == $USER_GID {print \$1}' /etc/group)
if [ -n "\$groupname" ]; then groupdel "\$groupname"; fi
EOF
RUN bash /delete-default-user.sh && rm /delete-default-user.sh

RUN groupadd --gid $USER_GID $USERNAME \
  && useradd --uid $USER_UID --gid $USER_GID -m $USERNAME \
  && echo $USERNAME ALL=\(root\) NOPASSWD:ALL > /etc/sudoers.d/$USERNAME \
  && chmod 0440 /etc/sudoers.d/$USERNAME
USER $USERNAME
WORKDIR /home/$USERNAME
ENV HOME=/home/$USERNAME

# Install Homebrew

RUN /bin/bash -c "$(curl -fsSL https://raw.githubusercontent.com/Homebrew/install/HEAD/install.sh)" \
  && eval "$(/home/linuxbrew/.linuxbrew/bin/brew shellenv)" \
  brew install gcc

# Dotfiles setup

RUN git clone https://github.com/fardjad/dotfiles.git $HOME/.dotfiles
COPY <<EOF $HOME/fetch-dotfiles-formulas.sh
#!/usr/bin/env bash

set -euo pipefail

eval "$(/home/linuxbrew/.linuxbrew/bin/brew shellenv)"

cd "\$HOME/.dotfiles" || exit 1

brewfiles=\$(find . -name Brewfile -type f -mindepth 2 -maxdepth 2)
for brewfile in \$brewfiles; do
    pushd \$(dirname \$brewfile)
    for formula in \$(brew bundle list); do
        brew fetch \$formula
    done
    popd
done
EOF
RUN bash $HOME/fetch-dotfiles-formulas.sh && sudo rm $HOME/fetch-dotfiles-formulas.sh

COPY <<EOF $HOME/setup-dotfiles.sh
#!/usr/bin/env bash

set -euo pipefail

eval "$(/home/linuxbrew/.linuxbrew/bin/brew shellenv)"

directories_to_keep=(
  "script" # needed for setup
  "bat"
  "eza"
  "golang"
  "mcfly"
  "starship"
  "tldr"
  "zsh"
)
for directory in \$HOME/.dotfiles/*; do
  dir_name=$(basename "\$directory")
  if [[ ! " \${directories_to_keep[@]} " =~ " \${dir_name} " ]]; then
    rm -rf "\$directory"
  fi
done

\$HOME/.dotfiles/script/setup

touch ~/.zsh_history # McFly expects this file to exist on startup
EOF
RUN bash $HOME/setup-dotfiles.sh && sudo rm $HOME/setup-dotfiles.sh

# Customizations

COPY <<EOF $HOME/setup-just.sh
#!/usr/bin/env bash

set -euo pipefail

eval "$(/home/linuxbrew/.linuxbrew/bin/brew shellenv)"

brew install just
EOF
RUN bash $HOME/setup-just.sh && sudo rm $HOME/setup-just.sh

COPY <<EOF $HOME/setup-editor.sh
#!/usr/bin/env bash

set -euo pipefail

eval "$(/home/linuxbrew/.linuxbrew/bin/brew shellenv)"

brew install vim neovim
echo "export EDITOR=nvim" >> \$HOME/.localrc
echo "set clipboard+=unnamedplus" >> \$HOME/.ideavimrc

EOF
RUN bash $HOME/setup-editor.sh && sudo rm $HOME/setup-editor.sh

COPY <<EOF $HOME/setup-starship.sh
#!/usr/bin/env bash

set -euo pipefail

eval "$(/home/linuxbrew/.linuxbrew/bin/brew shellenv)"

# IntelliJ terminal doesn't like non-plain-text symbols
starship preset plain-text-symbols > \$HOME/.starship/config.toml
printf "\\n[container]\\ndisabled = true\\n" >> \$HOME/.starship/config.toml
EOF
RUN bash $HOME/setup-starship.sh && sudo rm $HOME/setup-starship.sh

COPY <<EOF $HOME/setup-tmux-xpanes.sh
#!/usr/bin/env bash

set -euo pipefail

eval "$(/home/linuxbrew/.linuxbrew/bin/brew shellenv)"

brew install tmux-xpanes
EOF
RUN bash $HOME/setup-tmux-xpanes.sh && sudo rm $HOME/setup-tmux-xpanes.sh

COPY <<EOF $HOME/setup-viddy.sh
#!/usr/bin/env bash

set -euo pipefail

eval "$(/home/linuxbrew/.linuxbrew/bin/brew shellenv)"

brew install viddy
EOF
RUN bash $HOME/setup-viddy.sh && sudo rm $HOME/setup-viddy.sh

COPY <<EOF $HOME/setup-jq.sh
#!/usr/bin/env bash

set -euo pipefail

eval "$(/home/linuxbrew/.linuxbrew/bin/brew shellenv)"

brew install jq
EOF
RUN bash $HOME/setup-jq.sh && sudo rm $HOME/setup-jq.sh
