import { errResult, okResult, type ResultAsync } from "@mkvlrn/result";
import type { $ } from "bun";
import { runShellCommand } from "../utils/run-shell-command";

export interface YayArgs {
  yayCmds: [() => $.ShellPromise, string][];
}

export async function yay({ yayCmds }: YayArgs): ResultAsync<true, Error> {
  for await (const [cmd, name] of yayCmds) {
    const result = await runShellCommand(cmd, name, true);
    if (result.isError) {
      return errResult(result.error);
    }
  }

  return okResult(true);
}
