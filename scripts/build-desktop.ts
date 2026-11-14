import { cp, mkdir, rm } from "node:fs/promises";
import { resolve } from "node:path";

if (process.platform !== "win32") throw new Error("Windows 桌面版需要在 Windows 上构建");

const root = resolve(import.meta.dir, "..");
const output = resolve(root, "dist/xiangyin-desktop");
await rm(output, { recursive: true, force: true });
await mkdir(resolve(output, "database"), { recursive: true });

async function run(command: string[], cwd = root) {
  const child = Bun.spawn(command, { cwd, stdout: "inherit", stderr: "inherit" });
  if ((await child.exited) !== 0) throw new Error(`${command[0]} 构建失败`);
}

await run(["bun", "run", "build:client"]);
await cp(resolve(root, "dist/client"), resolve(output, "client"), { recursive: true });
await cp(resolve(root, "database/schema.sql"), resolve(output, "database/schema.sql"));
await run([
  "bun",
  "build",
  "./server/index.ts",
  "--compile",
  "--outfile",
  resolve(output, "server.exe"),
]);
await run(
  ["go", "build", "-ldflags=-H=windowsgui", "-o", resolve(output, "象映笔记.exe"), "."],
  resolve(root, "desktop"),
);
console.log(`桌面应用已生成：${output}`);
