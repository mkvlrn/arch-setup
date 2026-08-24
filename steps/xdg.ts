import os from "node:os";
import type { Step } from "../utils/types";

export function xdg(): Step {
  const mkdir = ["repos", "work", "documents", "downloads", "media"];
  const rmrf = [
    "Desktop",
    "Documents",
    "Downloads",
    "Music",
    "Pictures",
    "Projects",
    "Public",
    "Templates",
    "Videos",
  ];

  return {
    label: "Rework XDG user dirs",
    commands: [
      ["update xdg dirs", Bun.$`xdg-user-dirs-update`],
      ["create new xdg set", Bun.$`mkdir -p ${mkdir}`.cwd(os.homedir())],
      ["remove old xdg set", Bun.$`rm -rf ${rmrf}`.cwd(os.homedir())],
    ],
  };
}
