#!/bin/bash

set -e
ROOT_DIR="$(cd "$(dirname "$0")" && pwd)"
source "$ROOT_DIR/dev.sh"
go run "$ROOT_DIR/pkg/delta/main.go"
