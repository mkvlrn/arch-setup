import os from "node:os";
import path from "node:path";
import type { Step } from "../utils/types";

export function yay(): Step {
  const yaySrcDir = path.join(os.tmpdir(), "yay-bin");

  return {
    label: "Install yay and update mirrors",
    commands: [
      ["clone yay-bin", Bun.$`git clone https://aur.archlinux.org/yay-bin ${yaySrcDir}`],
      ["build yay", Bun.$`makepkg -si --noconfirm`.cwd(yaySrcDir)],
      ["remove yay-bin src dir", Bun.$`rm -rf ${yaySrcDir}`],
      ["track git packages", Bun.$`yay -Y --gendb`],
      ["enable dev packages updates", Bun.$`yay -Y --devel --save`],
      [
        "get best mirror list",
        Bun.$`sudo reflector --latest 20 --protocol https --sort rate --save /etc/pacman.d/mirrorlist`,
      ],
      ["update package data", Bun.$`yay -Syu --noconfirm`],
      ["remove *-debug packages", Bun.$`yay -Qq | grep -- '-debug$' | xargs -r yay -Rnsu`],
    ],
  };
}
