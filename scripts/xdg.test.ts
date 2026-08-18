import { afterEach, expect, mock, spyOn, test } from "bun:test";
import assert from "node:assert/strict";
import fs from "node:fs/promises";
import { xdg } from "./xdg";

afterEach(() => {
  mock.restore();
});

test("success", async () => {
  await using tempDir = await fs.mkdtempDisposable("/tmp/");
  const xdgCommand = mock<() => Promise<void>>(async () => Promise.resolve());

  const result = await xdg({
    homeDir: tempDir.path,
    mkdir: ["keep1", "keep2", "removeme"],
    rmrf: ["removeme"],
    xdgCmd: xdgCommand,
  });
  const existing = (await fs.readdir(tempDir.path, { withFileTypes: true }))
    .filter((f) => f.isDirectory())
    .map((f) => f.name);

  assert(!result.isError);
  expect(result.value).toBeTrue();
  expect(existing.toSorted()).toEqual(["keep1", "keep2"]);
  expect(xdgCommand).toHaveBeenCalled();
});

test("break fn.mkdir", async () => {
  const error = new Error("mkdir broke");
  const xdgCommand = mock<() => Promise<void>>(async () => Promise.resolve());
  spyOn(fs, "mkdir").mockRejectedValue(error);

  const result = await xdg({
    homeDir: "tempDir",
    mkdir: ["keep1", "keep2", "removeme"],
    rmrf: ["removeme"],
    xdgCmd: xdgCommand,
  });

  assert(result.isError);
  expect(result.error).toBeInstanceOf(Error);
  expect(result.error).toMatchObject({
    message: "could not create new xdg dirs",
    cause: error,
  });
  expect(xdgCommand).not.toHaveBeenCalled();
});

test("break fn.rm", async () => {
  const error = new Error("rm broke");
  await using tempDir = await fs.mkdtempDisposable("/tmp/");
  const xdgCommand = mock<() => Promise<void>>(async () => Promise.resolve());
  spyOn(fs, "rm").mockRejectedValue(error);

  const result = await xdg({
    homeDir: tempDir.path,
    mkdir: ["keep1", "keep2", "removeme"],
    rmrf: ["removeme"],
    xdgCmd: xdgCommand,
  });

  assert(result.isError);
  expect(result.error).toBeInstanceOf(Error);
  expect(result.error).toMatchObject({
    message: "could not remove default xdg dirs",
    cause: error,
  });
  expect(xdgCommand).not.toHaveBeenCalled();
});

test("break xdgCommand", async () => {
  const error = new Error("xdg is evil");
  await using tempDir = await fs.mkdtempDisposable("/tmp/");
  const xdgCommand = mock<() => Promise<void>>(async () => Promise.reject(error));

  const result = await xdg({
    homeDir: tempDir.path,
    mkdir: ["keep1", "keep2", "removeme"],
    rmrf: ["removeme"],
    xdgCmd: xdgCommand,
  });

  assert(result.isError);
  expect(result.error).toBeInstanceOf(Error);
  expect(result.error).toMatchObject({
    message: "could not run xdgCmd",
    cause: error,
  });
  expect(xdgCommand).toHaveBeenCalled();
});
