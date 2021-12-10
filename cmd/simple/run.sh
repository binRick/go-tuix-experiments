#!/usr/bin/env bash
set -e
[[ -d bin ]] || mkdir -p bin
BIN=./bin/simple
go build -o $BIN
reset

he() {
	reset || true
	if [[ -f .err ]]; then
		cat .err
		unlink .err
	fi
}

trap he EXIT
#while :; do 
eval $BIN ${@:-} 2>.err
#sleep 1
#done
