import { errResult, okResult, type ResultAsync } from "@mkvlrn/result";

export function createCleanup() {
  const tasks: (() => Promise<void>)[] = [];

  return {
    defer(fn: () => Promise<void>): void {
      tasks.push(fn);
    },

    async run(): ResultAsync<true, Error> {
      const errors: unknown[] = [];

      for await (const fn of tasks.toReversed()) {
        try {
          await fn();
        } catch (err) {
          errors.push(err);
        }
      }

      tasks.length = 0;

      if (errors.length > 0) {
        return errResult(new AggregateError(errors, "cleanup failed"));
      }

      return okResult(true);
    },
  } as const;
}
