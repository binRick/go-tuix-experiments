package sexpect

import (
	process "github.com/shirou/gopsutil/v3/process"
)

func SexpectSessions() ([]Session, error) {
	var ass []Session
	procs, err := process.Processes()
	if err != nil {
		panic(err)
	}
	for _, proc := range procs {
		n, err := proc.Name()
		if err != nil || n != `sexpect` {
			continue
		}
		ct, _ := proc.CreateTime()
		mp, _ := proc.MemoryPercent()
		cp, _ := proc.CPUPercent()
		cmdl, _ := proc.Cmdline()
		cwd, _ := proc.Cwd()
		st, _ := proc.Status()
		term, _ := proc.Terminal()
		conns, _ := proc.Connections()
		un, _ := proc.Username()
		of, _ := proc.OpenFiles()
		ass = append(ass, Session{
			PID:            int(proc.Pid),
			CreateTime:     ct,
			Cmdline:        cmdl,
			MemoryPercent:  mp,
			CPUPercent:     cp,
			Cwd:            cwd,
			Terminal:       term,
			Status:         st,
			Username:       un,
			OpenFilesQty:   int32(len(of)),
			ConnectionsQty: int32(len(conns)),
		})
	}
	return ass, nil
}
