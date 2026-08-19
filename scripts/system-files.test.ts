/** biome-ignore-all lint/performance/useTopLevelRegex: single use */

import { afterEach, beforeEach, expect, mock, spyOn, test } from "bun:test";
import assert from "node:assert/strict";
import fs from "node:fs/promises";
import os from "node:os";
import path from "node:path";
import { $ } from "bun";
import { systemFiles } from "./system-files";

let mockEtc = "";
let makepkgPath = "";
let pacmanPath = "";

beforeEach(async () => {
  mockEtc = await fs.mkdtemp(path.join(os.tmpdir(), "test-arch-setup-system-files-"));
  makepkgPath = path.join(mockEtc, "makepkg.conf");
  pacmanPath = path.join(mockEtc, "pacman.conf");
  await Bun.write(makepkgPath, `OPTIONS=bad\nMAKEFLAGS="bad"`);
  await Bun.write(pacmanPath, "#Color");
});

afterEach(async () => {
  await fs.rm(mockEtc, { recursive: true, force: true });
  mock.restore();
});

test("copy success", async () => {
  const source = path.join(mockEtc, "source.conf");
  const destination = path.join(mockEtc, "copied.conf");
  await Bun.write(source, "hello");

  const result = await systemFiles({
    copies: {
      [source]: destination,
    },
    install: (src: string, dest: string) => $`install -m 644 ${src} ${dest}`,
  });

  assert(!result.isError);
  expect(await Bun.file(destination).text()).toBe("hello");
});

test("update success", async () => {
  const result = await systemFiles({
    updates: {
      [makepkgPath]: [
        [/^OPTIONS=.*$/m, "OPTIONS=good"],
        [/^MAKEFLAGS=.*$/m, `MAKEFLAGS="good"`],
      ],
      [pacmanPath]: [[/^#Color$/, "Color"]],
    },
    install: (src: string, dest: string) => $`install -m 644 ${src} ${dest}`,
  });
  const updatedMakepkg = await Bun.file(path.join(mockEtc, "makepkg.conf")).text();
  const updatedPacman = await Bun.file(path.join(mockEtc, "pacman.conf")).text();

  assert(!result.isError);
  expect(updatedMakepkg).toEqual(`OPTIONS=good\nMAKEFLAGS="good"`);
  expect(updatedPacman).toEqual("Color");
});

test("break doUpdate", async () => {
  const error = new Error("mkdtempDisposable broke");
  spyOn(fs, "mkdtempDisposable").mockRejectedValue(error);
  const result = await systemFiles({
    updates: {
      [makepkgPath]: [
        [/^OPTIONS=.*$/m, "OPTIONS=good"],
        [/^MAKEFLAGS=.*$/m, `MAKEFLAGS="good"`],
      ],
      [pacmanPath]: [[/^#Color$/, "Color"]],
    },
    install: (src: string, dest: string) => $`install -m 644 ${src} ${dest}`,
  });

  assert(result.isError);
  expect(result.error).toBeInstanceOf(Error);
  expect(result.error).toMatchObject({
    message: `could not update ${path.join(mockEtc, "makepkg.conf")}`,
    cause: error,
  });
});

test("break doCopy", async () => {
  const source = path.join(mockEtc, "source.conf");
  const destination = path.join(mockEtc, "copied.conf");
  const error = new Error("install broke");
  await Bun.write(source, "hello");
  const mockInstall = mock<(src: string, dest: string) => $.ShellPromise>(() => $`false`);

  const result = await systemFiles({
    copies: {
      [source]: destination,
    },
    install: mockInstall,
  });

  assert(result.isError);
  expect(result.error).toBeInstanceOf(Error);
  expect(result.error).toMatchObject({
    message: "could not copy system files",
    cause: error,
  });
  expect(mockInstall).toHaveBeenCalledWith(source, destination);
});
