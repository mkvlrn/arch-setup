import { expect, test } from "bun:test";
import assert from "node:assert/strict";
import { $ } from "bun";
import { runShellCommand } from "./run-shell-command";

test("success", async () => {
  const result = await runShellCommand(() => $`true`, "cmd");
  const resultQuiet = await runShellCommand(() => $`true`, "quietCmd", true);

  assert(!result.isError);
  assert(!resultQuiet.isError);
  expect(result.value).toBeTrue();
  expect(resultQuiet.value).toBeTrue();
});

test("failure", async () => {
  const result = await runShellCommand(() => $`false`, "cmd");

  assert(result.isError);
  expect(result.error).toBeInstanceOf(Error);
  expect(result.error).toMatchObject({ message: "could not run cmd" });
});
