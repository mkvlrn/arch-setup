import process from "node:process";
import { $ } from "bun";
import { baseArgs, pkgArgs, repoArgs, systemFilesArgs, xdgArgs, yayArgs } from "./main-data";
import { systemFiles } from "./scripts/system-files";
import { xdg } from "./scripts/xdg";
import { formatError } from "./utils/format-error";
import { runShellCommands } from "./utils/run-shell-commands";

async function main() {
  await $`sudo -v`;

  const keepalive = Bun.spawn(["sh", "-c", "while sleep 30; do sudo -n true || exit; done"], {
    stdout: "ignore",
    stderr: "ignore",
  });

  const scripts = [
    [() => systemFiles(systemFilesArgs), "Update and copy new system files"],
    [() => runShellCommands(baseArgs), "Install minimal packages to build yay"],
    [() => runShellCommands(yayArgs), "Install yay and update mirrors"],
    [() => runShellCommands(pkgArgs), "Install packages with yay and do cleanup"],
    [() => runShellCommands(repoArgs), "Clone arch-setup repo to ~/repos/arch-setup"],
    [() => xdg(xdgArgs), "Configure XDG user directories"],
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
}

main();
