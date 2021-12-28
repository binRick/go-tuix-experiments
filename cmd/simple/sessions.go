package main

import (
	"sync"
	"time"

	sexpect "local.dev/sexpect"
)

var current_sessions = []sexpect.Session{}
var current_sessions_dur = time.Duration(0)
var cur_sessions_mutex sync.Mutex

func update_sexpect_sessions() {
	cur_sessions_mutex.Lock()
	s := time.Now()

	sessions, err := sexpect.SexpectSessions()
	if err != nil {
		panic(err)
	}
	current_sessions = sessions
	current_sessions_dur = time.Since(s)
	cur_sessions_mutex.Unlock()
}

func get_sexpect_sessions() []sexpect.Session {
	cur_sessions_mutex.Lock()
	cur := current_sessions
	cur_sessions_mutex.Unlock()
	return cur
}
