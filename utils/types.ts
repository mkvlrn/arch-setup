export type Command = [string, Bun.$.ShellPromise];

export type Step = { label: string; commands: Command[] };
