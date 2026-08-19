import os from "node:os";
import path from "node:path";
import { $ } from "bun";
import type { UpdateConfsArgs } from "./scripts/update-confs";
import type { XdgArgs } from "./scripts/xdg";
import type { Command } from "./utils/run-shell-commands";

export const updateConfsArgs = {
  etcDir: "/etc",
  files: {
    "makepkg.conf": [
      [
        /^OPTIONS=.*$/m,
        "OPTIONS=(strip docs !libtool !staticlibs emptydirs zipman purge !debug lto)",
      ],
      [/^MAKEFLAGS=.*$/m, `MAKEFLAGS="--jobs=$(nproc)"`],
    ],
    "pacman.conf": [[/^#Color$/m, "Color"]],
  },
  sudoInstal: (origin: string, destination: string) =>
    $`sudo install -m 644 ${origin} ${destination}`,
} satisfies UpdateConfsArgs;

export const baseArgs = [
  { exec: () => $`sudo -v`, name: "sudo timestamp", quiet: true },
  {
    exec: () => $`sudo pacmain -Syu --noconfirm --needed git base-devel reflector`,
    name: "install basic packages for yay",
    quiet: true,
  },
] satisfies Command[];

export const yayArgs = [
  {
    exec: () => $`git clone https://aur.archlinux.org/yay-bin`.cwd(os.homedir()),
    name: "clone yay-bin",
    quiet: true,
  },
  {
    exec: () => $`makepkg -si --noconfirm`.cwd(path.join(os.homedir(), "yay-bin")),
    name: "make yay",
    quiet: true,
  },
  {
    exec: () => $`rm -rf ${path.join(os.homedir(), "yay-bin")}`,
    name: "remove yay-bin src dir",
    quiet: true,
  },
  { exec: () => $`yay -Y --gendb`, name: "track git packages", quiet: true },
  { exec: () => $`yay -Y --devel --save`, name: "enable dev packages updates", quiet: true },
  {
    exec: () =>
      $`sudo reflector --latest 20 --protocol https --sort rate --save /etc/pacman.d/mirrorlist`,
    name: "get best mirror list",
    quiet: true,
  },
  { exec: () => $`yay -Syu --noconfirm`, name: "update package data", quiet: true },
  {
    exec: () => $`yay -Qq | grep -- '-debug$' | xargs -r yay -Rnsu`,
    name: "remove *-debug packages",
    quiet: true,
  },
] satisfies Command[];

const pkgs = [
  "brave-bin",
  "bruno-bin",
  "deluge-gtk",
  "deluge",
  "docker-buildx",
  "docker-compose",
  "docker",
  "fish",
  "kitty",
  "less",
  "okular",
  "openssh",
  "power-profiles-daemon",
  "qalculate-qt",
  "stow",
  "ttf-hack-nerd",
  "ttf-iosevkaterm-nerd",
  "ttf-zed-mono-nerd",
  "unzip",
  "vscodium-bin",
  "xdg-user-dirs",
  "zed",
];
export const pkgArgs = [
  {
    exec: () => $`yay -S --noconfirm --needed ${pkgs.join(" ")}`,
    name: "update package data",
    quiet: true,
  },
] satisfies Command[];

export const xdgArgs = {
  homeDir: os.homedir(),
  mkdir: ["repos", "sandbox", "work", "documents", "downloads", "downloads/torrents"],
  rmrf: ["Desktop", "Documents", "Downloads", "Music", "Pictures", "Templates", "Videos"],
  cmd: { exec: () => $`xdg-user-dirs-update`, name: "xdg dirs update", quiet: true },
} satisfies XdgArgs;

export const repoArgs = [
  {
    exec: () =>
      $`git clone https://github.com/mkvlrn/arch-setup ${path.join(os.homedir(), "repos", "arch-setup")}`,
    name: "clone repo",
    quiet: true,
  },
] satisfies Command[];
