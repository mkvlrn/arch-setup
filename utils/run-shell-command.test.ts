import { expect, mock, test } from "bun:test";
import assert from "node:assert/strict";
import { runShellCommand } from "./run-shell-command";

test("success", async () => {
  const cmd = mock<() => Promise<unknown>>(() => Promise.resolve());

  const result = await runShellCommand(cmd, "cmd");

  assert(!result.isError);
  expect(result.value).toBeTrue();
});

test("failure", async () => {
  const error = new Error("nope");
  const cmd = mock<() => Promise<unknown>>(() => Promise.reject(error));

  const result = await runShellCommand(cmd, "cmd");

  assert(result.isError);
  expect(result.error).toBeInstanceOf(Error);
  expect(result.error).toMatchObject({
    message: "could not run cmd",
    cause: error,
  });
});
