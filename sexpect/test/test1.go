package main

import (
	"os"

	"github.com/k0kubun/pp"
	sexpect "local.dev/sexpect"
)

func main() {
	sessions, _ := sexpect.SexpectSessions()
	pp.Fprintf(os.Stderr, "%s Sessions: %s\n", len(sessions), sessions)
	os.Exit(0)
}
