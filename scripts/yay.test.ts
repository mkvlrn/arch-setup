import { expect, test } from "bun:test";
import assert from "node:assert/strict";
import { yay } from "./yay";

test("success", async () => {
  const result = await yay({ yayCmds: [[() => Promise.resolve(), "test"]] });

  assert(!result.isError);
  expect(result.value).toBeTrue();
});

test("failure", async () => {
  const error = new Error("i'm not feeling it");
  const result = await yay({ yayCmds: [[() => Promise.reject(error), "test"]] });

  assert(result.isError);
  expect(result.error).toBeInstanceOf(Error);
  expect(result.error).toMatchObject({
    message: "could not run test",
    cause: error,
  });
});
