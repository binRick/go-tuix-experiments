package sexpect

import (
	"strings"

	"github.com/gobwas/glob"
	"github.com/shirou/gopsutil/process"
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
