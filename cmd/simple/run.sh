#!/usr/bin/env bash
set -e
[[ -d bin ]] || mkdir -p bin
BIN=./bin/simple

go build -o $BIN
exec $BIN ${@:-}
