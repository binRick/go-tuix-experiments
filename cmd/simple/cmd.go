package main

import (
	"context"
	"io"
	"os/exec"
	"sync"
	"time"

	"github.com/rivo/tview"
	logrus "github.com/sirupsen/logrus"
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

func run_cmd_in_context(ctx context.Context, cmd string, stdout, stderr *tview.TextView) {
	c := exec.CommandContext(ctx, "sh", "-c", cmd)
	e, err := c.StderrPipe()
	if err != nil {
		panic(err)
	}
	o, err := c.StdoutPipe()
	if err != nil {
		panic(err)
	}
	go copyAndCapture(stdout, o)
	go copyAndCapture(stderr, e)
	l.WithFields(logrus.Fields{"cmd": cmd}).Info("Command Started inside existing context")
	serr := c.Start()
	if serr != nil {
		panic(serr)
	}
	l.WithFields(logrus.Fields{"cmd": cmd}).Info("Command waited inside existing context")
	start := time.Now()
	err = c.Wait()
	l.WithFields(logrus.Fields{"cmd": cmd, "dur": time.Since(start)}).Info("Command ended inside existing context")
	if err != nil {
		panic(err)
	}
}

func run_cmd_in_new_context(cmd string, stdout, stderr *tview.TextView) {
	start := time.Now()
	cmd_signal := make(chan int)
	var wg sync.WaitGroup
	ctx, cancel_cmd_ctx := context.WithCancel(context.Background())
	wg.Add(1)
	go func() {
		l.WithFields(logrus.Fields{"cmd": cmd}).Info("Command Cancel Signal Listener Created. Waiting on cmd_signal")
		start := time.Now()
		recvd_cmd_signal := <-cmd_signal
		l.WithFields(logrus.Fields{"cmd": cmd, "dur": time.Since(start), "recvd_signal": recvd_cmd_signal}).Info("Command Cancel Signal Listener Activated")
		//	if recvd_cmd_signal == 0 { //  signal 0 is exit from wg.Done goroutine
		//} else {
		if false {
			cancel_cmd_ctx()
		}
		//		}
		l.WithFields(logrus.Fields{"cmd": cmd}).Info("Command Cancel Signal Listener Completed")
	}()

	go func() {
		start := time.Now()
		l.WithFields(logrus.Fields{"cmd": cmd}).Info("Starting Command in new Context")
		run_cmd_in_context(ctx, cmd, stdout, stderr)
		l.WithFields(logrus.Fields{"cmd": cmd, "dur": time.Since(start)}).Info("Command Ended")
		cmd_signal <- 0
		l.WithFields(logrus.Fields{"cmd": cmd, "dur": time.Since(start)}).Info("Command Wait Group Notified")
		wg.Done()
	}()

	//	<-time.After(3 * time.Second)
	//	log.Println("closing via ctx")
	//	cancel()

	// Wait for the child goroutine to finish, which will only occur when
	// the child process has stopped and the call to cmd.Wait has returned.
	// This prevents main() exiting prematurely.
	l.WithFields(logrus.Fields{"cmd": cmd, "dur": time.Since(start)}).Info("Command Runner Waiting")
	wg.Wait()
	l.WithFields(logrus.Fields{"cmd": cmd, "dur": time.Since(start)}).Info("Command Runner Completed")
}
