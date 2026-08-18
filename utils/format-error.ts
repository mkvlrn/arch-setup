export function formatError(error: unknown): string {
  if (!(error instanceof Error)) {
    return String(error);
  }

  const output = [error.stack ?? error.message];
  let current = error.cause;

  while (current instanceof Error) {
    output.push(`Caused by:\n${current.stack ?? current.message}`);
    current = current.cause;
  }

  return output.join("\n\n");
}
