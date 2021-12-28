module test1

go 1.16

replace local.dev/sexpect => ./../../sexpect/.

require (
	github.com/k0kubun/pp v3.0.1+incompatible
	github.com/shirou/gopsutil v3.21.11+incompatible // indirect
	local.dev/sexpect v0.0.0-00010101000000-000000000000
)
