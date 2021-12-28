module simple

go 1.16

require (
	github.com/binRick/abduco-dev/go/abducoctl v0.0.0-00010101000000-000000000000
	github.com/gdamore/tcell/v2 v2.4.1-0.20210905002822-f057f0a857a1
	github.com/k0kubun/pp v3.0.1+incompatible
	github.com/millerlogic/tuix v0.0.0-20210130203550-953ce41af824
	github.com/nxadm/tail v1.4.8
	github.com/rivo/tview v0.0.0-20211202162923-2a6de950f73b
	github.com/sirupsen/logrus v1.8.1
)

require local.dev/sexpect v0.0.0-00010101000000-000000000000

replace github.com/binRick/abduco-dev/go/abducoctl => ../../../abduco-dev/go/abducoctl

replace local.dev/sexpect => ./../../sexpect/.
