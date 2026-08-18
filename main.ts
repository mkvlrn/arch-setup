import process from "node:process";
import { $ } from "bun";
import { updateConfArgs, xdgArgs } from "./main-data";
import { updateConfs } from "./scripts/update-confs";
import { xdg } from "./scripts/xdg";

await $`sudo -v`;

const keepalive = Bun.spawn(["sh", "-c", "while sleep 30; do sudo -n true || exit; done"], {
  stdout: "ignore",
  stderr: "ignore",
});

try {
  for await (const script of [() => xdg(xdgArgs), () => updateConfs(updateConfArgs)]) {
    const result = await script();
    if (result.isError) {
      console.error(result.error);
      process.exit(1);
    }
  }
} finally {
  keepalive.kill();
}
