import os from "node:os";
import path from "node:path";
import process from "node:process";
import type { Step } from "../utils/types";

export function vscodiumExtensions(): Step {
  const repoDir = path.join(os.homedir(), "repos", "arch-setup");

  return {
    label: "Install vscodium extensions",
    commands: [
      [
        "install vscodium extensiosn",
        Bun.$`xargs -n1 codium --install-extension < vscode-extensions.txt`.cwd(repoDir).env({
          ...process.env,
          VSCODE_GALLERY_SERVICE_URL: "https://marketplace.visualstudio.com/_apis/public/gallery",
          VSCODE_GALLERY_ITEM_URL: "https://marketplace.visualstudio.com/items",
        }),
      ],
    ],
  };
}
