import packageJson from "../package.json" with { type: "json" };

export default {
  name: "象映笔记",
  slug: "xiangyin-notes",
  identifier: "com.xiangyin.notes.desktop",
  version: packageJson.version,
  macos: {
    minimumSystemVersion: "13.0",
    helperEntitlements: { "backend/server": "bun-entitlements.plist" },
  },
  icon: "../app/icon-512.png",
  out: "../dist/packages",
  resources: [
    "../dist/client",
    "../dist/desktop-resources/database",
    "../dist/desktop-resources/backend",
  ],
};
