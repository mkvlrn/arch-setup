import { errResult, okResult, type ResultAsync } from "@mkvlrn/result";
import type { $ } from "bun";

export async function runShellCommand(
  cmd: () => $.ShellPromise,
  name: string,
  quiet = false,
): ResultAsync<true, Error> {
  try {
    if (quiet) {
      await cmd().quiet();
    } else {
      await cmd();
    }

    return okResult(true);
  } catch (cause) {
    return errResult(new Error(`could not run ${name}`, { cause }));
  }
}
