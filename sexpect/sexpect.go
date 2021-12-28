package sexpect

import (
	"strings"

	"github.com/gobwas/glob"
	process "github.com/shirou/gopsutil/v3/process"
)

var EXTRACTED_ENV_GLOBS = []string{
	`STD*_LOG_FILE`,
}

func get_extracted_env(proc *process.Process) map[string]string {
	extracted_env := map[string]string{}
	environ, err := proc.Environ()
	if err != nil {
		panic(err)
	}
	for _, kv := range environ {
		if len(kv) < 3 || len(strings.Split(kv, `=`)) < 2 {
			continue
		}
		s := strings.Split(kv, `=`)
		k := s[0]
		v := s[1]
		for _, eg := range EXTRACTED_ENV_GLOBS {
			g := glob.MustCompile(eg)
			match := g.Match(k)
			if match {
				extracted_env[k] = v
			}
		}
	}
	return extracted_env
}

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
		env := get_extracted_env(proc)
		ass = append(ass, Session{
			PID:            int(proc.Pid),
			CreateTime:     ct,
			Cmdline:        cmdl,
			Environ:        env,
			StdoutLog:      env["STDOUT_LOG_FILE"],
			StderrLog:      env["STDERR_LOG_FILE"],
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
