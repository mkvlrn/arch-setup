import os from "node:os";
import path from "node:path";
import type { Step } from "../utils/types";

export function miscUser(): Step {
  return {
    label: "Set user shell and docker group situation",
    commands: [
      ["set user shell", Bun.$`sudo chsh -s $(which fish) ${os.userInfo().username}`],
      ["add user to docker group", Bun.$`sudo usermod -aG docker ${os.userInfo().username}`],
      [
        "set anonymous ftp user root dir",
        Bun.$`sudo usermod -d ${path.join(os.homedir(), "downloads", "torrents")} ftp`,
      ],
      ["allow ftp user to traverse to download dir", Bun.$`chmod o+x ${os.homedir()}`],
      ["start docker service", Bun.$`sudo systemctl enable --now docker.socket`],
      ["start pure-fptd service", Bun.$`sudo systemctl enable --now pure-ftpd.socket`],
      [
        "create fish completions dir",
        Bun.$`mkdir -p ${path.join(os.homedir(), ".config", "fish", "completions")}`,
      ],
    ],
  };
}
