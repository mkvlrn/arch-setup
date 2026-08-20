import os from "node:os";
import path from "node:path";
import type { Step } from "../utils/types";

export function mise(): Step {
  const misePath = path.join(os.homedir(), ".local", "bin", "mise");

  return {
    label: "Install tools with mise",
    commands: [["mise install", Bun.$`${misePath} install`]],
  };
}
