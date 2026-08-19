import { expect, test } from "bun:test";
import { formatError } from "./format-error";

test("formats non-error values", () => {
  expect(formatError("nope")).toBe("nope");
  expect(formatError(42)).toBe("42");
  expect(formatError(null)).toBe("null");
});

test("formats error", () => {
  const error = new Error("top error");

  const formatted = formatError(error);

  expect(formatted).toContain("Error: top error");
});

test("formats error cause", () => {
  const cause = new Error("inner error");
  const error = new Error("top error", { cause });

  const formatted = formatError(error);

  expect(formatted).toContain("Error: top error");
  expect(formatted).toContain("Caused by:\nError: inner error");
});

test("formats nested error causes", () => {
  const bottom = new Error("bottom error");
  const middle = new Error("middle error", { cause: bottom });
  const top = new Error("top error", { cause: middle });

  const formatted = formatError(top);

  expect(formatted).toContain("Error: top error");
  expect(formatted).toContain("Caused by:\nError: middle error");
  expect(formatted).toContain("Caused by:\nError: bottom error");
});

test("formats stdout and stderr", () => {
  const cause = Object.assign(new Error("shell broke"), {
    stdout: Buffer.from("some stdout\n"),
    stderr: Buffer.from("some stderr\n"),
  });
  const error = new Error("command failed", { cause });

  const formatted = formatError(error);

  expect(formatted).toContain("stdout:\nsome stdout");
  expect(formatted).toContain("stderr:\nsome stderr");
});

test("ignores empty stdout and stderr", () => {
  const cause = Object.assign(new Error("shell broke"), {
    stdout: Buffer.from(""),
    stderr: Buffer.from(""),
  });
  const error = new Error("command failed", { cause });

  const formatted = formatError(error);

  expect(formatted).not.toContain("stdout:");
  expect(formatted).not.toContain("stderr:");
});

test("stops at non-error cause", () => {
  const error = new Error("top error", { cause: "not an error" });

  const formatted = formatError(error);

  expect(formatted).toContain("Error: top error");
  expect(formatted).not.toContain("Caused by:");
});
