#!/bin/sh

set -eu

REPO="$HOME/repos/arch-setup"

fail() {
  printf 'FAILED: %s\n' "$1" >&2
  exit 1
}

check_stow_tree() {
  stow_dir="$1"
  target="$2"

  find "$stow_dir" -type f | while read -r source; do
    relative="${source#"$stow_dir"/}"
    relative="${relative#*/}"
    destination="$target/$relative"

    [ -L "$destination" ] ||
      fail "$destination is not a symlink"

    source_real="$(readlink -f "$source")"
    destination_real="$(readlink -f "$destination")"

    [ "$source_real" = "$destination_real" ] ||
      fail "$destination points to $destination_real instead of $source_real"
  done
}

printf 'Checking repo...\n'

[ -d "$REPO/.git" ] ||
  fail "arch-setup repo was not cloned"

remote="$(git -C "$REPO" remote get-url origin)"

[ "$remote" = "git@github.com:mkvlrn/arch-setup" ] ||
  fail "unexpected git remote: $remote"

[ -z "$(git -C "$REPO" status --porcelain)" ] ||
  fail "arch-setup repo is dirty"

printf 'Checking system stow...\n'

check_stow_tree "$REPO/stow/system" ""

printf 'Checking user stow...\n'

check_stow_tree "$REPO/stow/user" "$HOME"

printf 'Checking directories...\n'

for dir in repos work documents downloads media; do
  [ -d "$HOME/$dir" ] ||
    fail "$HOME/$dir does not exist"
done

for dir in Desktop Documents Downloads Music Pictures Public Templates Videos; do
  [ ! -e "$HOME/$dir" ] ||
    fail "$HOME/$dir still exists"
done

printf 'Checking packages...\n'

for pkg in \
  base-devel \
  git \
  reflector \
  stow \
  brave-bin \
  bruno-bin \
  deluge-gtk \
  deluge \
  docker-buildx \
  docker-compose \
  docker \
  fish \
  kitty \
  less \
  okular \
  openssh \
  power-profiles-daemon \
  qalculate-qt \
  ttf-hack-nerd \
  ttf-iosevkaterm-nerd \
  ttf-zed-mono-nerd \
  unzip \
  vscodium-bin \
  xdg-user-dirs \
  zed; do
  pacman -Q "$pkg" >/dev/null 2>&1 ||
    fail "$pkg is not installed"
done

command -v yay >/dev/null 2>&1 ||
  fail "yay is not installed"

if pacman -Qq | grep -q -- '-debug$'; then
  fail "*-debug packages are still installed"
fi

printf 'Checking mirrors...\n'

grep -q '^Server' /etc/pacman.d/mirrorlist ||
  fail "mirrorlist contains no mirrors"

printf 'Checking mise...\n'

MISE="$HOME/.local/bin/mise"

[ -x "$MISE" ] ||
  fail "mise is not installed"

for cmd in \
  usage \
  jq \
  glab \
  gh \
  eza \
  fastfetch \
  onefetch \
  lazydocker \
  biome \
  aws \
  devcontainer \
  node \
  oh-my-posh \
  ncu; do
  "$MISE" exec -- "$cmd" --version >/dev/null 2>&1 ||
    fail "mise tool $cmd is not available"
done

printf 'Checking user configuration...\n'

fish_path="$(command -v fish)"
user_shell="$(getent passwd "$USER" | cut -d: -f7)"

[ "$user_shell" = "$fish_path" ] ||
  fail "user shell is $user_shell instead of $fish_path"

id -nG "$USER" | grep -qw docker ||
  fail "$USER is not in the docker group"

systemctl is-enabled docker.socket >/dev/null ||
  fail "docker.socket is not enabled"

systemctl is-active docker.socket >/dev/null ||
  fail "docker.socket is not active"

printf 'Machine state checks passed.\n'
