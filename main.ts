import process from "node:process";
import { $ } from "bun";
import { updateConfArgs, xdgArgs, yayArgs } from "./main-data";
import { updateConfs } from "./scripts/update-confs";
import { xdg } from "./scripts/xdg";
import { yay } from "./scripts/yay";
import { formatError } from "./utils/format-error";
import { createSpinner } from "./utils/spinner";

await $`sudo -v`;

const keepalive = Bun.spawn(["sh", "-c", "while sleep 30; do sudo -n true || exit; done"], {
  stdout: "ignore",
  stderr: "ignore",
});

const scripts = [
  ["Configuring XDG user dirs", () => xdg(xdgArgs)],
  ["Configuring pacman", () => updateConfs(updateConfArgs)],
  ["Installing yay and refreshing mirrors", () => yay(yayArgs)],
] as const;

try {
  for await (const [index, script] of scripts.entries()) {
    const [label, cmd] = script;
    const spinner = createSpinner(`[${index + 1}/${scripts.length}] ${label}`);
    const result = await cmd();
    if (result.isError) {
      spinner.failure();
      console.error(formatError(result.error));
      process.exit(1);
    }

    spinner.success();
  }
} finally {
  console.info("Done.");
  keepalive.kill();
}
