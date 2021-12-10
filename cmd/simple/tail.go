package main

import (
	"fmt"
	"sync"

	"github.com/nxadm/tail"
	logrus "github.com/sirupsen/logrus"
)

var open_logs = []string{}
var open_logs_mutex sync.Mutex

func tail_log(session string) {
	log_file := get_session_log_path(session)
	open_logs_mutex.Lock()
	has := false
	for _, ol := range open_logs {
		if ol == log_file {
			has = true
		}
	}
	if !has {
		open_logs = append(open_logs, log_file)
	}
	open_logs_mutex.Unlock()
	if has {
		return
	}
	var seek = tail.SeekInfo{Offset: 0, Whence: 2}
	tailer, err := tail.TailFile(log_file, tail.Config{
		Follow:   true,
		ReOpen:   true,
		Location: &seek,
	})
	if err != nil {
		panic(err)
	}
	for line := range tailer.Lines {
		l.WithFields(logrus.Fields{
			"bytes": len(line.Text),
			"log":   log_file,
		}).Info(fmt.Sprintf(`TAIL- %d bytes: %s`, len(line.Text), line.Text))
		fmt.Fprintf(log_view, "%s> %s\n", session, line.Text)
		//	log_view.ScrollToEnd()
	}
}
