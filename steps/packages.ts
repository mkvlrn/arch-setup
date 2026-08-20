import type { Step } from "../utils/types";

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

export function packages(): Step {
  return {
    label: "Install packages with yay and mise via curl",
    commands: [
      ["install packages", Bun.$`yay -S --noconfirm --needed ${pkgs}`],
      ["install mise", Bun.$`curl https://mise.run | sh`],
    ],
  };
}
