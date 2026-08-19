import { errResult, okResult, type ResultAsync } from "@mkvlrn/result";
import type { $ } from "bun";

export interface Command {
  exec: () => $.ShellPromise;
  name: string;
  quiet: boolean;
}

export async function runShellCommands(cmds: Command[]): ResultAsync<true, Error> {
  for await (const cmd of cmds) {
    try {
      if (cmd.quiet) {
        await cmd.exec().quiet();
      } else {
        await cmd.exec();
      }
    } catch (cause) {
      return errResult(new Error(`could not run ${cmd.name}`, { cause }));
    }
  }

  return okResult(true);
}
