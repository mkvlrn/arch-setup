import os from "node:os";
import path from "node:path";
import type { Step } from "../utils/types";

export function cloneRepo(): Step {
  const repo = "mkvlrn/arch-setup";
  const repoDir = path.join(os.homedir(), "repos", "arch-setup");

  return {
    label: `Clone ${repo} to ${repoDir}`,
    commands: [
      ["clone arch-setup", Bun.$`git clone https://github.com/${repo} ${repoDir}`],
      ["set ssh upstream", Bun.$`git remote set-url origin git@github.com:${repo}`.cwd(repoDir)],
    ],
  };
}
