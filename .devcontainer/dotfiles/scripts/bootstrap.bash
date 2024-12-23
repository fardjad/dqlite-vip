pushd "$(dirname "$0")/.." > /dev/null
export DOTFILES=$(pwd -P)
popd > /dev/null

export NONINTERACTIVE=1
export HOMEBREW_BUNDLE_NO_LOCK=1

info() {
  printf "\r  [ \033[00;34m..\033[0m ] $1\n"
}

user() {
  printf "\r  [ \033[0;33m??\033[0m ] $1\n"
}

success() {
  printf "\r\033[2K  [ \033[00;32mOK\033[0m ] $1\n"
}

fail() {
  printf "\r\033[2K  [\033[0;31mFAIL\033[0m] $1\n"
  echo ''
  exit -1
}

check_command() {
  command -v "$1" > /dev/null 2>&1
}

brew_bundle_install() {
  if ! check_command brew; then
    fail 'brew must be installed'
  fi

  pushd "$(dirname "$0")" > /dev/null

  brew bundle install -q "$@" \
    | sed '/^Homebrew Bundle complete.*/d' \
    | sed 's/^/  [brew] /'

  popd > /dev/null
}

is_mac() {
  [[ "$OSTYPE" == "darwin"* ]]
}

is_codespaces() {
  [[ -n "$CODESPACES" ]]
}

link_file() {
  pushd "$(dirname "$0")" > /dev/null

  src="$(readlink -f "$1")"
  if [ $? -ne 0 ]; then
    fail "could not find $1"
  fi
  dst="$2"

  if [ -e "$2" ]; then
    if [ "$src" = "$(readlink "$dst")" ]; then
      success "skipped $src"
      return 0
    else
      mv "$dst" "$dst.backup"
      success "moved $dst to $dst.backup"
    fi
  fi
  ln -sf "$src" "$dst"
  success "linked $src to $dst"

  popd > /dev/null
}