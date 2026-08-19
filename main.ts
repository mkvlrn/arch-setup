import process from "node:process";
import { $ } from "bun";
import { updateConfsArgs, xdgArgs, yayArgs } from "./main-data";
import { updateConfs } from "./scripts/update-confs";
import { xdg } from "./scripts/xdg";
import { yay } from "./scripts/yay";
import { formatError } from "./utils/format-error";

await $`sudo -v`;

const keepalive = Bun.spawn(["sh", "-c", "while sleep 30; do sudo -n true || exit; done"], {
  stdout: "ignore",
  stderr: "ignore",
});

const scripts = [
  [() => xdg(xdgArgs), "Configure XDG user directories"],
  [() => updateConfs(updateConfsArgs), "Configure pacman"],
  [() => yay(yayArgs), "Install yay and update mirrors"],
] as const;

try {
  for await (const [index, script] of scripts.entries()) {
    const [cmd, label] = script;
    console.log(`[${index + 1}/${scripts.length}] ${label}`);
    const result = await cmd();
    if (result.isError) {
      console.error(formatError(result.error));
      process.exit(1);
    }
  }
} finally {
  console.info("Done.");
  keepalive.kill();
}
