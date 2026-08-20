import os from "node:os";
import path from "node:path";
import type { Command, Step } from "../utils/types";

type StowType = "system" | "user";

export function stow(type: StowType): Step {
  const repoDir = path.join(os.homedir(), "repos", "arch-setup");
  const stowDir = path.join(repoDir, "stow", type);
  const prefix: string[] = [];
  const stowArgs: string[] = [];
  const commands: Command[] = [];
  let target = os.homedir();

  if (type === "system") {
    prefix.push("sudo");
    target = "/";

    commands.push(["remove existing confs", Bun.$`sudo rm -f /etc/pacman.conf /etc/makepkg.conf`]);
  } else {
    stowArgs.push("--adopt");
  }

  commands.push([
    `stow ${type} files`,
    Bun.$`${prefix} stow -R --no-folding ${stowArgs} -d ${stowDir} -t ${target} $(ls ${stowDir})`,
  ]);

  if (type === "user") {
    commands.push(
      ["restore adopted files", Bun.$`git restore .`.cwd(repoDir)],
      ["clean git state", Bun.$`git clean -fd`.cwd(repoDir)],
    );
  }

  return {
    label: `Stow ${type} files`,
    commands,
  };
}
