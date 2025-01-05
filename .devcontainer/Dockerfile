# syntax=docker/dockerfile:1

FROM rust:bookworm AS rust-tools

RUN cargo install just
RUN cargo install viddy

FROM golang:1.23-bookworm AS builder

COPY --from=rust-tools /usr/local/cargo/bin/viddy /usr/local/bin/viddy
COPY --from=rust-tools /usr/local/cargo/bin/just /usr/local/bin/just
COPY --from=vektra/mockery:v2 /usr/local/bin/mockery /usr/local/bin/mockery

RUN apt-get update && apt-get install -y \
  locales \
  software-properties-common \
  python3-launchpadlib \
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
  tmux \
  jq \
  vim \
  neovim \
  zsh

RUN curl -sSLo /usr/local/bin/xpanes https://raw.githubusercontent.com/greymd/tmux-xpanes/v4.2.0/bin/xpanes \
  && chmod 0755 /usr/local/bin/xpanes

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
  && useradd -s /usr/bin/zsh --uid $USER_UID --gid $USER_GID -m $USERNAME \
  && echo $USERNAME ALL=\(root\) NOPASSWD:ALL > /etc/sudoers.d/$USERNAME \
  && chmod 0440 /etc/sudoers.d/$USERNAME
USER $USERNAME
WORKDIR /home/$USERNAME
ENV HOME=/home/$USERNAME

# Locale
ENV LC_ALL=en_US.UTF-8
ENV LANG=en_US.UTF-8
ENV LANGUAGE=en_US.UTF-8
ENV LC_CTYPE=en_US.UTF-8
RUN sudo sed -i -e 's/# en_US.UTF-8 UTF-8/en_US.UTF-8 UTF-8/' /etc/locale.gen && \
  sudo dpkg-reconfigure --frontend=noninteractive locales

RUN mkdir -p $HOME/.zsh/completion
RUN just --completions=zsh > ${HOME}/.zsh/completion/_just
COPY --chown=${USER_UID}:${USER_GID} <<EOF /home/${USERNAME}/.zshrc

mkdir -p \$HOME/.zsh/completion
fpath=(~/.zsh/completion \$fpath)
autoload -U compinit
compinit

EOF