package main

import (
	"io"
	"os/exec"

	"github.com/rivo/tview"
)

func run_cmd(shell_cmd string, stdout, stderr *tview.TextView) {
	cmd := exec.Command("env", "sh", "-c", shell_cmd)
	e, err := cmd.StderrPipe()
	if err != nil {
		panic(err)
	}
	o, err := cmd.StdoutPipe()
	if err != nil {
		panic(err)
	}
	go io.Copy(stdout, o)
	go io.Copy(stderr, e)
	_ = cmd.Start()
	_ = cmd.Wait()
}
