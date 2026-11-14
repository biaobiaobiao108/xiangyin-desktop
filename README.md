# 象映笔记桌面版

象映笔记是一款本地 Markdown 笔记应用，支持笔记本、全文搜索、双向链接、任务、随记和图片附件。桌面版使用 MyGo 的系统 WebView 和本地 Bun 后端，数据保存在本机。

## 下载与使用

从 GitHub Releases 下载对应平台的压缩包：

- Windows x64：解压后运行 `象映笔记.exe`，请保留同一目录下的所有文件。需要 Microsoft Edge WebView2 Runtime。
- macOS Apple Silicon：解压后将 `象映笔记.app` 移到应用程序文件夹并打开。macOS 版本未经过 Apple 公证，首次打开可能需要在系统设置的隐私与安全性中允许打开。

本地服务只监听 `127.0.0.1:39093`；如果端口已被其他程序占用，应用会提示启动失败。桌面版使用独立的本地数据库，不会自动连接或同步已有服务上的笔记。

## 数据和备份

退出应用后备份整个数据目录，其中包括数据库、图片附件和本地登录密钥：

- Windows：`%APPDATA%\XiangyinNotes`
- macOS：`~/Library/Application Support/XiangyinNotes`

更新或替换应用程序不会删除该目录。请定期备份，并在恢复前退出应用。

## 本地开发与构建

需要 Bun 1.4.4 和 Go 1.27.4。Windows x64 或 Apple Silicon Mac 上运行：

```sh
bun install --frozen-lockfile
bun run build:desktop
```

MyGo 构建产物位于 `dist/packages/`。该命令编译前端及 Bun 后端，并将所需资源和项目图标打包。运行 `bun run dev` 可在浏览器中开发前端与服务端。

提交前依次运行 `bun run format`、`bun run typecheck`、`bun test` 和 `bun run build`。推送与 `package.json` 版本相同的 `v` 前缀 tag 会触发 GitHub Actions，分别打包 Windows amd64 和 macOS arm64，并创建附带压缩包的 GitHub Release。
