import process from "node:process";
import { baseSystem } from "./steps/base-system";
import { cloneRepo } from "./steps/clone-repo";
import { miscUser } from "./steps/misc-user";
import { mise } from "./steps/mise";
import { packages } from "./steps/packages";
import { stow } from "./steps/stow";
import { xdg } from "./steps/xdg";
import { yay } from "./steps/yay";
import { formatError } from "./utils/format-error";
import { runShellCommands } from "./utils/run-shell-commands";
import type { Step } from "./utils/types";

async function main() {
  console.info("Starting.");

  await Bun.$`sudo -v`;
  const keepalive = Bun.spawn(["sh", "-c", "while sleep 30; do sudo -n -v || exit; done"], {
    stdout: "ignore",
    stderr: "ignore",
  });

  const steps: Step[] = [
    baseSystem(),
    cloneRepo(),
    stow("system"),
    yay(),
    packages(),
    xdg(),
    stow("user"),
    mise(),
    miscUser(),
  ];

  try {
    for await (const [index, { label, commands }] of steps.entries()) {
      console.log(`[${index + 1}/${steps.length}] ${label}`);
      const result = await runShellCommands(commands);
      if (result.isError) {
        console.error(formatError(result.error));
        process.exit(1);
      }
    }
  } finally {
    console.info("Done.");
    keepalive.kill();
  }
}

main();
