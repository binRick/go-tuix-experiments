#!/usr/bin/env bash
set -e
[[ -d bin ]] || mkdir -p bin
BIN=./bin/simple
go build -o $BIN|| exit 1
reset

he() {
	reset || true
	if [[ -f .err ]]; then
		cat .err|tail -n 10
		unlink .err
	fi
}

trap he EXIT
#while :; do 
eval $BIN ${@:-}
# 2>.err
#sleep 1
#done
