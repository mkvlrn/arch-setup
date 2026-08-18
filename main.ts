import process from "node:process";
import { $ } from "bun";
import { updateConfArgs, xdgArgs, yayArgs } from "./main-data";
import { updateConfs } from "./scripts/update-confs";
import { xdg } from "./scripts/xdg";
import { yay } from "./scripts/yay";

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
    console.info(`[${index + 1}/${scripts.length}] - ${label}`);
    const result = await cmd();
    if (result.isError) {
      console.error(result.error);
      process.exit(1);
    }
  }
} finally {
  keepalive.kill();
}
