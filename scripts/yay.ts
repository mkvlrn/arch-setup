import { errResult, okResult, type ResultAsync } from "@mkvlrn/result";
import { runShellCommand } from "../utils/run-shell-command";

export interface YayArgs {
  yayCmds: [() => Promise<unknown>, string][];
}

export async function yay({ yayCmds }: YayArgs): ResultAsync<true, Error> {
  for await (const [cmd, name] of yayCmds) {
    const result = await runShellCommand(cmd, name);
    if (result.isError) {
      return errResult(result.error);
    }
  }

  return okResult(true);
}
