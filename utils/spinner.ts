import process from "node:process";
export function createSpinner(text: string) {
  if (!process.stdout.isTTY) {
    console.log(text);

    return {
      success() {},
      failure() {},
    };
  }

  const frames = ["⠋", "⠙", "⠹", "⠸", "⠼", "⠴", "⠦", "⠧", "⠇", "⠏"];
  let index = 0;

  process.stdout.write(`${frames[index]} ${text}`);

  const timer = setInterval(() => {
    index = (index + 1) % frames.length;
    process.stdout.write(`\r${frames[index]} ${text}`);
  }, 80);

  return {
    success() {
      clearInterval(timer);
      process.stdout.write(`\r✓ ${text}\n`);
    },

    failure() {
      clearInterval(timer);
      process.stdout.write(`\r✗ ${text}\n`);
    },
  };
}
