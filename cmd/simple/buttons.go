package main

import (
	"fmt"
	"os"

	"github.com/rivo/tview"
	logrus "github.com/sirupsen/logrus"
)

var prev_cmd = ``

func AddCommandButton(form *tview.Form, title string, command string, cb func()) {
	form.AddButton(title, func() {
		cmd_mutex.Lock()
		defer cmd_mutex.Unlock()
		prev_cmd = fmt.Sprintf(`%s`, cur_cmd)
		cur_cmd = command
		go run_cmd_in_new_context(cur_cmd, tv, tvr)
		if cb != nil {
			cb()
		}
		l.WithFields(logrus.Fields{
			"prev_cmd": prev_cmd,
			"cur_cmd":  cur_cmd,
			"title":    title,
			"pid":      os.Getpid(),
		}).Info(fmt.Sprintf("New Command Loaded: %s", title))
	})
}
