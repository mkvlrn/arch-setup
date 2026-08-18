import fs from "node:fs/promises";
import os from "node:os";
import path from "node:path";
import { $ } from "bun";
import type { UpdateConfArgs } from "./scripts/update-confs";
import type { XdgArgs } from "./scripts/xdg";
import type { YayArgs } from "./scripts/yay";

export const xdgArgs = {
  homeDir: os.homedir(),
  mkdir: ["repos", "sandbox", "work", "documents", "downloads", "downloads/torrents"],
  rmrf: ["Desktop", "Documents", "Downloads", "Music", "Pictures", "Templates", "Videos"],
  xdgCmd: () => $`xdg-user-dirs-update`,
} satisfies XdgArgs;

export const updateConfArgs = {
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
} satisfies UpdateConfArgs;

export const yayArgs = {
  yayCmds: [
    [() => $`git clone https://aur.archlinux.org/yay-bin`, "clone yay-bin"],
    [() => $`makepkg -si`.cwd(path.join(os.homedir(), "yay-bin")), "make yay"],
    [
      () => fs.rm(path.join(os.homedir(), "yay-bin"), { recursive: true, force: true }),
      "remove yay-bin src dir",
    ],
    [() => $`yay -Y --gendb`, "track git packages"],
    [() => $`yay -Y --devel --save`, "enable dev packages updates"],
    [
      () =>
        $`sudo reflector --latest 20 --protocol https --sort rate --save /etc/pacman.d/mirrorlist`,
      "get best mirror list",
    ],
    [() => $`yay -Syu --noconfirm`, "update package data"],
  ],
} satisfies YayArgs;
