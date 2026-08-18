/** biome-ignore-all lint/performance/useTopLevelRegex: single use */

import { afterEach, beforeEach, expect, mock, spyOn, test } from "bun:test";
import assert from "node:assert/strict";
import fs from "node:fs/promises";
import os from "node:os";
import path from "node:path";
import { $ } from "bun";
import { updateConfs } from "./update-confs";

let mockEtc = "";

beforeEach(async () => {
  mockEtc = await fs.mkdtemp(path.join(os.tmpdir(), "arch-setup-etc-"));
  await Bun.write(path.join(mockEtc, "makepkg.conf"), `OPTIONS=bad\nMAKEFLAGS="bad"`);
  await Bun.write(path.join(mockEtc, "pacman.conf"), "#Color");
});

afterEach(async () => {
  await fs.rm(mockEtc, { recursive: true, force: true });
  mock.restore();
});

test("success", async () => {
  const result = await updateConfs({
    etcDir: mockEtc,
    files: {
      "makepkg.conf": [
        [/^OPTIONS=.*$/m, "OPTIONS=good"],
        [/^MAKEFLAGS=.*$/m, `MAKEFLAGS="good"`],
      ],
      "pacman.conf": [[/^#Color$/, "Color"]],
    },
    sudoInstal: (origin: string, destination: string) => $`install -m 644 ${origin} ${destination}`,
  });
  const updatedMakepkg = await Bun.file(path.join(mockEtc, "makepkg.conf")).text();
  const updatedPacman = await Bun.file(path.join(mockEtc, "pacman.conf")).text();

  assert(!result.isError);
  expect(updatedMakepkg).toEqual(`OPTIONS=good\nMAKEFLAGS="good"`);
  expect(updatedPacman).toEqual("Color");
});

test("break updateConf", async () => {
  const error = new Error("mkdtempDisposable broke");
  spyOn(fs, "mkdtempDisposable").mockRejectedValue(error);
  const result = await updateConfs({
    etcDir: mockEtc,
    files: {
      "makepkg.conf": [
        [/^OPTIONS=.*$/m, "OPTIONS=good"],
        [/^MAKEFLAGS=.*$/m, `MAKEFLAGS="good"`],
      ],
      "pacman.conf": [[/^#Color$/, "Color"]],
    },
    sudoInstal: (origin: string, destination: string) => $`install -m 644 ${origin} ${destination}`,
  });

  assert(result.isError);
  expect(result.error).toBeInstanceOf(Error);
  expect(result.error).toMatchObject({
    message: "could not update /etc/makepkg.conf",
    cause: error,
  });
});
