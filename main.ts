import process from "node:process";
import { $ } from "bun";
import { updateConfArgs, xdgArgs, yayArgs } from "./main-data";
import { updateConfs } from "./scripts/update-confs";
import { xdg } from "./scripts/xdg";
import { yay } from "./scripts/yay";
import { labels, steps } from "./utils/progress";

await $`sudo -v`;

const keepalive = Bun.spawn(["sh", "-c", "while sleep 30; do sudo -n true || exit; done"], {
  stdout: "ignore",
  stderr: "ignore",
});

try {
  for await (const [index, script] of [
    () => xdg(xdgArgs),
    () => updateConfs(updateConfArgs),
    () => yay(yayArgs),
  ].entries()) {
    const label = `[${index + 1}/${steps}] ${labels[script.name] ?? ""}`;
    console.info(label);
    const result = await script();
    if (result.isError) {
      console.error(result.error);
      process.exit(1);
    }
  }
} finally {
  keepalive.kill();
}
