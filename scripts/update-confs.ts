import fs from "node:fs/promises";
import os from "node:os";
import path from "node:path";
import { errResult, okResult, type ResultAsync } from "@mkvlrn/result";
import type { $ } from "bun";

export interface UpdateConfsArgs {
  etcDir: string;
  files: Record<string, [RegExp, string][]>;
  sudoInstal: (origin: string, destination: string) => $.ShellPromise;
}

export async function updateConfs({
  etcDir,
  files,
  sudoInstal,
}: UpdateConfsArgs): ResultAsync<true, Error> {
  for await (const [file, update] of Object.entries(files)) {
    const conf = await doUpdate(etcDir, file, update, sudoInstal);
    if (conf.isError) {
      return errResult(conf.error);
    }
  }

  return okResult(true);
}

async function doUpdate(
  etcDir: string,
  fileName: string,
  updates: [RegExp, string][],
  sudoInstal: (origin: string, destination: string) => Promise<unknown>,
): ResultAsync<true, Error> {
  try {
    await using tempDir = await fs.mkdtempDisposable(path.join(os.tmpdir(), "arch-setup-"));

    const current = await Bun.file(path.join(etcDir, fileName)).text();
    let updated = current;
    for await (const [regex, update] of updates) {
      updated = updated.replace(regex, update);
    }

    const tempFilePath = path.join(tempDir.path, fileName);
    await Bun.write(tempFilePath, updated);
    await sudoInstal(tempFilePath, path.join(etcDir, fileName));

    return okResult(true);
  } catch (cause) {
    return errResult(new Error(`could not update /etc/${fileName}`, { cause }));
  }
}
