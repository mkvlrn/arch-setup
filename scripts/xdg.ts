import fs from "node:fs/promises";
import path from "node:path";
import { errResult, okResult, type ResultAsync } from "@mkvlrn/result";
import { type Command, runShellCommands } from "../utils/run-shell-commands";

export interface XdgArgs {
  homeDir: string;
  mkdir: string[];
  rmrf: string[];
  cmd: Command;
}

export async function xdg({ homeDir, mkdir, rmrf, cmd }: XdgArgs): ResultAsync<true, Error> {
  const create = await createDirs(homeDir, mkdir);
  if (create.isError) {
    return errResult(create.error);
  }

  const remove = await removeDirs(homeDir, rmrf);
  if (remove.isError) {
    return errResult(remove.error);
  }

  const update = await runShellCommands([cmd]);
  if (update.isError) {
    return errResult(update.error);
  }

  return okResult(true);
}

async function createDirs(homeDir: string, dirs: string[]): ResultAsync<true, Error> {
  try {
    for await (const dir of dirs) {
      await fs.mkdir(path.join(homeDir, dir), { recursive: true });
    }

    return okResult(true);
  } catch (cause) {
    return errResult(new Error("could not create new xdg dirs", { cause }));
  }
}

async function removeDirs(homeDir: string, dirs: string[]): ResultAsync<true, Error> {
  try {
    for await (const dir of dirs) {
      await fs.rm(path.join(homeDir, dir), { recursive: true, force: true });
    }

    return okResult(true);
  } catch (cause) {
    return errResult(new Error("could not remove default xdg dirs", { cause }));
  }
}
