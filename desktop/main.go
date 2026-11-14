package main

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"io"
	"net"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	"github.com/egoist/mygo"
)

func run() error {
	exe, err := os.Executable()
	if err != nil {
		return err
	}
	root := filepath.Dir(exe)
	if runtime.GOOS == "darwin" {
		root = filepath.Clean(filepath.Join(root, "..", "Resources"))
	}
	data, err := os.UserConfigDir()
	if err != nil {
		return err
	}
	data = filepath.Join(data, "XiangyinNotes")
	if err = os.MkdirAll(data, 0700); err != nil {
		return err
	}
	passwordPath := filepath.Join(data, "desktop.key")
	password, err := os.ReadFile(passwordPath)
	if os.IsNotExist(err) {
		key := make([]byte, 32)
		if _, err = rand.Read(key); err != nil {
			return err
		}
		password = []byte(hex.EncodeToString(key))
		err = os.WriteFile(passwordPath, password, 0600)
	}
	if err != nil {
		return err
	}
	tokenBytes := make([]byte, 32)
	if _, err = rand.Read(tokenBytes); err != nil {
		return err
	}
	token := hex.EncodeToString(tokenBytes)
	listener, err := net.Listen("tcp", "127.0.0.1:39093")
	if err != nil {
		return fmt.Errorf("无法启动本地服务（端口 39093 已被占用）：%w", err)
	}
	port := listener.Addr().(*net.TCPAddr).Port
	listener.Close()
	serverName := "server"
	if runtime.GOOS == "windows" {
		serverName += ".exe"
	}
	server := exec.Command(filepath.Join(root, "backend", serverName))
	server.Dir = root
	configureServer(server)
	for _, entry := range os.Environ() {
		key, _, _ := strings.Cut(entry, "=")
		key = strings.ToUpper(key)
		if strings.HasPrefix(key, "XIANGYING_") || key == "DATABASE_PATH" || key == "ASSETS_PATH" || key == "HOST" || key == "PORT" || key == "PUBLIC_URL" || key == "CLIENT_ROOT" || key == "NODE_ENV" || key == "COOKIE_SECURE" {
			continue
		}
		server.Env = append(server.Env, entry)
	}
	server.Env = append(server.Env, "HOST=127.0.0.1", fmt.Sprintf("PORT=%d", port), "NODE_ENV=production", "XIANGYING_DESKTOP_TOKEN="+token, "XIANGYING_USERNAME=desktop", "XIANGYING_PASSWORD="+string(password), "DATABASE_PATH="+filepath.Join(data, "xiangying-notes.sqlite"), "CLIENT_ROOT="+filepath.Join(root, "client"), "COOKIE_SECURE=false")
	logFile, err := os.OpenFile(filepath.Join(data, "desktop.log"), os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0600)
	if err != nil {
		return err
	}
	defer logFile.Close()
	server.Stdout = logFile
	server.Stderr = logFile
	if err = server.Start(); err != nil {
		return err
	}
	defer func() { _ = server.Process.Kill(); _, _ = server.Process.Wait() }()
	base := fmt.Sprintf("http://127.0.0.1:%d", port)
	client := &http.Client{Timeout: time.Second}
	ready := false
	for i := 0; i < 100; i++ {
		response, e := client.Get(base + "/api/health")
		if e == nil {
			_, _ = io.Copy(io.Discard, response.Body)
			response.Body.Close()
			if response.StatusCode == 200 {
				ready = true
				break
			}
		}
		time.Sleep(100 * time.Millisecond)
	}
	if !ready {
		return fmt.Errorf("后端启动失败，请查看 %s", logFile.Name())
	}
	mygo.App.WhenReady(func() {
		window := mygo.NewWindow(mygo.WindowOptions{Title: "象映笔记", URL: base + "/api/desktop-login?token=" + token, Width: 1280, Height: 820, MinWidth: 760, MinHeight: 520, TitleBarStyle: mygo.TitleBarHiddenInset, TitleBarHeight: 40})
		if icon, err := os.ReadFile(filepath.Join(root, "client", "icon-192.png")); err == nil {
			if err := window.SetIcon(icon); err != nil {
				fmt.Fprintln(logFile, "设置窗口图标失败:", err)
			}
		} else {
			fmt.Fprintln(logFile, "读取窗口图标失败:", err)
		}
	})
	return mygo.App.Run()
}
func main() {
	if err := run(); err != nil {
		showError(err)
	}
}
