package main

import (
	"os"

	"github.com/k0kubun/pp"
	sexpect "local.dev/sexpect"
)

func main() {
	sessions, _ := sexpect.SexpectSessions()
	pp.Println(sessions)
	os.Exit(0)
}
