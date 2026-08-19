import { errResult, okResult, type ResultAsync } from "@mkvlrn/result";
import { type Command, runShellCommands } from "../utils/run-shell-commands";

export interface YayArgs {
  cmds: Command[];
}

export async function yay({ cmds }: YayArgs): ResultAsync<true, Error> {
  const result = await runShellCommands(cmds.map((cmd) => ({ ...cmd, quiet: true })));
  if (result.isError) {
    return errResult(result.error);
  }

  return okResult(true);
}
