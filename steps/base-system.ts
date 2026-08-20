import type { Step } from "../utils/types";

const basePackages = ["git", "base-devel", "reflector", "stow"];

export function baseSystem(): Step {
  return {
    label: "Install minimal packages to build yay and stow files",
    commands: [
      [
        "install basic packages for yay and stow",
        Bun.$`sudo pacman -Syu --noconfirm --needed ${basePackages}`,
      ],
    ],
  };
}
