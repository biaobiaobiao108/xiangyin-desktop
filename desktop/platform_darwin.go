package main

import (
	"fmt"
	"os/exec"
)

func configureServer(_ *exec.Cmd) {}
func showError(err error) {
	_ = exec.Command("osascript", "-e", `on run argv
display dialog (item 1 of argv) with title "象映笔记" buttons {"好"} default button 1 with icon stop
end run`, err.Error()).Run()
	fmt.Println(err)
}
