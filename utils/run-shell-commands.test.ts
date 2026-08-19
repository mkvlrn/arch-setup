import { expect, test } from "bun:test";
import assert from "node:assert/strict";
import { $ } from "bun";
import { runShellCommands } from "./run-shell-commands";

test("runs commands successfully", async () => {
  const result = await runShellCommands([
    { exec: () => $`true`, name: "cmd", quiet: false },
    { exec: () => $`true`, name: "quietCmd", quiet: true },
  ]);

  assert(!result.isError);
  expect(result.value).toBeTrue();
});

test("returns error when command fails", async () => {
  const result = await runShellCommands([{ exec: () => $`false`, name: "cmd", quiet: false }]);

  assert(result.isError);
  expect(result.error).toBeInstanceOf(Error);
  expect(result.error.message).toEqual("could not run cmd");
});
