#!/bin/sh

set -u

REPO="$HOME/repos/arch-setup"
FAILURES=0

fail() {
  printf 'FAILED: %s\n' "$1" >&2
  FAILURES=$((FAILURES + 1))
}

check_stow_tree() {
  stow_dir="$1"
  target="$2"
  list="$(mktemp)"

  find "$stow_dir" -type f >"$list"

  while IFS= read -r source; do
    relative="${source#"$stow_dir"/}"
    relative="${relative#*/}"
    destination="$target/$relative"

    if [ ! -L "$destination" ]; then
      fail "$destination is not a symlink"
      continue
    fi

    source_real="$(readlink -f "$source")"
    destination_real="$(readlink -f "$destination")"

    if [ "$source_real" != "$destination_real" ]; then
      fail "$destination points to $destination_real instead of $source_real"
    fi
  done <"$list"

  rm -f "$list"
}

printf 'Checking repo...\n'

if [ ! -d "$REPO/.git" ]; then
  fail "arch-setup repo was not cloned"
else
  remote="$(git -C "$REPO" remote get-url origin 2>/dev/null || true)"

  if [ "$remote" != "git@github.com:mkvlrn/arch-setup" ]; then
    fail "unexpected git remote: $remote"
  fi

  dirty="$(git -C "$REPO" status --porcelain 2>/dev/null || true)"

  if [ -n "$dirty" ]; then
    fail "arch-setup repo is dirty"
  fi
fi

printf 'Checking system stow...\n'

if [ -d "$REPO/stow/system" ]; then
  check_stow_tree "$REPO/stow/system" ""
else
  fail "$REPO/stow/system does not exist"
fi

printf 'Checking user stow...\n'

if [ -d "$REPO/stow/user" ]; then
  check_stow_tree "$REPO/stow/user" "$HOME"
else
  fail "$REPO/stow/user does not exist"
fi

printf 'Checking directories...\n'

for dir in repos work documents downloads media; do
  if [ ! -d "$HOME/$dir" ]; then
    fail "$HOME/$dir does not exist"
  fi
done

for dir in Desktop Documents Downloads Music Pictures Public Templates Videos; do
  if [ -e "$HOME/$dir" ]; then
    fail "$HOME/$dir still exists"
  fi
done

printf 'Checking packages...\n'

for pkg in \
  base-devel \
  git \
  reflector \
  stow \
  age \
  bruno-bin \
  deluge-gtk \
  deluge \
  docker-buildx \
  docker-compose \
  docker \
  firefox \
  fish \
  keychain \
  kitty \
  less \
  okular \
  openssh \
  power-profiles-daemon \
  pure-ftpd \
  qalculate-qt \
  ttf-hack-nerd \
  ttf-iosevkaterm-nerd \
  ttf-zed-mono-nerd \
  unzip \
  xdg-user-dirs \
  zed; do
  if ! pacman -Q "$pkg" >/dev/null 2>&1; then
    fail "$pkg is not installed"
  fi
done

if ! command -v yay >/dev/null 2>&1; then
  fail "yay is not installed"
fi

if pacman -Qq | grep -q -- '-debug$'; then
  fail "*-debug packages are still installed"
fi

printf 'Checking mirrors...\n'

if ! grep -q '^Server' /etc/pacman.d/mirrorlist; then
  fail "mirrorlist contains no mirrors"
fi

printf 'Checking mise...\n'

MISE="$HOME/.local/bin/mise"

if [ ! -x "$MISE" ]; then
  fail "mise is not installed"
else
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
    if ! "$MISE" exec -- "$cmd" --version >/dev/null 2>&1; then
      fail "mise tool $cmd is not available"
    fi
  done
fi

printf 'Checking user configuration...\n'

fish_path="$(command -v fish 2>/dev/null || true)"
user_shell="$(getent passwd "$USER" | cut -d: -f7)"

if [ -z "$fish_path" ]; then
  fail "fish is not available"
elif [ "$user_shell" != "$fish_path" ]; then
  fail "user shell is $user_shell instead of $fish_path"
fi

if ! id -nG "$USER" | grep -qw docker; then
  fail "$USER is not in the docker group"
fi

if ! systemctl is-enabled docker.socket >/dev/null 2>&1; then
  fail "docker.socket is not enabled"
fi

if ! systemctl is-active docker.socket >/dev/null 2>&1; then
  fail "docker.socket is not active"
fi

if ! systemctl is-active pure-ftpd.socket >/dev/null 2>&1; then
  fail "pure-ftpd.socket is not active"
fi

if [ "$FAILURES" -gt 0 ]; then
  printf '\n%d machine state check(s) failed.\n' "$FAILURES" >&2
  exit 1
fi

printf '\nMachine state checks passed.\n'
