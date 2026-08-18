import { expect, mock, test } from "bun:test";
import assert from "node:assert/strict";
import { createCleanup } from "./cleanup";

test("success", async () => {
  const cleanup = createCleanup();
  const fn = mock<() => Promise<void>>(async () => Promise.resolve());
  cleanup.defer(() => fn());

  const result = await cleanup.run();

  assert(!result.isError);
  expect(result.value).toBeTrue();
  expect(fn).toHaveBeenCalled();
});

test("break run", async () => {
  const cleanup = createCleanup();
  const fn = mock<() => Promise<void>>(async () => Promise.resolve()).mockImplementation(() => {
    throw new Error("fn failed");
  });
  cleanup.defer(() => fn());

  const result = await cleanup.run();

  assert(result.isError);
  expect(result.error).toBeInstanceOf(AggregateError);
  expect(result.error).toMatchObject({
    message: "cleanup failed",
    errors: expect.arrayContaining([new Error("fn failed")]),
  });
  expect(fn).toHaveBeenCalled();
});
