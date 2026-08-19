import fs from "node:fs/promises";
import os from "node:os";
import path from "node:path";
import { errResult, okResult, type ResultAsync } from "@mkvlrn/result";
import type { $ } from "bun";

export interface SystemFilesArgs {
  copies?: Record<string, string>;
  updates?: Record<string, [RegExp, string][]>;
  install: (src: string, dest: string) => $.ShellPromise;
}

export async function systemFiles({
  copies,
  updates,
  install,
}: SystemFilesArgs): ResultAsync<true, Error> {
  if (copies) {
    for await (const [src, dest] of Object.entries(copies)) {
      const copied = await doCopy(src, dest, install);
      if (copied.isError) {
        return errResult(copied.error);
      }
    }
  }

  if (updates) {
    for await (const [file, update] of Object.entries(updates)) {
      const updated = await doUpdate(file, update, install);
      if (updated.isError) {
        return errResult(updated.error);
      }
    }
  }

  return okResult(true);
}

async function doCopy(
  src: string,
  dest: string,
  install: (src: string, dest: string) => $.ShellPromise,
): ResultAsync<true, Error> {
  try {
    await using tempDir = await fs.mkdtempDisposable(
      path.join(os.tmpdir(), "arch-setup-system-files-"),
    );
    const tempFile = path.join(tempDir.path, path.basename(src));
    await Bun.write(tempFile, Bun.file(src));
    await install(tempFile, dest);

    return okResult(true);
  } catch (err) {
    return errResult(new Error("could not copy system files", { cause: err }));
  }
}

async function doUpdate(
  fileName: string,
  updates: [RegExp, string][],
  install: (src: string, dest: string) => $.ShellPromise,
): ResultAsync<true, Error> {
  try {
    await using tempDir = await fs.mkdtempDisposable(
      path.join(os.tmpdir(), "arch-setup-system-files-"),
    );

    const current = await Bun.file(fileName).text();
    let updated = current;
    for await (const [regex, update] of updates) {
      updated = updated.replace(regex, update);
    }

    const tempFilePath = path.join(tempDir.path, fileName);
    await Bun.write(tempFilePath, updated);
    await install(tempFilePath, fileName);

    return okResult(true);
  } catch (cause) {
    return errResult(new Error(`could not update ${fileName}`, { cause }));
  }
}
