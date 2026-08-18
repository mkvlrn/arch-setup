import fs from "node:fs/promises";

export const steps = (await fs.readdir("scripts")).filter((f) => !f.endsWith("test.ts")).length;

export const labels: Record<string, string> = {
  xdg: "Configuring XDG user dirs",
  "update-conf": "Configuring pacman",
  yay: "Installing yay and refreshing mirrors",
};
