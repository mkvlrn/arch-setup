import { errResult, okResult, type ResultAsync } from "@mkvlrn/result";
import type { Command } from "./types";

export async function runShellCommands(cmds: Command[]): ResultAsync<true, Error> {
  for await (const [name, cmd] of cmds) {
    try {
      await cmd.quiet();
    } catch (cause) {
      return errResult(new Error(`could not run: ${name}`, { cause }));
    }
  }

  return okResult(true);
}
