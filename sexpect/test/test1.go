package main

import (
	"os"

	"github.com/k0kubun/pp"
	sexpect "local.dev/sexpect"
)

func main() {
	pp.Println(`ok2`)
	pp.Fprintf(os.Stderr, "ok3\n")
	sessions, _ := sexpect.SexpectSessions()
	pp.Println(sessions)
	os.Exit(1)
}
