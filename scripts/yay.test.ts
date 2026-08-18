import { expect, test } from "bun:test";
import assert from "node:assert/strict";
import { $ } from "bun";
import { yay } from "./yay";

test("success", async () => {
  const result = await yay({ yayCmds: [[() => $`true`, "test"]] });

  assert(!result.isError);
  expect(result.value).toBeTrue();
});

test("failure", async () => {
  const result = await yay({ yayCmds: [[() => $`false`, "test"]] });

  assert(result.isError);
  expect(result.error).toBeInstanceOf(Error);
  expect(result.error).toMatchObject({ message: "could not run test" });
});
