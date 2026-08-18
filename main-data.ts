import os from "node:os";
import { $ } from "bun";
import type { UpdateConfArgs } from "./scripts/update-confs";
import type { XdgArgs } from "./scripts/xdg";

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
