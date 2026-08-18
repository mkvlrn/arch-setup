import { errResult, okResult, type ResultAsync } from "@mkvlrn/result";

export async function runShellCommand(
  cmd: () => Promise<unknown>,
  name: string,
): ResultAsync<true, Error> {
  try {
    await cmd();

    return okResult(true);
  } catch (cause) {
    return errResult(new Error(`could not run ${name}`, { cause }));
  }
}
