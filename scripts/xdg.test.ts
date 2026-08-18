import { afterEach, expect, mock, spyOn, test } from "bun:test";
import assert from "node:assert/strict";
import fs from "node:fs/promises";
import { $ } from "bun";
import { xdg } from "./xdg";

afterEach(() => {
  mock.restore();
});

test("success", async () => {
  await using tempDir = await fs.mkdtempDisposable("/tmp/");
  const xdgCmd = mock<() => $.ShellPromise>(() => $`true`);

  const result = await xdg({
    homeDir: tempDir.path,
    mkdir: ["keep1", "keep2", "removeme"],
    rmrf: ["removeme"],
    xdgCmd,
  });
  const existing = (await fs.readdir(tempDir.path, { withFileTypes: true }))
    .filter((f) => f.isDirectory())
    .map((f) => f.name);

  assert(!result.isError);
  expect(result.value).toBeTrue();
  expect(existing.toSorted()).toEqual(["keep1", "keep2"]);
  expect(xdgCmd).toHaveBeenCalled();
});

test("break fn.mkdir", async () => {
  const error = new Error("mkdir broke");
  const xdgCmd = mock<() => $.ShellPromise>(() => $`true`);
  spyOn(fs, "mkdir").mockRejectedValue(error);

  const result = await xdg({
    homeDir: "tempDir",
    mkdir: ["keep1", "keep2", "removeme"],
    rmrf: ["removeme"],
    xdgCmd,
  });

  assert(result.isError);
  expect(result.error).toBeInstanceOf(Error);
  expect(result.error).toMatchObject({
    message: "could not create new xdg dirs",
    cause: error,
  });
  expect(xdgCmd).not.toHaveBeenCalled();
});

test("break fn.rm", async () => {
  const error = new Error("rm broke");
  await using tempDir = await fs.mkdtempDisposable("/tmp/");
  const xdgCmd = mock<() => $.ShellPromise>(() => $`true`);
  spyOn(fs, "rm").mockRejectedValue(error);

  const result = await xdg({
    homeDir: tempDir.path,
    mkdir: ["keep1", "keep2", "removeme"],
    rmrf: ["removeme"],
    xdgCmd,
  });

  assert(result.isError);
  expect(result.error).toBeInstanceOf(Error);
  expect(result.error).toMatchObject({
    message: "could not remove default xdg dirs",
    cause: error,
  });
  expect(xdgCmd).not.toHaveBeenCalled();
});

test("break xdgCommand", async () => {
  await using tempDir = await fs.mkdtempDisposable("/tmp/");
  const xdgCmd = mock<() => $.ShellPromise>(() => $`false`);

  const result = await xdg({
    homeDir: tempDir.path,
    mkdir: ["keep1", "keep2", "removeme"],
    rmrf: ["removeme"],
    xdgCmd,
  });

  assert(result.isError);
  expect(result.error).toBeInstanceOf(Error);
  expect(result.error).toMatchObject({ message: "could not run xdg dirs update" });
  expect(xdgCmd).toHaveBeenCalled();
});
