export function formatError(error: unknown): string {
  if (!(error instanceof Error)) {
    return String(error);
  }

  const output = [error.stack ?? error.message];

  let current = error.cause;

  while (current instanceof Error) {
    output.push(`Caused by:\n${current.stack ?? current.message}`);

    if ("stdout" in current && current.stdout instanceof Uint8Array && current.stdout.length > 0) {
      output.push(`stdout:\n${Buffer.from(current.stdout).toString().trim()}`);
    }

    if ("stderr" in current && current.stderr instanceof Uint8Array && current.stderr.length > 0) {
      output.push(`stderr:\n${Buffer.from(current.stderr).toString().trim()}`);
    }

    current = current.cause;
  }

  return output.join("\n\n");
}
