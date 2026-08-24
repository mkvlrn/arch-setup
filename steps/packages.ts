import type { Step } from "../utils/types";

export function packages(): Step {
  const pkgs = [
    "age",
    "bruno-bin",
    "deluge-gtk",
    "deluge",
    "docker-buildx",
    "docker-compose",
    "docker",
    "firefox",
    "fish",
    "keychain",
    "kitty",
    "less",
    "okular",
    "openssh",
    "power-profiles-daemon",
    "pure-ftpd",
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
