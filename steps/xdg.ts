import os from "node:os";
import type { Step } from "../utils/types";

const mkdir = ["repos", "work", "documents", "downloads", "media"];
const rmrf = [
  "Desktop",
  "Documents",
  "Downloads",
  "Music",
  "Pictures",
  "Public",
  "Templates",
  "Videos",
];

export function xdg(homeDir = os.homedir(), update = Bun.$`xdg-user-dirs-update`): Step {
  return {
    label: "Rework XDG user dirs",
    commands: [
      ["create new xdg set", Bun.$`mkdir -p ${mkdir}`.cwd(homeDir)],
      ["remove old xdg set", Bun.$`rm -rf ${rmrf}`.cwd(homeDir)],
      ["update xdg dirs", update],
    ],
  };
}
