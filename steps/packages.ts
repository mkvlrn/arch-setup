import type { Step } from "../utils/types";

export function packages(): Step {
  const pkgs = [
    "brave-bin",
    "bruno-bin",
    "deluge-gtk",
    "deluge",
    "docker-buildx",
    "docker-compose",
    "docker",
    "fish",
    "keychain",
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
    "xdg-user-dirs",
    "zed",
  ];

  return {
    label: "Install packages",
    commands: [
      ["install packages", Bun.$`yay -S --noconfirm --needed ${pkgs}`],
      ["install mise", Bun.$`curl https://mise.run | sh`],
    ],
  };
}
