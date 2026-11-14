import { cp, mkdir, rm } from "node:fs/promises";
import { resolve } from "node:path";

if (process.platform !== "win32" && process.platform !== "darwin")
  throw new Error("桌面版需要在 Windows 或 macOS 上构建");

const root = resolve(import.meta.dir, "..");
const resources = resolve(root, "dist/desktop-resources");
await rm(resources, { recursive: true, force: true });
await mkdir(resolve(resources, "database"), { recursive: true });
await mkdir(resolve(resources, "backend"), { recursive: true });

async function run(command: string[], cwd = root) {
  const child = Bun.spawn(command, { cwd, stdout: "inherit", stderr: "inherit" });
  if ((await child.exited) !== 0) throw new Error(`${command[0]} 构建失败`);
}

await run(["bun", "run", "build:client"]);
await cp(resolve(root, "database/schema.sql"), resolve(resources, "database/schema.sql"));
await run([
  "bun",
  "build",
  "./server/index.ts",
  "--compile",
  "--outfile",
  resolve(resources, "backend", process.platform === "win32" ? "server.exe" : "server"),
]);
await run(
  ["go", "run", "github.com/egoist/mygo/cmd/mygo@v0.9.0", "build"],
  resolve(root, "desktop"),
);
